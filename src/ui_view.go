package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
)

type Drawable interface {
	Render(v *UIView)
	Init(v *UIView)
	Destroy(v *UIView)
}

type DrawableText struct {
	text      *text.Text
	transform pixel.Matrix
}

type UIView struct {
	elements  []Drawable
	textQueue []DrawableText

	imd       *imdraw.IMDraw
	evHandler *EventHandler
}

func (v *UIView) AddElement(element Drawable) {
	v.elements = append(v.elements, element)
}

func (v *UIView) enqueueText(text *text.Text, transform pixel.Matrix) {
	v.textQueue = append(v.textQueue, DrawableText{
		text:      text,
		transform: transform,
	})
}

func (v *UIView) activate() {
	for _, elem := range v.elements {
		elem.Init(v)
	}
}

func (v *UIView) deactivate() {
	for _, elem := range v.elements {
		elem.Destroy(v)
	}
}

func (v *UIView) render() {
	for _, elem := range v.elements {
		elem.Render(v)
	}
	v.imd.Draw(v.evHandler.win)

	for _, txtEl := range v.textQueue {
		txtEl.text.Draw(v.evHandler.win, txtEl.transform)
	}
}

func NewUIView(imd *imdraw.IMDraw, evHandler *EventHandler) *UIView {
	view := &UIView{
		elements:  make([]Drawable, 0, 10),
		textQueue: make([]DrawableText, 0, 10),
		imd:       imd,
		evHandler: evHandler,
	}

	return view
}
