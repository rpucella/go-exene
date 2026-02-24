package exene

import (
)

type Pile struct {
	id WId
	widgets []Widget
	win Window
	insertChan chan Pair[int, Widget]
	dropChan chan int
	activateChan chan int
}

func NewPile(widgets ...Widget) *Pile {
	// Shortcut!!!
	id := NewId()
	return &Pile{id: id, widgets: widgets}
}

func (w *Pile) Realize(win Window, size Size, resizeChan chan Size) *Html {
	if w.win != nil {
		return nil
	}
	w.win = win
	insertChan := make(chan Pair[int, Widget])
	w.insertChan = insertChan
	dropChan := make(chan int)
	w.dropChan = dropChan
	activateChan := make(chan int)
	w.activateChan = activateChan
	bounds := w.BoundsOf()
	rSize := ClampBounds(bounds, size)
	subHtmls := make([]*Html, len(w.widgets))
	subResizeChans := make([]chan Size, len(w.widgets))
	for i, sw := range w.widgets {
		subResizeChans[i] = make(chan Size)
		subHtmls[i] = sw.Realize(win, size, subResizeChans[i])
		if i > 0 {
			// Mimic how we hide via the SDK.
			savedDisplay := subHtmls[i].GetStyle("display")
			subHtmls[i] = subHtmls[i].
				Style("display", "none").
				Attr("data-display-save", savedDisplay)
		}
	}
	go func() {
		currActive := 0
		currSize := size
		//currentWidgets := boxes  // Shared.
		//currentResizeChans := subResizeChans
		currSize = currSize
		for {
			select {
			case newSize := <- resizeChan:
				currSize = newSize
				rSize := ClampBounds(bounds, newSize)
				win.UpdateSize(w.id, rSize)

			case index := <- activateChan:
				win.HideChild(w.id, currActive)
				currActive = index
				win.UnhideChild(w.id, currActive)

			case <- insertChan:
			case <- dropChan:
				
				/*
			case idxBox := <- insertChan:
				subResizeChan := make(chan Size)
				currentBoxes = insertInto(currentBoxes, idxBox.index, idxBox.box)
				currentResizeChans = insertInto(currentResizeChans, idxBox.index, subResizeChan)
				bounds := b.boxBoundsOf(parentDir)
				rSize := ClampBounds(bounds, currSize)
				sizes := calculateSizes(currentBoxes, dir, rSize)
				subHtml := idxBox.box.boxRealize(win, sizes[idxBox.index], dir, subResizeChan, nil, nil)
				if idxBox.index == len(currentBoxes) - 1 {
					win.AppendChild(id, subHtml)
				} else {
					win.InsertChild(id, idxBox.index, subHtml)
				}
				for i, ch := range currentResizeChans {
					if i != idxBox.index {
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
				*/
			}
		}
	}()
	return NewHtml("div").
		Id(w.id.String()).
		Styles(DefaultStyle(rSize)).
		AppendAll(subHtmls)
}

func (w *Pile) BoundsOf() Bounds {
	// max across all piles?
	width := FixDim(0)
	height := FixDim(0)
	for _, sw := range w.widgets {
		b := sw.BoundsOf()
		width = MaxDim(width, b.Width)
		height = MaxDim(height, b.Height)
	}
	return Bounds{width, height}
}

func (w *Pile) Insert(idx int, widget Widget) {
	w.insertChan <- NewPair(idx, widget)
}

func (w *Pile) Delete(idx int) {
	w.dropChan <- idx
}

func (w *Pile) Activate(idx int) {
	w.activateChan <- idx
}
