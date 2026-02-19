package exene

import (
	"fmt"
)

type Dim struct {
	Min int
	Nat int
	Max int   // Use <0 for "no max"
}

type Bounds struct {
	Width Dim
	Height Dim
}

func FixDim(v int) Dim {
	return Dim{v, v, v}
}

func CompatibleDim(d Dim, size int) bool {
	if size < d.Min {
		return false
	}
	if d.Max >= 0 && size > d.Max {
		return false
	}
	return true
}

func ClampDim(d Dim, size int) int {
	if size < d.Min {
		return d.Min
	}
	if d.Max >= 0 && size > d.Max {
		return d.Max
	}
	return size
}

func MaxDim(d1 Dim, d2 Dim) Dim {
	newMin := max(d1.Min, d2.Min)
	newNat := max(d1.Nat, d2.Nat)
	newMax := -1
	if d1.Max >= 0 && d2.Max >= 0 {
		newMax = max(d1.Max, d2.Max)
	}
	return Dim{newMin, newNat, newMax}
}

func AddDim(d1 Dim, d2 Dim) Dim {
	newMin := d1.Min + d2.Min
	newNat := d1.Nat + d2.Nat
	newMax := -1
	if d1.Max >= 0 && d2.Max >= 0 {
		newMax = d1.Max + d2.Max
	}
	return Dim{newMin, newNat, newMax}
}

func FixBounds(w int, h int) Bounds {
	return Bounds{FixDim(w), FixDim(h)}
}

func CompatibleBounds(b Bounds, width int, height int) bool {
	return CompatibleDim(b.Width, width) && CompatibleDim(b.Height, height)
}

func ClampBounds(b Bounds, width int, height int) (int, int) {
	return ClampDim(b.Width, width), ClampDim(b.Height, height)
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
	rWidth, rHeight := ClampBounds(w.bounds, width, height)
	return Html{
		w.Id,
		"button",
		nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", rHeight),
			"width": fmt.Sprintf("%dpx", rWidth),
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
	rWidth, rHeight := ClampBounds(w.bounds, width, height)
	return Html{
		w.Id,
		"div",
		nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", rHeight),
			"width": fmt.Sprintf("%dpx", rWidth),
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
	subHtml := w.widget.Realize(webIfc, width, height)
	rWidth, rHeight := ClampBounds(w.BoundsOf(), width, height)
	return Html{
		w.Id,
		"div",
		nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", rHeight),
			"width": fmt.Sprintf("%dpx", rWidth),
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
	dim2 := FixDim(2 * w.size)
	newBounds := Bounds{
		AddDim(bounds.Width, dim2),
		AddDim(bounds.Height, dim2),
	}
	return newBounds
}

func NewFrame(size int, widget Widget) *Frame {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	frame := &Frame{strId, size, widget, nil}
	return frame
}
