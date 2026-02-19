package exene

import (
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

func TransDim(d1 Dim, v int) Dim {
	// Account for the fact that v could be negative.
	newMin := max(0, d1.Min + v)
	newNat := max(0, d1.Nat + v)
	newMax := -1
	if d1.Max >= 0 {
		newMax = max(0, d1.Max + v)
	}
	return Dim{newMin, newNat, newMax}
}

func FixBounds(w int, h int) Bounds {
	return Bounds{FixDim(w), FixDim(h)}
}

func CompatibleBounds(b Bounds, size Size) bool {
	return CompatibleDim(b.Width, size.Width) && CompatibleDim(b.Height, size.Height)
}

func ClampBounds(b Bounds, size Size) Size {
	return Size{ClampDim(b.Width, size.Width), ClampDim(b.Height, size.Height)}
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
	Realize(*WebInterface, Size, chan Size) Html
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

func (sh Shell) Init(webIfc *WebInterface, size Size, resizeChan chan Size) Html {
	return sh.widget.Realize(webIfc, size, resizeChan)
}


