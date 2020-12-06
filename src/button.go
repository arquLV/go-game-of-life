package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type button struct {
	label      string
	position   pixel.Vec
	dimensions pixel.Vec
	onClick    func()

	textRenderer    *text.Text
	eventAreaHandle int
	isHovered       bool
}

func (b *button) getRect() pixel.Rect {
	return pixel.R(
		b.position.X,
		b.position.Y,
		b.position.X+b.dimensions.X,
		b.position.Y+b.dimensions.Y,
	)
}

func (b *button) Render(view *UIView) {
	imd := view.imd

	if b.isHovered {
		imd.Color = color.RGBA{0x00, 0xdd, 0x00, 0xff}
	} else {
		imd.Color = color.RGBA{0x00, 0xff, 0x00, 0xff}
	}
	imd.Push(b.position, b.position.Add(b.dimensions))
	imd.Rectangle(0)
}

func (b *button) Init(view *UIView) {
	btnRect := b.getRect()

	view.enqueueText(b.textRenderer, pixel.IM.Scaled(b.textRenderer.Orig, 2))

	evHandle := view.evHandler.AddEventArea(btnRect)
	b.eventAreaHandle = evHandle

	view.evHandler.OnClick(evHandle, b.onClick)

	view.evHandler.OnMouseEnter(evHandle, func() {
		b.isHovered = true
	})

	view.evHandler.OnMouseLeave(evHandle, func() {
		b.isHovered = false
	})
}

func (b *button) Destroy(view *UIView) {
	view.evHandler.RemoveAllHandlers(b.eventAreaHandle)
}

func NewButton(label string, position pixel.Vec, dimensions pixel.Vec, onClick func()) *button {
	btn := &button{
		label:      label,
		position:   position,
		dimensions: dimensions,
		onClick:    onClick,

		eventAreaHandle: -1,
		isHovered:       false,
	}

	// labelPosition := position.Add(pixel.V(dimensions.X ))
	fontAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	btn.textRenderer = text.New(position, fontAtlas)

	labelBounds := btn.textRenderer.BoundsOf(label)

	btn.textRenderer.Dot.X += (dimensions.X - labelBounds.W()) / 4
	btn.textRenderer.Dot.Y += (dimensions.Y - labelBounds.H()) / 4

	btn.textRenderer.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	fmt.Fprintln(btn.textRenderer, label)

	return btn
}
