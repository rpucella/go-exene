package exene

import (
)

type Label struct {
    bounds Bounds
	id WId
	text string
	style *Style
	// Realized info
	win Window
	labelChan chan string
}


func NewLabel(bounds Bounds, text string, styles ...StyleOption) *Label {
	id := NewId()
	style := &Style{}
	for _, s := range styles {
		s(style)
	}
	textWidget := &Label{bounds, id, text, style, nil, nil}
	return textWidget
}


func (w *Label) Realize(win Window, size Size, resizeChan chan Size) Html {
	if w.win != nil {
		return Html{}
	}
	w.win = win
	labelChan := make(chan string)	
	w.labelChan = labelChan
	rSize := ClampBounds(w.bounds, size)
	go func() {
		for {
			select {
			case newSize := <- resizeChan:
				rSize := ClampBounds(w.bounds, newSize)
				win.UpdateSize(w.id, rSize)

			case newText := <- labelChan:
				w.win.UpdateText(w.id, newText)
			}
		}
	}()
	styling := CreateDefaultStyle(rSize)
	w.style.ExtendStyle(styling)
	return Html{
		w.id.String(),
		"div",
		nil,
		styling,
		w.text,
		nil,
		nil,
	}
}

func (w *Label) BoundsOf() Bounds {
	return w.bounds
}

func (w *Label) UpdateText(text string) {
	w.labelChan <- text
}
