package exene

import (
	"fmt"
)

type Text struct {
    bounds Bounds
	Id string 
	Text string
	webIfc *WebInterface
}


func NewText(bounds Bounds, text string) *Text {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	textWidget := &Text{bounds, strId, text, nil}
	return textWidget
}


func (w *Text) Realize(webIfc *WebInterface, size Size, resizeChan chan Size) Html {
	w.webIfc = webIfc
	rSize := ClampBounds(w.bounds, size)
	go func() {
		for {
			select {
			case newSize := <- resizeChan:
				rSize := ClampBounds(w.bounds, newSize)
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
		},
		w.Text,
		nil,
		nil,
	}
}

func (w *Text) BoundsOf() Bounds {
	return w.bounds
}

func (w *Text) UpdateText(text string) {
	w.Text = text
	w.webIfc.updateChan <- map[string]any{"target": w.Id, "type": "update-text", "text": text}
}
