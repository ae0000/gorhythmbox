package main

import (
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
			PageType: "artists",
			Albums:   rb.GetArtists(),
		}
		r.HTML(200, "albums", p)
	})

	m.Get("/genres", func(r render.Render) {
		p := PageData{Name: "Albums", Albums: rb.GetAlbums()}
		r.HTML(200, "albums", p)
	})

	m.Get("/albums/:album", func(r render.Render, params martini.Params) {
		album := params["album"]
		p := PageData{Name: "Album", Album: rb.GetAlbum(album), PageId: album}

		r.HTML(200, "album", p)
	})

	m.Get("/albums/play/:album", func(r render.Render, params martini.Params) {
		album := params["album"]
		p := PageData{Name: "Albums", Albums: rb.GetAlbums(), Album: rb.GetAlbum(album)}
		rb.PlayAlbum(album)

		r.HTML(200, "album", p)
	})

	m.Get("/album/:albumid/track/:trackid", func(r render.Render, params martini.Params) {
		albumid := params["albumid"]
		trackid := params["trackid"]

		// Need to convert trackid to Int
		id, _ := strconv.ParseInt(trackid, 10, 0)
		idi := int(id)

		album := rb.GetAlbum(albumid)
		album.SelectTrack(idi)

		p := PageData{
			Name:   "Albums",
			Album:  album,
			PageId: albumid,
		}

		rb.PlayTrack(idi)

		r.HTML(200, "album", p)
	})

	m.Get("/albums/enqueue/:album", func(r render.Render, params martini.Params) {
		album := params["album"]
		p := PageData{Name: "Albums", Albums: rb.GetAlbums(), Album: rb.GetAlbum(album)}
		rb.EnqueueAlbum(album)

		r.HTML(200, "album", p)
	})

	m.Run()

}
