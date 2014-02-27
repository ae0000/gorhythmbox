package rhythmbox

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Library string
	Db      Rhythmdb
	Artists []Item
	Albums  []Item
	Genres  []Item
}

const (
	RhythmboxClient     = "rhythmbox-client"                          // The actual client to run commands through
	RhythmboxXmlLibrary = "$HOME/.local/share/rhythmbox/rhythmdb.xml" // Do not write to this
)

/*
<rhythmdb version="1.8">
  <entry type="song">
    <title>Timeline</title>
    <genre>DJ Set</genre>
    <artist>Kaempfer &amp; Dietze</artist>
    <album>Timeline</album>
    <duration>4755</duration>
    <file-size>190284861</file-size>
    <location>file:///home/ae/Music/progressive_psy/Kaempfer%20&amp;%20Dietze%20-%20Timeline.mp3</location>
    <mtime>1346541814</mtime>
    <first-seen>1359918484</first-seen>
    <last-seen>1393272093</last-seen>
    <bitrate>320</bitrate>
    <date>733773</date>
    <media-type>audio/mpeg</media-type>
  </entry>
*/

type Rhythmdb struct {
	XMLName xml.Name `xml:"rhythmdb"`
	Version string   `xml:"version,attr"`
	Entries []Entry  `xml:"entry"`
}

type Tracks struct {
	Entries []Entry
}

type Albums struct {
	Items []Item
}

type Entry struct {
	Id          int
	Type        string `xml:"type,attr"`
	Title       string `xml:"title"`
	Genre       string `xml:"genre"`
	Artist      string `xml:"artist"`
	Album       string `xml:"album"`
	Duration    int    `xml:"duration"`
	TrackNumber int    `xml:"track-number"`
	Rating      int    `xml:"rating"`
	PlayCount   int    `xml:"play-count"`
	// FileSize string `xml:"file-size"`
	Location string `xml:"location"`
	// Mtime string `xml:"mtime"`
	FirstSeen int `xml:"first-seen"`
	LastSeen  int `xml:"last-seen"`
	// Bitrate string `xml:"bitrate"`
	// Date string `xml:"date"`
	MediaType string `xml:"media-type"`
	Selected  bool
}

type Item struct {
	Id       int
	Name     string
	Type     string
	Count    int
	Image    string
	HasImage bool
	Entry    Entry
	Tracks   []Entry
}

func (i *Item) SelectTrack(trackid int) {
	for j, t := range i.Tracks {
		if t.Id == trackid {
			i.Tracks[j].Selected = true
			return
		}
	}
}

// Sorters
type ByTrackNumber []Entry
type ByRandom []Entry
type ByArtistE []Entry
type ByArtist []Item
type ByAlbum []Item
type ByGenre []Item

func (a ByTrackNumber) Len() int           { return len(a) }
func (a ByTrackNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTrackNumber) Less(i, j int) bool { return a[i].TrackNumber < a[j].TrackNumber }

func (a ByRandom) Len() int           { return len(a) }
func (a ByRandom) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRandom) Less(i, j int) bool { return RandBool() }

func (a ByArtist) Len() int           { return len(a) }
func (a ByArtist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByArtist) Less(i, j int) bool { return a[i].Entry.Artist < a[j].Entry.Artist }

func (a ByArtistE) Len() int           { return len(a) }
func (a ByArtistE) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByArtistE) Less(i, j int) bool { return a[i].Artist < a[j].Artist }

func (a ByAlbum) Len() int           { return len(a) }
func (a ByAlbum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAlbum) Less(i, j int) bool { return a[i].Entry.Album < a[j].Entry.Album }

func (a ByGenre) Len() int           { return len(a) }
func (a ByGenre) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGenre) Less(i, j int) bool { return a[i].Entry.Genre < a[j].Entry.Genre }

// Read in the library and set everything up for browsing
func (r *Client) Setup() {
	file, err := ioutil.ReadFile(r.Library)
	if err != nil {
		fmt.Printf("[ERRO] Could not load library: %v\n", err)
		return
	}

	err = xml.Unmarshal([]byte(file), &r.Db)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	// Add Id
	for i := 0; i < len(r.Db.Entries); i++ {
		r.Db.Entries[i].Id = i
	}

	// Sort out the unique artists, albums and genres
	for _, e := range r.Db.Entries {
		if len(e.Album) > 0 {
			if !r.AlbumExists(e.Album) {
				item := Item{
					Id:    e.Id,
					Name:  e.Album,
					Type:  "Album",
					Entry: e,
				}
				// Try and get a pic
				item.Image, item.HasImage = r.GetAlbumImage(e.Id)

				r.Albums = append(r.Albums, item)
			}
		}
		if len(e.Artist) > 0 {
			if !r.ArtistExists(e.Artist) {
				item := Item{
					Id:    e.Id,
					Name:  e.Artist,
					Type:  "Artist",
					Count: 1,
					Entry: e,
				}
				r.Artists = append(r.Artists, item)
			} else {
				r.IncrementArtistCount(e.Artist)
			}
		}

		if len(e.Genre) > 0 {
			if !r.GenreExists(e.Genre) {
				item := Item{
					Id:    e.Id,
					Name:  e.Genre,
					Type:  "Genre",
					Count: 1,
					Entry: e,
				}
				r.Genres = append(r.Genres, item)
			} else {
				r.IncrementGenreCount(e.Genre)
			}
		}
	}
}

func (r *Client) IncrementGenreCount(s string) {
	for i, g := range r.Genres {
		if g.Name == s {
			r.Genres[i].Count++
		}
	}
}

