package exene

import (
	"fmt"
)

type Padding struct {
	id WId
	thick int
	widget Widget
	win Window
}

func NewPadding(thick int, widget Widget) *Padding {
	id := NewId()
	padding := &Padding{id, thick, widget, nil}
	return padding
}

func (w *Padding) Realize(win Window, size Size, resizeChan chan Size) Html {
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
			"padding": fmt.Sprintf("%dpx", w.thick),
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

func (w *Padding) BoundsOf() Bounds {
	bounds := w.widget.BoundsOf()
	dim2 := FixDim(2 * w.thick)
	newBounds := Bounds{
		AddDim(bounds.Width, dim2),
		AddDim(bounds.Height, dim2),
	}
	return newBounds
}
