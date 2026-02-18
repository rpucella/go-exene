package exene

import (
	"fmt"
)

type Dim struct {
	Min int
	Nat int
}

type Bounds struct {
	Width Dim
	Height Dim
}

func FixedDim(v int) Dim {
	return Dim{v, v}
}

func FixedBounds(w int, h int) Bounds {
	return Bounds{FixedDim(w), FixedDim(h)}
}

type Html struct {
	Id string `json:"id"`
	Tag string `json:"tag"`
	Attrs map[string]string `json:"attrs"`
	Style map[string]string `json:"style"`
	Text string `json:"text"`
	Children []Html `json:"children"`
	Events []string `json:"events"`
}

type Widget interface {
	BoundsOf() Bounds
	// May also want to pass the environment?
	Realize(*WebInterface, int, int) Html
}


var id int = 0

func NewId() int {
	newId := id
	id += 1
	return newId
}


type Shell struct {
	root bool
	widget Widget
}

func NewShell(w Widget) Shell {
	return Shell{false, w}
}

func (sh Shell) Init(webIfc *WebInterface, width int, height int) Html {
	return sh.widget.Realize(webIfc, width, height)
}

/*
   ************************************************************
   
     Widget library

   ************************************************************
*/

type Button struct {
    bounds Bounds
	Id string
	Label string
	Action func()
	webIfc *WebInterface
}

func (w *Button) Realize(webIfc *WebInterface, width int, height int) Html {
	w.webIfc = webIfc
	eventDispatch := make(chan bool)
	go func() {
		for {
			// Also: handle destroy messages?
			select {
			case <- eventDispatch:
				w.Action()
			}
		}
	}()
	webIfc.dispatchMap[w.Id] = eventDispatch
	return Html{
		w.Id,
		"button",
		nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", w.bounds.Height.Nat),
			"width": fmt.Sprintf("%dpx", w.bounds.Width.Nat),
			"overflow": "hidden",
			"border": "none",
		},
		w.Label,
		nil,
		[]string{"click"},
	}
}

func (w *Button) BoundsOf() Bounds {
	return w.bounds
}


func NewButton(bounds Bounds, label string, act func()) *Button {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	button := &Button{bounds, strId, label, act, nil}
	return button
}


type Text struct {
    bounds Bounds
	Id string 
	Text string
	webIfc *WebInterface
}

func (w *Text) Realize(webIfc *WebInterface, width int, height int) Html {
	w.webIfc = webIfc
	return Html{
		w.Id,
		"div",
		nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", w.bounds.Height.Nat),
			"width": fmt.Sprintf("%dpx", w.bounds.Width.Nat),
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

func NewText(bounds Bounds, text string) *Text {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	textWidget := &Text{bounds, strId, text, nil}
	return textWidget
}



type Frame struct {
	Id string
	size int
	widget Widget
	webIfc *WebInterface
}

func (w *Frame) Realize(webIfc *WebInterface, width int, height int) Html {
	w.webIfc = webIfc
	bounds := w.BoundsOf()
	subHtml := w.widget.Realize(webIfc, width, height)
	return Html{
		w.Id,
		"div",
		nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", bounds.Height.Nat),
			"width": fmt.Sprintf("%dpx", bounds.Width.Nat),
			"overflow": "hidden",
			"border": fmt.Sprintf("%dpx solid black", w.size),
		},
		"",
		[]Html{
			subHtml,
		},
		nil,
	}
}

func (w *Frame) BoundsOf() Bounds {
	bounds := w.widget.BoundsOf()
	bounds.Width.Nat = 2 * w.size + bounds.Width.Nat
	bounds.Height.Nat = 2 * w.size + bounds.Height.Nat
	return bounds
}

func NewFrame(size int, widget Widget) *Frame {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	frame := &Frame{strId, size, widget, nil}
	return frame
}

