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


func (w *Button) Realize(win Window, size Size, env Environment) *Html {
	if w.win != nil {
		return nil
	}
	w.win = win
	rSize := ClampBounds(w.bounds, size)
	eventChan := make(chan bool)
	win.RegisterEventChan(w.id, "click", eventChan)
	go func() {
		for {
			// Also: handle destroy messages?
			select {
			case <- eventChan:
				w.action()

			case newSize := <- env.ResizeChan:
				rSize := ClampBounds(w.bounds, newSize)
				win.UpdateSize(w.id, rSize)
			}
		}
	}()
	return NewHtml("button").
		Id(w.id.String()).
		Text(w.label).
		Styles(DefaultStyle(rSize)).
		Style("cursor", "pointer").
		Event("click")
}

func (w *Button) BoundsOf() Bounds {
	return w.bounds
}
