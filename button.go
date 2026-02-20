package exene

import (
	"fmt"
)

type Button struct {
    bounds Bounds
	Id string
	Label string
	Action func()
	webIfc *WebInterface
}


func NewButton(bounds Bounds, label string, act func()) *Button {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	button := &Button{bounds, strId, label, act, nil}
	return button
}


func (w *Button) Realize(webIfc *WebInterface, size Size, resizeChan chan Size) Html {
	w.webIfc = webIfc
	rSize := ClampBounds(w.bounds, size)
	eventChan := make(chan bool)
	webIfc.dispatchMap[w.Id] = eventChan
	go func() {
		for {
			// Also: handle destroy messages?
			select {
			case <- eventChan:
				w.Action()

			case newSize := <- resizeChan:
				rSize := ClampBounds(w.bounds, newSize)
				webIfc.UpdateSize(w.Id, rSize)
			}
		}
	}()
	return Html{
		w.Id,
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
		w.Label,
		nil,
		[]string{"click"},
	}
}

func (w *Button) BoundsOf() Bounds {
	return w.bounds
}
