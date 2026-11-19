package main

import (
	"log"

	"kapp/internal/app"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {

	a := app.New()
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:     "kapp",
			Width:     760,
			Height:    520,
			MinWidth:  600,
			MinHeight: 420,
			// Opens where the user left it last time.
			StateKey: "main",
			// The window shows the interface MyGo draws, not a web page.
			Content: ui.View(a.View),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
