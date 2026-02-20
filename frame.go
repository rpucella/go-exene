package exene

import (
	"fmt"
)

type Frame struct {
	id WId
	thick int
	color string
	widget Widget
	win Window
}

func NewFrame(thick int, color Color, widget Widget) *Frame {
	id := NewId()
	frame := &Frame{id, thick, string(color), widget, nil}
	return frame
}

func (w *Frame) Realize(win Window, size Size, resizeChan chan Size) Html {
	w.win = win
	subResizeChan := make(chan Size)
	rSize := ClampBounds(w.BoundsOf(), size)
	subSize := AddSize(rSize, Size{-2 * w.thick, -2 * w.thick})
	subHtml := w.widget.Realize(win, subSize, subResizeChan)
	go func() {
		for {
			select {
			case newSize := <- resizeChan:
				rSize := ClampBounds(w.BoundsOf(), newSize)
				subSize := AddSize(rSize, Size{-2 * w.thick, -2 * w.thick})
				subResizeChan <- subSize
				win.UpdateSize(w.id, rSize)
			}
		}
	}()
	return Html{
		w.id.String(),
		"div",
		nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", rSize.Height),
			"width": fmt.Sprintf("%dpx", rSize.Width),
			"overflow": "hidden",
			"border": fmt.Sprintf("%dpx solid %s", w.thick, w.color),
			"transition": "height 0.1s, width 0.1s",
			"boxSizing": "border-box",
		},
		"",
		[]Html{
			subHtml,
		},
		nil,
	}
}

func (w *Frame) BoundsOf() Bounds {
	bounds := w.widget.BoundsOf()
	dim2 := FixDim(2 * w.thick)
	newBounds := Bounds{
		AddDim(bounds.Width, dim2),
		AddDim(bounds.Height, dim2),
	}
	return newBounds
}
