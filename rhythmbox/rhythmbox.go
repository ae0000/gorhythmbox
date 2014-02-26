package rhythmbox

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"os/user"
	"sort"
	"strconv"
	"strings"
)

type Client struct {
	Library string
	Db      Rhythmdb
	Artists []string
	Albums  []Item
	Genres  []string
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
	Id     int
	Name   string
	Type   string
	Entry  Entry
	Tracks []Entry
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
type ByArtist []Item
type ByAlbum []Item

func (a ByTrackNumber) Len() int           { return len(a) }
func (a ByTrackNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTrackNumber) Less(i, j int) bool { return a[i].TrackNumber < a[j].TrackNumber }

func (a ByArtist) Len() int           { return len(a) }
func (a ByArtist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByArtist) Less(i, j int) bool { return a[i].Entry.Artist < a[j].Entry.Artist }

func (a ByAlbum) Len() int           { return len(a) }
func (a ByAlbum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAlbum) Less(i, j int) bool { return a[i].Entry.Album < a[j].Entry.Album }

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
				r.Albums = append(r.Albums, item)
			}
		}
		if len(e.Artist) > 0 {
			if !r.ArtistExists(e.Artist) {
				r.Artists = append(r.Artists, e.Artist)
			}
		}
		if len(e.Genre) > 0 {
			if !r.GenreExists(e.Genre) {
				r.Genres = append(r.Genres, e.Genre)
			}
		}
	}

	// Sort albums by artist

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
		if a == s {
			return true
		}
	}
	return false
}

func (r *Client) GenreExists(s string) bool {
	for _, a := range r.Genres {
		if a == s {
			return true
		}
	}
	return false
}

func (r *Client) GetAlbums() []Item {
	sort.Sort(ByArtist(r.Albums))
	return r.Albums
}

func (r *Client) GetAlbum(albumId string) Item {
	// Need to convert albumId to Int
	id, _ := strconv.ParseInt(albumId, 10, 0)
	idi := int(id)
	album := Item{}
	albumName := r.Db.Entries[idi].Album

	for _, a := range r.Db.Entries {
		if a.Album == albumName {
			if len(album.Entry.Album) == 0 {
				album.Entry = a
			}
			album.Tracks = append(album.Tracks, a)
		}
	}

	sort.Sort(ByTrackNumber(album.Tracks))
	return album
}

func (r *Client) GetArtists() []string {
	return r.Artists
}

func (r *Client) GetGenres() {
	for _, a := range r.Genres {
		fmt.Println(a)
	}
}

func (r *Client) PlayAlbum(album string) {
	r.ClearQueue()
	r.EnqueueAlbum(album)
}

func (r *Client) EnqueueAlbum(album string) {
	// Get all tracks that match this album
	t := Tracks{}

	for _, e := range r.Db.Entries {
		if e.Album == album {
			t.Entries = append(t.Entries, e)
		}
	}

	if len(t.Entries) == 0 {
		fmt.Println("No tracks found")
		return
	}

	// Sort tracks by tracknumber
	sort.Sort(ByTrackNumber(t.Entries))

	for _, e := range t.Entries {
		fmt.Println("Enqueue: ", e.Title)
		r.Enqueue(e.Location)
	}

}

func (r *Client) PlayTrack(id int) {
	r.ClearQueue()
	r.Enqueue(r.Db.Entries[id].Location)
	r.Next()
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
