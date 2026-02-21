package exene

import (
)

func Center(widget Widget) *Box {
	return NewBox(
		NewVtCenter(
			NewGlue(Dim{0, 0, -1}),
			NewHzCenter(
				NewGlue(Dim{0, 0, -1}),
				NewWBox(widget),
				NewGlue(Dim{0, 0, -1}),
			),
			NewGlue(Dim{0, 0, -1}),
		),
	)
}
