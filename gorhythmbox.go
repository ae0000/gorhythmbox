package main

import (
	"fmt"
	"strconv"

	"github.com/ae0000/gorhythmbox/rhythmbox"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
)

type PageData struct {
	PageId    string
	Name      string
	Albums    []rhythmbox.Item
	Album     rhythmbox.Item
	Artist    rhythmbox.Entry
	PageType  string
	Selected  int
	ShowTitle bool
}

type AjaxReturn struct {
	A,
	B,
	C,
	D,
	E string
}

func main() {
	fmt.Println("***************************************************** [START]")
	// Setup Rhythmbox
	rb := rhythmbox.Client{}
	rb.GuessLibrary()
	rb.Setup()

	// Setup martini
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
	}))

	m.Get("/", func(r render.Render) {
		p := PageData{Name: "Home"}
		r.HTML(200, "home", p)
	})

	m.Get("/ajax/:do", func(r render.Render, params martini.Params) {
		switch params["do"] {
		case "previous":
			rb.Previous()
		case "play":
			rb.Play()
		case "pause":
			rb.Pause()
		case "next":
			rb.Play()
			rb.Next()
		case "volumeup":
			rb.VolumeUp()
		case "volumedown":
			rb.VolumeDown()
		case "current":
			r.JSON(200, AjaxReturn{A: rb.PrintPlayingFormat("<strong>%" + "aa:</strong><em> " + "%" + "tt</em>")})
			return
		}

		r.JSON(200, PageData{Name: "Next"}) //  HTML(200, "home", p)
	})

	m.Get("/albums", func(r render.Render) {
		p := PageData{
			Name:      "Albums",
			PageType:  "albums",
			Albums:    rb.GetAlbums(),
			ShowTitle: true,
		}
		r.HTML(200, "albums", p)
	})

	m.Get("/artists", func(r render.Render) {
		p := PageData{
			Name:     "Artists",
			PageType: "artist",
			Albums:   rb.GetArtists(),
		}
		r.HTML(200, "artists", p)
	})

	m.Get("/genres", func(r render.Render) {
		p := PageData{
			Name:     "Genres",
			PageType: "genre",
			Albums:   rb.GetGenres(),
		}
		r.HTML(200, "genres", p)
	})

	m.Get("/albums/:albumid", func(r render.Render, params martini.Params) {
		albumid := params["albumid"]

		// Need to convert to Int
		aid, _ := strconv.ParseInt(albumid, 10, 0)
		adi := int(aid)

		p := PageData{Name: "Album", Album: rb.GetAlbum(adi), PageId: albumid}

		r.HTML(200, "album", p)
	})

	m.Get("/album/:albumid/track/:trackid", func(r render.Render, params martini.Params) {
		albumid := params["albumid"]
		trackid := params["trackid"]

		// Need to convert to Int
		aid, _ := strconv.ParseInt(albumid, 10, 0)
		adi := int(aid)
		id, _ := strconv.ParseInt(trackid, 10, 0)
		idi := int(id)

		album := rb.GetAlbum(adi)
		album.SelectTrack(idi)

		p := PageData{
			Name:   "Albums",
			Album:  album,
			PageId: albumid,
		}

		rb.PlayTrack(idi)

		r.HTML(200, "album", p)
	})

	m.Get("/album/enqueue/:albumid", func(r render.Render, params martini.Params) {
		albumid := params["albumid"]

		// Need to convert artistid to Int
		id, _ := strconv.ParseInt(albumid, 10, 0)
		idi := int(id)

		p := PageData{Name: "Album", Album: rb.GetAlbum(idi), PageId: albumid}

		rb.EnqueueAlbum(idi)
		r.HTML(200, "album", p)
	})

	m.Get("/album/play/:albumid", func(r render.Render, params martini.Params) {
		albumid := params["albumid"]

		// Need to convert artistid to Int
		id, _ := strconv.ParseInt(albumid, 10, 0)
		idi := int(id)

		p := PageData{Name: "Album", Album: rb.GetAlbum(idi), PageId: albumid}

		rb.PlayAlbum(idi)
		r.HTML(200, "album", p)
	})

	m.Get("/album/random/:albumid", func(r render.Render, params martini.Params) {
		albumid := params["albumid"]

		// Need to convert artistid to Int
		id, _ := strconv.ParseInt(albumid, 10, 0)
		idi := int(id)

		p := PageData{Name: "Album", Album: rb.GetAlbum(idi), PageId: albumid}

		rb.PlayAlbumRandomly(idi)
		r.HTML(200, "album", p)
	})

	m.Get("/artist/:artistid", func(r render.Render, params martini.Params) {
		artistid := params["artistid"]

		// Need to convert artistid to Int
		id, _ := strconv.ParseInt(artistid, 10, 0)
		idi := int(id)

		p := PageData{
			Name:   "Album",
			Albums: rb.GetArtistsAlbums(idi),
			PageId: artistid,
			Artist: rb.GetArtist(idi),
		}

		r.HTML(200, "artist", p)
	})

	m.Get("/artist/enqueue/:artistid", func(r render.Render, params martini.Params) {
		artistid := params["artistid"]

		// Need to convert artistid to Int
		id, _ := strconv.ParseInt(artistid, 10, 0)
		idi := int(id)

		p := PageData{
			Name:   "Album",
			Albums: rb.GetArtistsAlbums(idi),
			PageId: artistid,
			Artist: rb.GetArtist(idi),
		}

		rb.EnqueueArtist(idi)
		r.HTML(200, "artist", p)
	})

	m.Get("/artist/play/:artistid", func(r render.Render, params martini.Params) {
		artistid := params["artistid"]

		// Need to convert artistid to Int
		id, _ := strconv.ParseInt(artistid, 10, 0)
		idi := int(id)

		p := PageData{
			Name:   "Album",
			Albums: rb.GetArtistsAlbums(idi),
			PageId: artistid,
			Artist: rb.GetArtist(idi),
		}

		rb.PlayArtist(idi)
		r.HTML(200, "artist", p)
	})

	m.Get("/artist/random/:artistid", func(r render.Render, params martini.Params) {
		artistid := params["artistid"]

		// Need to convert artistid to Int
		id, _ := strconv.ParseInt(artistid, 10, 0)
		idi := int(id)

		p := PageData{
			Name:   "Album",
			Albums: rb.GetArtistsAlbums(idi),
			PageId: artistid,
			Artist: rb.GetArtist(idi),
		}

		rb.PlayArtistRandomly(idi)
		r.HTML(200, "artist", p)
	})

	m.Get("/genre/:genreid", func(r render.Render, params martini.Params) {
		genreid := params["genreid"]

		// Need to convert genreid to Int
		id, _ := strconv.ParseInt(genreid, 10, 0)
		idi := int(id)

		p := PageData{
			Name:   "Genre",
			Album:  rb.GetGenreTracks(idi),
			PageId: genreid,
		}

		r.HTML(200, "genre", p)
	})

	m.Get("/genre/:genreid/track/:trackid", func(r render.Render, params martini.Params) {
		genreid := params["genreid"]
		trackid := params["trackid"]

		// Need to convert to Int
		gid, _ := strconv.ParseInt(genreid, 10, 0)
		gidi := int(gid)
		id, _ := strconv.ParseInt(trackid, 10, 0)
		idi := int(id)

		album := rb.GetGenreTracks(gidi)
		album.SelectTrack(idi)

		p := PageData{
			Name:   "Genre",
			Album:  album,
			PageId: genreid,
		}

		rb.PlayTrack(idi)

		r.HTML(200, "genre", p)
	})

	m.Get("/genre/enqueue/:genreid", func(r render.Render, params martini.Params) {
		genreid := params["genreid"]

		// Need to convert genreid to Int
		id, _ := strconv.ParseInt(genreid, 10, 0)
		idi := int(id)

		p := PageData{
			Name:   "Genre",
			Album:  rb.GetGenreTracks(idi),
			PageId: genreid,
		}

		rb.EnqueueGenre(idi)
		r.HTML(200, "genre", p)
	})

	m.Get("/genre/play/:genreid", func(r render.Render, params martini.Params) {
		genreid := params["genreid"]

		// Need to convert genreid to Int
		id, _ := strconv.ParseInt(genreid, 10, 0)
		idi := int(id)

		p := PageData{
			Name:   "Genre",
			Album:  rb.GetGenreTracks(idi),
			PageId: genreid,
		}

		rb.PlayGenre(idi)
		r.HTML(200, "genre", p)
	})

	m.Get("/genre/random/:genreid", func(r render.Render, params martini.Params) {
		genreid := params["genreid"]

		// Need to convert genreid to Int
		id, _ := strconv.ParseInt(genreid, 10, 0)
		idi := int(id)

		p := PageData{
			Name:   "Genre",
			Album:  rb.GetGenreTracks(idi),
			PageId: genreid,
		}

		rb.PlayGenreRandomly(idi)
		r.HTML(200, "genre", p)
	})

	m.Run()

}
