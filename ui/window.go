package ui

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"
	"image/color"
	"sync"

	//"gioui.org/layout"
	"gioui.org/op"
	//"gioui.org/text"
	"gioui.org/widget/material"
	//"image/color"
	"log"
	"strconv"
	"time"
)

var SomString = "something"

func InitUi(wg *sync.WaitGroup) {

	go func() {
		defer wg.Done()
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Print(err)
		}
	}()

	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops
	i := 0

	temp := Button{
		Name:   "test",
		Width:  100,
		Height: 20,
		Theme:  th,
	}

	// Update the board 3 times per second.
	advanceBoard := time.NewTicker(time.Second / 3)
	defer advanceBoard.Stop()

	go func() {
		<-advanceBoard.C
		i = i + 1
		w.Invalidate()
	}()

	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			SomString = "s" + strconv.Itoa(i)
			title := material.H1(th, SomString)
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon
			title.Alignment = text.Middle
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(title.Layout),
				layout.Rigid(temp.Init().Render().Button.Layout),
			)

			e.Frame(gtx.Ops)

		}
	}
}
