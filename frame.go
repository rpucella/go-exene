package exene

import (
	"fmt"
)

type Frame struct {
	Id string
	thick int
	widget Widget
	webIfc *WebInterface
}

func NewFrame(thick int, widget Widget) *Frame {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	frame := &Frame{strId, thick, widget, nil}
	return frame
}

func (w *Frame) Realize(webIfc *WebInterface, size Size, resizeChan chan Size) Html {
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
			"border": fmt.Sprintf("%dpx solid black", w.thick),
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
