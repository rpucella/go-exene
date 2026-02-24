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
	frame := &Frame{id, thick, color.String(), widget, nil}
	return frame
}

func (w *Frame) Realize(win Window, size Size, env Environment) *Html {
	if w.win != nil {
		return nil
	}
	w.win = win
	subResizeChan := make(chan Size)
	rSize := ClampBounds(w.BoundsOf(), size)
	subSize := AddSize(rSize, Size{-2 * w.thick, -2 * w.thick})
	subEnv := Environment{subResizeChan, nil, nil}
	subHtml := w.widget.Realize(win, subSize, subEnv)
	go func() {
		for {
			select {
			case newSize := <- env.ResizeChan:
				rSize := ClampBounds(w.BoundsOf(), newSize)
				subSize := AddSize(rSize, Size{-2 * w.thick, -2 * w.thick})
				subResizeChan <- subSize
				win.UpdateSize(w.id, rSize)
			}
		}
	}()
	return NewHtml("div").
		Id(w.id.String()).
		Styles(DefaultStyle(rSize)).
		Style("border", fmt.Sprintf("%dpx solid %s", w.thick, w.color)).
		Append(subHtml)
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
