package exene

import (
	"math"
)

type Box struct {
	id WId
	box BoxEntry
	win Window
	insertChan chan Pair[int, BoxEntry]
	dropChan chan int
}

func NewBox(box BoxEntry) *Box {
	id := NewId()
	return &Box{id, box, nil, nil, nil}
}

func (w *Box) Realize(win Window, size Size, env Environment) *Html {
	if w.win != nil {
		return nil
	}
	w.win = win
	insertChan := make(chan Pair[int, BoxEntry])
	w.insertChan = insertChan
	dropChan := make(chan int)
	w.dropChan = dropChan
	html := w.box.boxRealize(win, size, dirNone, env, insertChan, dropChan)
	return html
}

func (w *Box) BoundsOf() Bounds {
	return w.box.boxBoundsOf(dirNone)
}

func (w *Box) Insert(idx int, box BoxEntry) {
	w.insertChan <- NewPair(idx, box)
}

func (w *Box) Delete(idx int) {
	w.dropChan <- idx
}


type BoxEntry interface{
	boxRealize(Window, Size, direction, Environment, chan Pair[int, BoxEntry], chan int) *Html
	boxBoundsOf(direction) Bounds
}

type BoxList struct {
	align alignment
	dir direction
	boxes []BoxEntry
}

type BoxWidget struct {
	widget Widget
}

type BoxGlue struct {
	dim Dim
}

type direction int

const (
	dirVertical direction = iota
	dirHorizontal
	dirNone
)

func (d direction) String() string {
	if d == dirVertical {
		return "column"
	} else {
		// Default out to row, though dirNone should never be a flex box.
		return "row"
	} 
}

type alignment int

const (
	alignStart alignment = iota
	alignCenter 
	alignEnd
)

func (n alignment) String() string {
	if n == alignStart {
		return "flex-start"
	} else if n == alignCenter {
		return "center"
	} else {
		return "flex-end"
	}
}

func NewVtLeft(boxes ...BoxEntry) BoxList {
	return BoxList{alignStart, dirVertical, boxes}
}
	
func NewVtCenter(boxes ...BoxEntry) BoxList {
	return BoxList{alignCenter, dirVertical, boxes}
}
	
func NewVtRight(boxes ...BoxEntry) BoxList {
	return BoxList{alignEnd, dirVertical, boxes}
}
	
func NewHzTop(boxes ...BoxEntry) BoxList {
	return BoxList{alignStart, dirHorizontal, boxes}
}
	
func NewHzCenter(boxes ...BoxEntry) BoxList {
	return BoxList{alignCenter, dirHorizontal, boxes}
}
	
func NewHzBottom(boxes ...BoxEntry) BoxList {
	return BoxList{alignEnd, dirHorizontal, boxes}
}

func NewWBox(widget Widget) BoxWidget {
	return BoxWidget{widget}
}

func NewGlue(d Dim) BoxGlue {
	return BoxGlue{d}
}
	
