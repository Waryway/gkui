package gui

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/markdown"
	"gioui.org/x/richtext"
	"github.com/inkeliz/giohyperlink"
	"gkui/env"
	"golang.org/x/net/context"
	"log"
	"sync"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func InitUi(wg *sync.WaitGroup, ctx context.Context, cancelFunc context.CancelFunc) {
	th := NewTheme(gofont.Collection())
	ui := UI{
		Window:   app.NewWindow(),
		Renderer: markdown.NewRenderer(),
		Theme:    th,
		Resize:   component.Resize{Ratio: 0.5},
		Ctx:      ctx,
		Cancel:   cancelFunc,
	}
	//gui.Renderer.Config.MonospaceFont.Typeface = "Go Mono"
	go func() {
		defer wg.Done()
		if err := ui.Loop(); err != nil {
			log.Print(err)
		}
	}()

	app.Main()
}

// UI specifies the user interface.
type UI struct {
	// External systems.
	// Window provides access to the OS window.
	Window *app.Window
	// Theme contains semantic style data. Extends `material.Theme`.
	Theme *Theme
	// Renderer transforms raw text containing markdown into rich text.
	Renderer *markdown.Renderer

	// Core state.
	// Editor retains raw text in an edit buffer.
	Editor widget.Editor
	// TextState retains rich text interactions: clicks, hovers and long presses.
	TextState richtext.InteractiveText
	// Resize state retains the split between the editor and the rendered text.
	component.Resize

	Ctx context.Context

	Cancel context.CancelFunc
}

func (ui *UI) Loop() error {
	var ops op.Ops

	for {
		e := ui.Window.NextEvent()
		giohyperlink.ListenEvents(e)
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ui.Layout(gtx) // Render the layout each time.
			e.Frame(gtx.Ops)
		}
	}
}

// Update processes events from the previous frame, updating state accordingly.
func (ui *UI) Update(gtx C) {
	go func() {
		for {
			select {
			case logMsg := <-*env.Logger.OutChan:
				ui.Editor.SetText(ui.Editor.Text() + "\r\n" + fmt.Sprint(logMsg))
				break
			}
		}
	}()

	for {
		o, event, ok := ui.TextState.Update(gtx)
		if !ok {
			break
		}
		switch event.Type {
		case richtext.Click:
			if url, ok := o.Get(markdown.MetadataURL).(string); ok && url != "" {
				if err := giohyperlink.Open(url); err != nil {
					// TODO(jfm): display UI element explaining the error to the user.
					env.Logger.WarnLog("error: opening hyperlink: %v", err)
				}
			}
		case richtext.Hover:
		case richtext.LongPress:
			env.Logger.DebugLog("long press")
			if url, ok := o.Get(markdown.MetadataURL).(string); ok && url != "" {
				ui.Window.Option(app.Title(url))
			}
		default:
			env.Logger.WarnLog("error: unhandled event: %v", event.Type)
		}
	}
	for {
		event, ok := ui.Editor.Update(gtx)
		if !ok {
			break
		}
		if _, ok := event.(widget.ChangeEvent); ok {
			var err error
			ui.Theme.cache, err = ui.Renderer.Render([]byte(ui.Editor.Text()))
			if err != nil {
				// TODO(jfm): display UI element explaining the error to the user.
				env.Logger.WarnLog("error: rendering markdown: %v", err)
			}
		}
	}
}

// Layout renders the current frame.
func (ui *UI) Layout(gtx C) D {
	ui.Update(gtx)
	widgets := []layout.Widget{
		func(gtx C) D {
			gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(200))
			return material.Editor(ui.Theme.Base, &ui.Editor, "Hint").Layout(gtx)
		},
	}
	list := &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	return material.List(ui.Theme.Base, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}
