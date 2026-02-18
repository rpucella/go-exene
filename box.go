package exene

import (
	"fmt"
)

type BoxEntry interface{
	boxRealize(*WebInterface, int, int, direction) Html
	boxBoundsOf(direction) Bounds
}

type direction int

const (
	verticalDir direction = iota
	horizontalDir
	noDir
)

type BoxVtLeft struct {
	Boxes []BoxEntry
}

type BoxVtCenter struct {
	Boxes []BoxEntry
}

type BoxVtRight struct {
	Boxes []BoxEntry
}

type BoxHzTop struct {
	Boxes []BoxEntry
}

type BoxHzCenter struct {
	Boxes []BoxEntry
}

type BoxHzBottom struct {
	Boxes []BoxEntry
}

type BoxWidget struct {
	Widget Widget
}

type BoxGlue struct {
	Dim Dim
}

func verticalLayout(webIfc *WebInterface, width int, height int, dir direction, boxes []BoxEntry, align string) Html {
	subHtmls := make([]Html, len(boxes))
	boxWidth := 0
	boxHeight := 0
	for i, w := range boxes {
		boxWidth = max(boxWidth, w.boxBoundsOf(dir).Width.Nat)
		boxHeight += w.boxBoundsOf(dir).Height.Nat
		subHtmls[i] = w.boxRealize(webIfc, width, height / len(boxes), dir)
	}
	return Html{
		"",
		"div",
	    nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", boxHeight),
			"width": fmt.Sprintf("%dpx", boxWidth),
			"overflow": "hidden",
			"display": "flex",
			"flex-direction": "column",
			"align-items": align,
		},
		"",
		subHtmls,
	    nil,
	}
}

func layoutBoundsOf(boxes []BoxEntry, dir direction) Bounds {
	boxWidth := 0
	boxHeight := 0
	for _, w := range boxes {
		if dir == verticalDir {
			boxWidth = max(boxWidth, w.boxBoundsOf(dir).Width.Nat)
			boxHeight += w.boxBoundsOf(dir).Height.Nat
		}
		if dir == horizontalDir {
			boxWidth += w.boxBoundsOf(dir).Width.Nat
			boxHeight = max(boxHeight, w.boxBoundsOf(dir).Height.Nat)
		}
	}
	return FixedBounds(boxWidth, boxHeight)
}

func (b BoxVtLeft) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return verticalLayout(webIfc, width, height, verticalDir, b.Boxes, "flex-start")
}

func (b BoxVtLeft) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, verticalDir)
}

func (b BoxVtCenter) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return verticalLayout(webIfc, width, height, verticalDir, b.Boxes, "center")
}

func (b BoxVtCenter) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, verticalDir)
}

func (b BoxVtRight) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return verticalLayout(webIfc, width, height, verticalDir, b.Boxes, "flex-end")
}

func (b BoxVtRight) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, verticalDir)
}


func horizontalLayout(webIfc *WebInterface, width int, height int, dir direction, boxes []BoxEntry, align string) Html {
	subHtmls := make([]Html, len(boxes))
	boxWidth := 0
	boxHeight := 0
	for i, w := range boxes {
		boxWidth += w.boxBoundsOf(dir).Width.Nat
		boxHeight = max(boxHeight, w.boxBoundsOf(dir).Height.Nat)
		subHtmls[i] = w.boxRealize(webIfc, width / len(boxes), height, dir)
	}
	return Html{
		"",
		"div",
	    nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", boxHeight),
			"width": fmt.Sprintf("%dpx", boxWidth),
			"overflow": "hidden",
			"display": "flex",
			"flex-direction": "row",
			"align-items": align,
		},
		"",
		subHtmls,
	    nil,
	}
}

func (b BoxHzTop) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return horizontalLayout(webIfc, width, height, horizontalDir, b.Boxes, "flex-start")
}

func (b BoxHzTop) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, horizontalDir)
}

func (b BoxHzCenter) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return horizontalLayout(webIfc, width, height, horizontalDir, b.Boxes, "center")
}

func (b BoxHzCenter) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, horizontalDir)
}

func (b BoxHzBottom) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return horizontalLayout(webIfc, width, height, horizontalDir, b.Boxes, "flex-end")
}

func (b BoxHzBottom) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, horizontalDir)
}

func (b BoxWidget) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return b.Widget.Realize(webIfc, width, height)
}

func (b BoxWidget) boxBoundsOf(dir direction) Bounds {
	return b.Widget.BoundsOf()
}

func (b BoxGlue) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	boxWidth := 0
	boxHeight := 0
	if dir == verticalDir {
		boxHeight = b.Dim.Nat
	}
	if dir == horizontalDir {
		boxWidth = b.Dim.Nat
	}
	return Html{
		"",
		"div",
	    nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", boxHeight),
			"width": fmt.Sprintf("%dpx", boxWidth),
		},
		"",
		nil,
	    nil,
	}
}

func (b BoxGlue) boxBoundsOf(dir direction) Bounds {
	var width Dim
	var height Dim
	if dir == verticalDir {
		height = b.Dim
	}
	if dir == horizontalDir {
		width = b.Dim
	}
	return Bounds{width, height}
}

type Box struct {
	bounds Bounds
	Id string
	Box BoxEntry
	webIfc *WebInterface
}

func (w *Box) Realize(webIfc *WebInterface, width int, height int) Html {
	w.webIfc = webIfc
	html := w.Box.boxRealize(webIfc, width, height, noDir)
	return html
}

func (w *Box) BoundsOf() Bounds {
	return w.bounds
}

func NewBox(box BoxEntry) *Box {
	// Shortcut!!!
	bounds := box.boxBoundsOf(noDir)
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	return &Box{bounds, strId, box, nil}
}