func (r *Client) IncrementArtistCount(s string) {
	for i, a := range r.Artists {
		if a.Name == s {
			r.Artists[i].Count++
		}
	}
}

func (r *Client) AlbumExists(s string) bool {
	for _, a := range r.Albums {
		if a.Name == s {
			return true
		}
	}
	return false
}

func (r *Client) ArtistExists(s string) bool {
	for _, a := range r.Artists {
		if a.Name == s {
			return true
		}
	}
	return false
}

func (r *Client) GenreExists(s string) bool {
	for _, a := range r.Genres {
		if a.Name == s {
			return true
		}
	}
	return false
}

func (r *Client) GetAlbums() []Item {
	sort.Sort(ByArtist(r.Albums))
	return r.Albums
}

func (r *Client) GetArtists() []Item {
	sort.Sort(ByArtist(r.Artists))
	return r.Artists
}

func (r *Client) GetGenres() []Item {
	sort.Sort(ByGenre(r.Genres))
	return r.Genres
}

func (r *Client) GetAlbum(id int) Item {
	album := Item{}
	albumName := r.Db.Entries[id].Album

	for _, a := range r.Db.Entries {
		if a.Album == albumName {
			if len(album.Entry.Album) == 0 {
				album.Entry = a
			}
			album.Tracks = append(album.Tracks, a)
		}
	}

	// Try and get a pic
	album.Image, album.HasImage = r.GetAlbumImage(id)

	sort.Sort(ByTrackNumber(album.Tracks))
	return album
}

// Try and get a pic
func (r *Client) GetAlbumImage(id int) (image string, hasImage bool) {
	// Get the first image in the dir if there is one
	location := r.Db.Entries[id].Location
	strId := strconv.Itoa(id)
	imagePath := "/albums/a" + strId + ".jpg"

	// Check if already exists
	if _, err := os.Stat("public" + imagePath); err == nil {
		// File already exists
		image = imagePath
		hasImage = true
		return
	}

	// Remove the bits we dont want
	location = strings.TrimLeft(location, "file:/")
	lastSlash := strings.LastIndex(location, "/")
	location = location[:lastSlash+1]
	// fmt.Println(html.UnescapeString(location))

	e, _ := url.QueryUnescape(location)

	filepath.Walk("/"+e, func(path string, _ os.FileInfo, _ error) error {

		lastFour := path[len(path)-4:]
		if lastFour == ".jpg" || lastFour == "jpeg" || lastFour == ".png" {
			Copy("public/albums/a"+strId+".jpg", path)
			image = imagePath
			hasImage = true
			return nil
		}
		return nil
	})

	return
}

func Copy(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

func (r *Client) GetArtistsAlbums(id int) []Item {
	artist := r.Db.Entries[id].Artist
	albums := Albums{}

	// Get the first
	for _, a := range r.Albums {

		if a.Entry.Artist == artist {
			albums.Items = append(albums.Items, a)
		}

	}

	return albums.Items
}

func (r *Client) GetGenreTracks(id int) Item {
	genre := r.Db.Entries[id].Genre
	album := Item{Name: genre}

	// Get the first
	for _, e := range r.Db.Entries {

		if e.Genre == genre {
			album.Tracks = append(album.Tracks, e)
		}

	}

	sort.Sort(ByArtistE(album.Tracks))

	return album
}

func (r *Client) GetArtist(id int) Entry {
	return r.Db.Entries[id]
}

func (r *Client) PlayAlbum(id int) {
	r.ClearQueue()
	r.EnqueueAlbum(id)
	r.Play()
}

func (r *Client) PlayAlbumRandomly(id int) {
	a := r.GetAlbum(id)

	// Sort tracks randomly
	sort.Sort(ByRandom(a.Tracks))

	r.ClearQueue()
	for _, e := range a.Tracks {
		r.Enqueue(e.Location)
		fmt.Println(e.Title)
	}
	r.Play()
}

func (r *Client) EnqueueAlbum(id int) {

	a := r.GetAlbum(id)

	// Sort tracks by tracknumber
	sort.Sort(ByTrackNumber(a.Tracks))

	for _, e := range a.Tracks {
		r.Enqueue(e.Location)
	}
}

func (r *Client) PlayTrack(id int) {
	r.ClearQueue()
	r.Enqueue(r.Db.Entries[id].Location)
	r.Play()
}

// Assume that we are running from the users account which has Rhythmbox
// installed - therefore as long as the version of Rhythmbox is recentish, the
// lib should be located in .local/share/....
func (r *Client) GuessLibrary() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	r.Library = strings.Replace(RhythmboxXmlLibrary, "$HOME", usr.HomeDir, 1)
	fmt.Println(r.Library)
}

// Executes the options against the actual client
func (r *Client) Execute(s ...string) {

	out := r.ExecuteAndReturn(s...)

	print(out)
}

// Executes the options against the actual client
func (r *Client) ExecuteAndReturn(s ...string) string {

	cmd := exec.Command(RhythmboxClient, s...) //s[0], s[1]) //"--enqueue", "file:///home/ae/Music/Doolittle%20%5BMFSL%5D/Pixies%20-%20Doolittle%20(MFSL)%20-%2002%20-%20Tame.flac")
	out, err := cmd.Output()
	// fmt.Println(s)
	// fmt.Println(out)
	if err != nil {
		return string(err.Error())
	}

	return string(out)
}

var randSeed int64 = 1

func random(min, max int) int {
	randSeed++
	rand.Seed(time.Now().Unix() + randSeed)
	return rand.Intn(max-min) + min
}

func RandBool() bool {
	a := random(1, 3)
	return a == 1
}
