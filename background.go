package exene

import (
)

type Background struct {
	id WId
	bgColor string
	fgColor string
	widget Widget
	win Window
}

func NewBackground(bgColor Color, fgColor Color, widget Widget) *Background {
	id := NewId()
	frame := &Background{id, bgColor.String(), fgColor.String(), widget, nil}
	return frame
}

func (w *Background) Realize(win Window, size Size, env Environment) *Html {
	if w.win != nil {
		return nil
	}
	w.win = win
	subResizeChan := make(chan Size)
	subEnv := Environment{subResizeChan, nil, nil}
	rSize := ClampBounds(w.BoundsOf(), size)
	subHtml := w.widget.Realize(win, rSize, subEnv)
	go func() {
		for {
			select {
			case newSize := <- env.ResizeChan:
				rSize := ClampBounds(w.BoundsOf(), newSize)
				subResizeChan <- rSize
				win.UpdateSize(w.id, rSize)
			}
		}
	}()
	return NewHtml("div").
		Id(w.id.String()).
		Styles(DefaultStyle(rSize)).
		Style("backgroundColor", w.bgColor).
		Style("color", w.fgColor).
		Append(subHtml)
}

func (w *Background) BoundsOf() Bounds {
	bounds := w.widget.BoundsOf()
	return bounds
}
