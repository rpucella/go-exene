package exene

import (
	"fmt"
	"math"
)

type Box struct {
	bounds Bounds
	Id string
	Box BoxEntry
	webIfc *WebInterface
}

func NewBox(box BoxEntry) *Box {
	// Shortcut!!!
	bounds := box.boxBoundsOf(noDir)
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	return &Box{bounds, strId, box, nil}
}

func (w *Box) Realize(webIfc *WebInterface, size Size, resizeChan chan Size) Html {
	w.webIfc = webIfc
	html := w.Box.boxRealize(webIfc, size, noDir, resizeChan)
	return html
}

func (w *Box) BoundsOf() Bounds {
	return w.bounds
}


type BoxEntry interface{
	boxRealize(*WebInterface, Size, direction, chan Size) Html
	boxBoundsOf(direction) Bounds
}

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

type direction int

const (
	verticalDir direction = iota
	horizontalDir
	noDir
)

func (b BoxVtLeft) boxRealize(webIfc *WebInterface, size Size, dir direction, resizeChan chan Size) Html {
	return layout(webIfc, size, verticalDir, b.boxBoundsOf(dir), b.Boxes, "flex-start", resizeChan)
}

func (b BoxVtLeft) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, verticalDir)
}

func (b BoxVtCenter) boxRealize(webIfc *WebInterface, size Size, dir direction, resizeChan chan Size) Html {
	return layout(webIfc, size, verticalDir, b.boxBoundsOf(dir), b.Boxes, "center", resizeChan)
}

func (b BoxVtCenter) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, verticalDir)
}

func (b BoxVtRight) boxRealize(webIfc *WebInterface, size Size, dir direction, resizeChan chan Size) Html {
	return layout(webIfc, size, verticalDir, b.boxBoundsOf(dir), b.Boxes, "flex-end", resizeChan)
}

func (b BoxVtRight) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, verticalDir)
}

func (b BoxHzTop) boxRealize(webIfc *WebInterface, size Size, dir direction, resizeChan chan Size) Html {
	return layout(webIfc, size, horizontalDir, b.boxBoundsOf(dir), b.Boxes, "flex-start", resizeChan)
}

func (b BoxHzTop) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, horizontalDir)
}

func (b BoxHzCenter) boxRealize(webIfc *WebInterface, size Size, dir direction, resizeChan chan Size) Html {
	return layout(webIfc, size, horizontalDir, b.boxBoundsOf(dir), b.Boxes, "center", resizeChan)
}

func (b BoxHzCenter) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, horizontalDir)
}

func (b BoxHzBottom) boxRealize(webIfc *WebInterface, size Size, dir direction, resizeChan chan Size) Html {
	return layout(webIfc, size, horizontalDir, b.boxBoundsOf(dir), b.Boxes, "flex-end", resizeChan)
}

func (b BoxHzBottom) boxBoundsOf(dir direction) Bounds {
	return layoutBoundsOf(b.Boxes, horizontalDir)
}

func (b BoxWidget) boxRealize(webIfc *WebInterface, size Size, dir direction, resizeChan chan Size) Html {
	return b.Widget.Realize(webIfc, size, resizeChan)
}

func (b BoxWidget) boxBoundsOf(dir direction) Bounds {
	return b.Widget.BoundsOf()
}

func (b BoxGlue) boxRealize(webIfc *WebInterface, size Size, dir direction, resizeChan chan Size) Html {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	bounds := b.boxBoundsOf(dir)
	boxSize := ClampBounds(bounds, size)
	go func() {
		for {
			select {
			case newSize := <- resizeChan:
				rSize := ClampBounds(bounds, newSize)
				webIfc.UpdateSize(strId, rSize)
			}
		}
	}()
	
	return Html{
		strId,
		"div",
	    nil,
		map[string]string{
			"height": fmt.Sprintf("%dpx", boxSize.Height),
			"width": fmt.Sprintf("%dpx", boxSize.Width),
			"overflow": "hidden",
			"transition": "height 0.1s, width 0.1s",
			"boxSizing": "border-box",
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

func layout(webIfc *WebInterface, size Size, dir direction, bounds Bounds, boxes []BoxEntry, align string, resizeChan chan Size) Html {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	subHtmls := make([]Html, len(boxes))
	rSize := ClampBounds(bounds, size)
	var sizes = calculateSizes(boxes, dir, rSize)
	subResizeChans := make([]chan Size, len(boxes))
	for i, b := range boxes {
		subResizeChans[i] = make(chan Size)
		subHtmls[i] = b.boxRealize(webIfc, sizes[i], dir, subResizeChans[i])
	}
	go func() {
		for {
			select {
			case newSize := <- resizeChan:
				rSize := ClampBounds(bounds, newSize)
				subSizes := calculateSizes(boxes, dir, rSize)
				for i, ch := range subResizeChans {
					ch <- subSizes[i]
				}
				webIfc.UpdateSize(strId, rSize)
			}
		}
	}()
	flexDirection := ""
	if dir == verticalDir {
		flexDirection = "column"
	}
	if dir == horizontalDir {
		flexDirection = "row"
	}
	return Html{
		strId,
		"div",
	    nil,
		map[string]string{
			"width": fmt.Sprintf("%dpx", rSize.Width),
			"height": fmt.Sprintf("%dpx", rSize.Height),
			"overflow": "hidden",
			"display": "flex",
			"flex-direction": flexDirection,
			"align-items": align,
			"transition": "height 0.1s, width 0.1s",
			"boxSizing": "border-box",
		},
		"",
		subHtmls,
	    nil,
	}
}

func calculateSizes(boxes []BoxEntry, dir direction, size Size) []Size {
	var lengths []int
	result := make([]Size, len(boxes))
	if dir == verticalDir {
		dims := make([]Dim, len(boxes))
		for i, b := range boxes {
			dims[i] = b.boxBoundsOf(dir).Height
		}
		lengths = calculatePartition(dims, size.Height)
		for i, n := range lengths {
			result[i] = Size{size.Width, n}
		}
	}
	if dir == horizontalDir {
		dims := make([]Dim, len(boxes))
		for i, b := range boxes {
			dims[i] = b.boxBoundsOf(dir).Width
		}
		lengths = calculatePartition(dims, size.Width)
		for i, n := range lengths {
			result[i] = Size{n, size.Height}
		}
	}
	return result
}


func calculatePartition(bnds []Dim, length int) []int {
	result := make([]int, len(bnds))
	bTotal := Dim{0, 0, 0}
	for i, bb := range bnds {
		bTotal = AddDim(bTotal, bb)
		result[i] = bb.Nat
	}
	//fmt.Println("----------------------------------------")
	if length < bTotal.Nat {
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
			excess := current - length
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
	if length > bTotal.Nat {
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
			excess := length - current
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

