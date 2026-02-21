package exene

import (
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
	if w.win != nil {
		return Html{}
	}
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
	styling := CreateDefaultStyle(rSize)
	styling["cursor"] = "pointer"
	return Html{
		w.id.String(),
		"button",
		nil,
		styling,
		w.label,
		nil,
		[]string{"click"},
	}
}

func (w *Button) BoundsOf() Bounds {
	return w.bounds
}