func (b BoxList) boxRealize(win Window, size Size, parentDir direction, env Environment, insertChan chan Pair[int, BoxEntry], dropChan chan int) *Html {
	dir := b.dir
	align := b.align
	boxes := b.boxes
	bounds := b.boxBoundsOf(parentDir)
	id := NewId()
	subHtmls := make([]*Html, len(boxes))
	rSize := ClampBounds(bounds, size)
	var sizes = calculateSizes(boxes, dir, rSize)
	subResizeChans := make([]chan Size, len(boxes))
	for i, b := range boxes {
		subResizeChans[i] = make(chan Size)
		env := Environment{subResizeChans[i], nil, nil}
		subHtmls[i] = b.boxRealize(win, sizes[i], dir, env, nil, nil)
	}
	go func() {
		currSize := size
		currentBoxes := boxes  // Shared.
		currentResizeChans := subResizeChans
		// probably need to rethink the whole thing.
		// instead of thread per box, maybe I can use the fact that I have IDs for all the subwidgets
		// and only track sizes for them?
		for {
			select {
			case newSize := <- env.ResizeChan:
				currSize = newSize
				rSize := ClampBounds(bounds, newSize)
				subSizes := calculateSizes(currentBoxes, dir, rSize)
				for i, ch := range subResizeChans {
					ch <- subSizes[i]
				}
				win.UpdateSize(id, rSize)

			case pair := <- insertChan:
				index, box := pair.Get()
				subResizeChan := make(chan Size)
				currentBoxes = insertInto(currentBoxes, index, box)
				currentResizeChans = insertInto(currentResizeChans, index, subResizeChan)
				bounds := b.boxBoundsOf(parentDir)
				rSize := ClampBounds(bounds, currSize)
				sizes := calculateSizes(currentBoxes, dir, rSize)
				env := Environment{subResizeChan, nil, nil}
				subHtml := box.boxRealize(win, sizes[index], dir, env, nil, nil)
				if index == len(currentBoxes) - 1 {
					win.AppendChild(id, subHtml)
				} else {
					win.InsertChild(id, index, subHtml)
				}
				for i, ch := range currentResizeChans {
					if i != index {
						// Skip the one we just inserted!
						ch <- sizes[i]
					}
				}
				// That's not enough:
				win.UpdateSize(id, rSize)
				// Resizing due to internal changes need to propagate up the hierarchy!
				// You don't get that on resize, obviously, because there everybody gets the resize require
				// But this is a resize request that's propagating up!

			case index := <- dropChan:
				currentBoxes = deleteFrom(currentBoxes, index)
				currentResizeChans = deleteFrom(currentResizeChans, index)
				bounds := b.boxBoundsOf(parentDir)
				rSize := ClampBounds(bounds, currSize)
				sizes := calculateSizes(currentBoxes, dir, rSize)
				win.DeleteChild(id, index)
				// Need to destroy the deleted box!
				for i, ch := range currentResizeChans {
					ch <- sizes[i]
				}
				win.UpdateSize(id, rSize)
			}
		}
	}()
	return NewHtml("div").
		Id(id.String()).
		Styles(DefaultStyle(rSize)).
		Style("display", "flex").
		Style("flexDirection", dir.String()).
		Style("alignItems", align.String()).
		AppendAll(subHtmls)
}

func (b BoxList) boxBoundsOf(parentDir direction) Bounds {
	boxes := b.boxes
	dir := b.dir
	width := FixDim(0)
	height := FixDim(0)
	for _, b := range boxes {
		bb := b.boxBoundsOf(dir)
		if dir == dirVertical {
			width = MaxDim(width, bb.Width)
			height = AddDim(height, bb.Height)
		}
		if dir == dirHorizontal {
			width = AddDim(width, bb.Width)
			height = MaxDim(height, bb.Height)
		}
	}
	return Bounds{width, height}
}

func (b BoxWidget) boxRealize(win Window, size Size, parentDir direction, env Environment, insertChan chan Pair[int, BoxEntry], dropChan chan int) *Html {
	go func(){
		for {
			select{
				case <- insertChan:
				case <- dropChan:
			}
		}
	}()
	return b.widget.Realize(win, size, env)
}

func (b BoxWidget) boxBoundsOf(parentDir direction) Bounds {
	return b.widget.BoundsOf()
}

func (b BoxGlue) boxRealize(win Window, size Size, parentDir direction, env Environment, insertChan chan Pair[int, BoxEntry], dropChan chan int) *Html {
	id := NewId()
	bounds := b.boxBoundsOf(parentDir)
	boxSize := ClampBounds(bounds, size)
	go func() {
		for {
			select {
			case newSize := <- env.ResizeChan:
				rSize := ClampBounds(bounds, newSize)
				win.UpdateSize(id, rSize)

			case <- insertChan:
			case <- dropChan:
			}
		}
	}()
	return NewHtml("div").
		Id(id.String()).
		Styles(DefaultStyle(boxSize))
}

func (b BoxGlue) boxBoundsOf(parentDir direction) Bounds {
	var width Dim
	var height Dim
	if parentDir == dirVertical {
		height = b.dim
	}
	if parentDir == dirHorizontal {
		width = b.dim
	}
	return Bounds{width, height}
}

func calculateSizes(boxes []BoxEntry, dir direction, size Size) []Size {
	var lengths []int
	result := make([]Size, len(boxes))
	if dir == dirVertical {
		dims := make([]Dim, len(boxes))
		for i, b := range boxes {
			dims[i] = b.boxBoundsOf(dir).Height
		}
		lengths = calculatePartition(dims, size.Height)
		for i, n := range lengths {
			result[i] = Size{size.Width, n}
		}
	}
	if dir == dirHorizontal {
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

func insertInto[T any](slice []T, index int, item T) []T {
	if index < len(slice) {
		return append(slice[:index], append([]T{item}, slice[index:]...)...)
	}
	return append(slice, item)
}

func deleteFrom[T any](slice []T, index int) []T {
	if index < len(slice) {
		return append(slice[:index], slice[index + 1:]...)
	}
	return slice
}
