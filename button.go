package exene

import (
)

type Button struct {
    bounds Bounds
	id WId
	label string
	style *Style
	action func()
	win Window
}


func NewButton(bounds Bounds, label string, act func(), styles ...StyleOption) *Button {
	id := NewId()
	style := &Style{}
	for _, s := range styles {
		s(style)
	}
	button := &Button{bounds, id, label, style, act, nil}
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
		Styles(w.style.AsMap()).
		Style("cursor", "pointer").
		Event("click")
}

func (w *Button) BoundsOf() Bounds {
	return w.bounds
}
