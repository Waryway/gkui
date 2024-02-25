package gui

import (
	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/widget/material"
	"gioui.org/x/richtext"
)

// NewTheme instantiates a theme, extending material theme.
func NewTheme(font []font.FontFace) *Theme {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(font))
	return &Theme{
		Base: th,
	}
}

// Theme contains semantic style data.
type Theme struct {
	// Base theme to extend.
	Base *material.Theme
	// cache of processed markdown.
	cache []richtext.SpanStyle
}
