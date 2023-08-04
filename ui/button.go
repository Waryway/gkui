package ui

import (
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Button struct {
	Name   string
	Hint   string
	Event  *widget.Clickable
	Height int
	Width  int
	Theme  *material.Theme
	Button *material.ButtonStyle
}

func (b *Button) Init() *Button {
	b.Event = new(widget.Clickable)
	return b
}

func (b *Button) Render() *Button {
	temp := material.Button(b.Theme, b.Event, b.Name)
	b.Button = &temp
	return b
}
