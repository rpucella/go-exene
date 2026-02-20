package exene

import (
	"fmt"
)

type Button struct {
    bounds Bounds
	id WId
	label string
	action func()
	win Window
}


func NewButton(bounds Bounds, label string, act func()) *Button {
	id := NewId()
	button := &Button{bounds, id, label, act, nil}
	return button
}


func (w *Button) Realize(win Window, size Size, resizeChan chan Size) Html {
	w.win = win
	rSize := ClampBounds(w.bounds, size)
	eventChan := make(chan bool)
	win.RegisterEventChan(w.id, eventChan)
	go func() {
		for {
			// Also: handle destroy messages?
			select {
			case <- eventChan:
				w.action()

			case newSize := <- resizeChan:
				rSize := ClampBounds(w.bounds, newSize)
				win.UpdateSize(w.id, rSize)
			}
		}
	}()
	return Html{
		w.id.String(),
		"button",
		nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", rSize.Height),
			"width": fmt.Sprintf("%dpx", rSize.Width),
			"overflow": "hidden",
			"border": "none",
			"transition": "height 0.1s, width 0.1s",
			"boxSizing": "border-box",
			"cursor": "pointer",
		},
		w.label,
		nil,
		[]string{"click"},
	}
}

func (w *Button) BoundsOf() Bounds {
	return w.bounds
}
