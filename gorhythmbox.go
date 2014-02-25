package main

import (
	"fmt"

	"github.com/ae0000/gorhythmbox/rhythmbox"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
)

type PageData struct {
	Name     string
	Albums   []string
	Album    rhythmbox.Album
	PageType string
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

	m.Get("/next", func(r render.Render) {
		p := PageData{Name: "Next"}
		rb.Play()
		rb.Next()
		r.HTML(200, "home", p)
	})

	m.Get("/previous", func(r render.Render) {
		p := PageData{Name: "Previous"}
		rb.Play()
		rb.Previous()
		r.HTML(200, "home", p)
	})

	m.Get("/pause", func(r render.Render) {
		p := PageData{Name: "pause"}
		rb.PlayPause()
		r.HTML(200, "home", p)
	})

	m.Get("/albums", func(r render.Render) {
		p := PageData{
			Name:     "Albums",
			PageType: "albums",
			Albums:   rb.GetAlbums(),
		}
		r.HTML(200, "albums", p)
	})

	m.Get("/artists", func(r render.Render) {
		p := PageData{Name: "Artists", Albums: rb.GetArtists()}
		r.HTML(200, "albums", p)
	})

	m.Get("/genres", func(r render.Render) {
		p := PageData{Name: "Albums", Albums: rb.GetAlbums()}
		r.HTML(200, "albums", p)
	})

	m.Get("/albums/:album", func(r render.Render, params martini.Params) {
		album := params["album"]
		p := PageData{Name: "Albums", Albums: rb.GetAlbums(), Album: rb.GetAlbum(album)}

		fmt.Println(album)
		r.HTML(200, "album", p)
	})

	m.Get("/albums/play/:album", func(r render.Render, params martini.Params) {
		album := params["album"]
		p := PageData{Name: "Albums", Albums: rb.GetAlbums(), Album: rb.GetAlbum(album)}
		rb.PlayAlbum(album)

		r.HTML(200, "album", p)
	})

	m.Get("/album/:album/track/:track", func(r render.Render, params martini.Params) {
		album := params["album"]
		track := params["track"]
		p := PageData{Name: "Albums", Albums: rb.GetAlbums(), Album: rb.GetAlbum(album)}
		rb.PlayTrack(album, track)

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
