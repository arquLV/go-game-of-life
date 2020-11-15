package main

import (
	"fmt"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type eventHandler struct {
	win *pixelgl.Window

	registeredAreas []pixel.Rect
	hoverHandlers   map[int]func()
	clickHandlers   map[int]func()
}

var (
	handlerInstantiator sync.Once
	handlerInstance     eventHandler
)

func (handler *eventHandler) start() {
	for !handler.win.Closed() {
		mousePos := handler.win.MousePosition()

		for areaIndex, area := range handler.registeredAreas {
			inArea := (mousePos.X >= area.Min.X) &&
				(mousePos.X <= area.Max.X) &&
				(mousePos.Y >= area.Min.Y) &&
				(mousePos.Y <= area.Max.Y)

			if inArea {
				onHover, hasHoverHandler := handler.hoverHandlers[areaIndex]
				if hasHoverHandler {
					onHover()
				}

				onClick, hasClickHandler := handler.clickHandlers[areaIndex]
				if hasClickHandler && handler.win.JustPressed(pixelgl.MouseButtonLeft) {
					onClick()
				}
			}
		}
	}
}

func (handler *eventHandler) OnHover(target pixel.Rect, callback func()) int {
	handler.registeredAreas = append(handler.registeredAreas, target)
	areaIndex := len(handler.registeredAreas) - 1

	handler.hoverHandlers[areaIndex] = callback
	return areaIndex
}

func (handler *eventHandler) OnClick(target pixel.Rect, callback func()) int {
	handler.registeredAreas = append(handler.registeredAreas, target)
	areaIndex := len(handler.registeredAreas) - 1
	fmt.Println("Adding click handler to area ", areaIndex)

	handler.clickHandlers[areaIndex] = callback
	return areaIndex
}

// NewEventHandler creates and starts new event handler if not done already
func NewEventHandler(window *pixelgl.Window) eventHandler {
	handlerInstantiator.Do(func() {
		fmt.Println("Instantiating")
		handlerInstance := eventHandler{win: window}

		handlerInstance.registeredAreas = make([]pixel.Rect, 0, 10)
		handlerInstance.hoverHandlers = make(map[int]func())
		handlerInstance.clickHandlers = make(map[int]func())

		go handlerInstance.start()
	})

	return handlerInstance
}
