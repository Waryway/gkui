package ui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"image/color"
	"log"
	"os"
	"strconv"
	"time"
)

var SomString = "something"

func InitUi() {
	go func() {
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	i := 0

	// Update the board 3 times per second.
	advanceBoard := time.NewTicker(time.Second / 3)
	defer advanceBoard.Stop()

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				SomString = "s" + strconv.Itoa(i)
				title := material.H1(th, SomString)
				maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
				title.Color = maroon
				title.Alignment = text.Middle
				title.Layout(gtx)

				e.Frame(gtx.Ops)
			}
		case <-advanceBoard.C:
			i = i + 1
			w.Invalidate()
		}
	}
}
