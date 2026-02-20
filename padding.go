package exene

import (
	"fmt"
)

type Padding struct {
	Id string
	thick int
	widget Widget
	webIfc *WebInterface
}

func NewPadding(thick int, widget Widget) *Padding {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	padding := &Padding{strId, thick, widget, nil}
	return padding
}

func (w *Padding) Realize(webIfc *WebInterface, size Size, resizeChan chan Size) Html {
	w.webIfc = webIfc
	subResizeChan := make(chan Size)
	rSize := ClampBounds(w.BoundsOf(), size)
	subSize := AddSize(rSize, Size{-2 * w.thick, -2 * w.thick})
	subHtml := w.widget.Realize(webIfc, subSize, subResizeChan)
	go func() {
		for {
			select {
			case newSize := <- resizeChan:
				rSize := ClampBounds(w.BoundsOf(), newSize)
				subSize := AddSize(rSize, Size{-2 * w.thick, -2 * w.thick})
				subResizeChan <- subSize
				webIfc.UpdateSize(w.Id, rSize)
			}
		}
	}()
	return Html{
		w.Id,
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
