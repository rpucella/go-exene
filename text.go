package exene

import (
	"fmt"
)

type Text struct {
    bounds Bounds
	id WId
	text string
	style *Style
	win Window
}


func NewText(bounds Bounds, text string, styles ...StyleOption) *Text {
	id := NewId()
	style := &Style{}
	for _, s := range styles {
		s(style)
	}
	textWidget := &Text{bounds, id, text, style, nil}
	return textWidget
}


func (w *Text) Realize(win Window, size Size, resizeChan chan Size) Html {
	w.win = win
	rSize := ClampBounds(w.bounds, size)
	go func() {
		for {
			select {
			case newSize := <- resizeChan:
				rSize := ClampBounds(w.bounds, newSize)
				win.UpdateSize(w.id, rSize)
			}
		}
	}()
	styleMap := w.style.mapOf()
	styleMap["height"] = fmt.Sprintf("%dpx", rSize.Height)
	styleMap["width"] = fmt.Sprintf("%dpx", rSize.Width)
	styleMap["overflow"] = "hidden"
	styleMap["transition"] = "height 0.1s, width 0.1s"
	styleMap["boxSizing"] = "border-box"
	return Html{
		w.id.String(),
		"div",
		nil,
		styleMap,
		w.text,
		nil,
		nil,
	}
}

func (w *Text) BoundsOf() Bounds {
	return w.bounds
}

func (w *Text) UpdateText(text string) {
	w.text = text
	w.win.UpdateText(w.id, text)
}
