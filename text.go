package exene

import (
	"fmt"
)

type Text struct {
    bounds Bounds
	Id string 
	Text string
	style *Style
	webIfc *WebInterface
}


func NewText(bounds Bounds, text string, styles ...StyleOption) *Text {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	style := &Style{}
	for _, s := range styles {
		s(style)
	}
	textWidget := &Text{bounds, strId, text, style, nil}
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
	styleMap := w.style.mapOf()
	styleMap["height"] = fmt.Sprintf("%dpx", rSize.Height)
	styleMap["width"] = fmt.Sprintf("%dpx", rSize.Width)
	styleMap["overflow"] = "hidden"
	styleMap["transition"] = "height 0.1s, width 0.1s"
	styleMap["boxSizing"] = "border-box"
	return Html{
		w.Id,
		"div",
		nil,
		styleMap,
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
