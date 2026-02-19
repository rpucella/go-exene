package exene

import (
	"fmt"
	"math"
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

func layout(webIfc *WebInterface, width int, height int, dir direction, bounds Bounds, boxes []BoxEntry, align string) Html {
	subHtmls := make([]Html, len(boxes))
	rWidth, rHeight := ClampBounds(bounds, width, height)
	var sizes []int
	if dir == verticalDir {
		dims := make([]Dim, len(boxes))
		for i, b := range boxes {
			dims[i] = b.boxBoundsOf(dir).Height
		}
		sizes = calculatePartition(dims, rHeight)
	}
	if dir == horizontalDir {
		dims := make([]Dim, len(boxes))
		for i, b := range boxes {
			dims[i] = b.boxBoundsOf(dir).Width
		}
		sizes = calculatePartition(dims, rWidth)
	}
	for i, b := range boxes {
		if dir == verticalDir {
			subHtmls[i] = b.boxRealize(webIfc, width, sizes[i], dir)
		}
		if dir == horizontalDir {
			subHtmls[i] = b.boxRealize(webIfc, sizes[i], height, dir)
		}
	}
	flexDirection := ""
	if dir == verticalDir {
		flexDirection = "column"
	}
	if dir == horizontalDir {
		flexDirection = "row"
	}
	return Html{
		"",
		"div",
	    nil,
		map[string]string{
			"width": fmt.Sprintf("%dpx", rWidth),
			"height": fmt.Sprintf("%dpx", rHeight),
			"overflow": "hidden",
			"display": "flex",
			"flex-direction": flexDirection,
			"align-items": align,
		},
		"",
		subHtmls,
	    nil,
	}
}

func calculatePartition(bnds []Dim, size int) []int {
	result := make([]int, len(bnds))
	bTotal := Dim{0, 0, 0}
	for i, bb := range bnds {
		bTotal = AddDim(bTotal, bb)
		result[i] = bb.Nat
	}
	//fmt.Println("----------------------------------------")
	if size < bTotal.Nat {
		for {
			//fmt.Println("  ", result)
			current := 0
			countAboveMin := 0
			delta := bTotal.Nat
			for i, bb := range bnds {
				current += result[i]
				if result[i] > bb.Min {
					delta = min(delta, result[i] - bb.Min)
					countAboveMin += 1
				}
			}
			if countAboveMin == 0 {
				// Everything is at min, so let's bail.
				return result
			}
			excess := current - size
			if excess < len(bnds) {
				// We got it - stop.
				return result
			}
			toSubtract := delta
			if excess < countAboveMin * delta {
				toSubtract = int(math.Ceil(float64(excess) / float64(countAboveMin)))
			}
			for i, bb := range bnds {
				if result[i] > bb.Min {
					result[i] = result[i] - toSubtract
				}
			}
		}
		return result
	}
	if size > bTotal.Nat {
		for {
			//fmt.Println("  ", result)
			current := 0
			countBelowMax := 0
			delta := -1
			for i, bb := range bnds {
				current += result[i]
				if bb.Max < 0 {
					countBelowMax += 1
				} else if result[i] < bb.Max {
					if delta < 0 {
						delta = bb.Max - result[i]
					} else {
						delta = min(delta, bb.Max - result[i])
					}
					countBelowMax += 1
				}
			}
			if countBelowMax == 0 {
				// Everything is at min, so let's bail.
				return result
			}
			excess := size - current
			if excess < len(bnds) {
				// We got it - stop.
				return result
			}
			toAdd := delta
			// delta = -1 means that no bounds has a max, so we can just allocate the excess uniformly.
			if delta == -1 || excess < countBelowMax * delta {
				toAdd = int(math.Floor(float64(excess) / float64(countBelowMax)))
			}
			for i, bb := range bnds {
				if bb.Max < 0 || result[i] < bb.Max {
					result[i] = result[i] + toAdd
				}
			}
		}
		return result
	}
	return result
}

func layoutBoundsOf(boxes []BoxEntry, dir direction) Bounds {
	width := FixDim(0)
	height := FixDim(0)
	for _, b := range boxes {
		if dir == verticalDir {
			bb := b.boxBoundsOf(dir)
			width = MaxDim(width, bb.Width)
			height = AddDim(height, bb.Height)
		}
		if dir == horizontalDir {
			bb := b.boxBoundsOf(dir)
			width = AddDim(width, bb.Width)
			height = MaxDim(height, bb.Height)
		}
	}
	return Bounds{width, height}
}

func (b BoxVtLeft) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return layout(webIfc, width, height, verticalDir, b.boxBoundsOf(dir), b.Boxes, "flex-start")
}

func (b BoxVtLeft) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, verticalDir)
}

func (b BoxVtCenter) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return layout(webIfc, width, height, verticalDir, b.boxBoundsOf(dir), b.Boxes, "center")
}

func (b BoxVtCenter) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, verticalDir)
}

func (b BoxVtRight) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return layout(webIfc, width, height, verticalDir, b.boxBoundsOf(dir), b.Boxes, "flex-end")
}

func (b BoxVtRight) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, verticalDir)
}

func (b BoxHzTop) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return layout(webIfc, width, height, horizontalDir, b.boxBoundsOf(dir), b.Boxes, "flex-start")
}

func (b BoxHzTop) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, horizontalDir)
}

func (b BoxHzCenter) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return layout(webIfc, width, height, horizontalDir, b.boxBoundsOf(dir), b.Boxes, "center")
}

func (b BoxHzCenter) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, horizontalDir)
}

func (b BoxHzBottom) boxRealize(webIfc *WebInterface, width int, height int, dir direction) Html {
	return layout(webIfc, width, height, horizontalDir, b.boxBoundsOf(dir), b.Boxes, "flex-end")
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
	bounds := b.boxBoundsOf(dir)
	boxWidth := ClampDim(bounds.Width, width)
	boxHeight := ClampDim(bounds.Height, height)
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

