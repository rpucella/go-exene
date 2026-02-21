package exene

import (
	"fmt"
)

type Background struct {
	id WId
	bgColor string
	fgColor string
	widget Widget
	win Window
}

func NewBackground(bgColor Color, fgColor Color, widget Widget) *Background {
	id := NewId()
	frame := &Background{id, bgColor.String(), fgColor.String(), widget, nil}
	return frame
}

func (w *Background) Realize(win Window, size Size, resizeChan chan Size) Html {
	if w.win != nil {
		return Html{}
	}
	w.win = win
	subResizeChan := make(chan Size)
	rSize := ClampBounds(w.BoundsOf(), size)
	subHtml := w.widget.Realize(win, rSize, subResizeChan)
	go func() {
		for {
			select {
			case newSize := <- resizeChan:
				rSize := ClampBounds(w.BoundsOf(), newSize)
				subResizeChan <- rSize
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
			"backgroundColor": w.bgColor,
			"color": w.fgColor,
			"overflow": "hidden",
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

func (w *Background) BoundsOf() Bounds {
	bounds := w.widget.BoundsOf()
	return bounds
}
