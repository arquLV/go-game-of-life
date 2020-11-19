package main

import (
	"fmt"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type EventHandler struct {
	win *pixelgl.Window

	registeredAreas []pixel.Rect
	hoveredAreas    map[int]bool

	hoverHandlers      map[int]func()
	mouseEnterHandlers map[int]func()
	mouseLeaveHandlers map[int]func()
	clickHandlers      map[int]func()
}

var (
	handlerInstantiator sync.Once
	handlerInstance     EventHandler
)

func (handler *EventHandler) start() {
	for !handler.win.Closed() {
		mousePos := handler.win.MousePosition()

		for areaIndex, area := range handler.registeredAreas {
			inArea := (mousePos.X >= area.Min.X) &&
				(mousePos.X <= area.Max.X) &&
				(mousePos.Y >= area.Min.Y) &&
				(mousePos.Y <= area.Max.Y)

			if inArea {
				onMouseEnter, hasMouseEnterHandler := handler.mouseEnterHandlers[areaIndex]
				if hasMouseEnterHandler && handler.hoveredAreas[areaIndex] == false {
					onMouseEnter()
				}

				handler.hoveredAreas[areaIndex] = true

				onHover, hasHoverHandler := handler.hoverHandlers[areaIndex]
				if hasHoverHandler {
					onHover()
				}

				onClick, hasClickHandler := handler.clickHandlers[areaIndex]
				if hasClickHandler && handler.win.JustPressed(pixelgl.MouseButtonLeft) {
					onClick()
				}
			} else if !inArea && handler.hoveredAreas[areaIndex] == true {
				handler.hoveredAreas[areaIndex] = false
				onMouseLeave, hasMouseLeaveHandler := handler.mouseLeaveHandlers[areaIndex]
				if hasMouseLeaveHandler {
					onMouseLeave()
				}
			}
		}
	}
}

func (handler *EventHandler) AddEventArea(area pixel.Rect) int {
	handler.registeredAreas = append(handler.registeredAreas, area)
	areaIndex := len(handler.registeredAreas) - 1

	return areaIndex
}

func (handler *EventHandler) OnMouseEnter(target int, callback func()) bool {
	if len(handler.registeredAreas)-1 > target {
		return false
	}
	handler.mouseEnterHandlers[target] = callback
	return true
}

func (handler *EventHandler) OnMouseLeave(target int, callback func()) bool {
	if len(handler.registeredAreas)-1 > target {
		return false
	}
	handler.mouseLeaveHandlers[target] = callback
	return true
}

func (handler *EventHandler) OnHover(target int, callback func()) bool {
	if len(handler.registeredAreas)-1 > target {
		return false
	}
	handler.hoverHandlers[target] = callback
	return true
}

func (handler *EventHandler) OnClick(target int, callback func()) bool {
	if len(handler.registeredAreas)-1 > target {
		return false
	}
	handler.clickHandlers[target] = callback
	return true
}

func (handler *EventHandler) RemoveAllHandlers(target int) bool {
	if len(handler.registeredAreas)-1 > target {
		return false
	}

	delete(handler.mouseEnterHandlers, target)
	delete(handler.mouseLeaveHandlers, target)
	delete(handler.hoverHandlers, target)
	delete(handler.clickHandlers, target)

	return true
}

// NewEventHandler creates and starts new event handler if not done already
func NewEventHandler(window *pixelgl.Window) *EventHandler {
	handlerInstantiator.Do(func() {
		fmt.Println("Instantiating")
		handlerInstance = EventHandler{
			win:             window,
			registeredAreas: make([]pixel.Rect, 0, 10),
			hoveredAreas:    make(map[int]bool),

			mouseEnterHandlers: make(map[int]func()),
			mouseLeaveHandlers: make(map[int]func()),
			hoverHandlers:      make(map[int]func()),
			clickHandlers:      make(map[int]func()),
		}

		go handlerInstance.start()
	})

	return &handlerInstance
}
