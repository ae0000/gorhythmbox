package rhythmbox

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"os/user"
	"sort"
	"strings"
)

type Client struct {
	Library string
	Db      Rhythmdb
	Artists []string
	Albums  []string
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
}

type Album struct {
	Entry  Entry
	Tracks []Entry
}

// ByTrackNumber implements sort.Interface for []Entry based on the track number
type ByTrackNumber []Entry

func (a ByTrackNumber) Len() int           { return len(a) }
func (a ByTrackNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTrackNumber) Less(i, j int) bool { return a[i].TrackNumber < a[j].TrackNumber }

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

	// Sort out the unique artists, albums and genres
	for _, e := range r.Db.Entries {
		if len(e.Album) > 0 {
			if !r.AlbumExists(e.Album) {
				r.Albums = append(r.Albums, e.Album)
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
}

func (r *Client) AlbumExists(s string) bool {
	for _, a := range r.Albums {
		if a == s {
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

func (r *Client) GetAlbums() []string {
	return r.Albums
	// for _, a := range r.Albums {
	// 	fmt.Println(a)
	// }
}

func (r *Client) GetAlbum(searchAlbum string) Album {
	album := Album{}

	for _, a := range r.Db.Entries {
		if a.Album == searchAlbum {
			if len(album.Entry.Album) == 0 {
				album.Entry = a
			}
			album.Tracks = append(album.Tracks, a)
		}
	}

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

func (r *Client) PlayTrack(album, track string) {
	for _, e := range r.Db.Entries {
		if e.Album == album && e.Title == track {
			r.ClearQueue()
			r.Enqueue(e.Location)
			r.Play()
		}
	}
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

	cmd := exec.Command(RhythmboxClient, s...) //s[0], s[1]) //"--enqueue", "file:///home/ae/Music/Doolittle%20%5BMFSL%5D/Pixies%20-%20Doolittle%20(MFSL)%20-%2002%20-%20Tame.flac")
	out, err := cmd.Output()
	// fmt.Println(s)
	// fmt.Println(out)
	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))

}