package exene

import (
)

func Center(widget Widget) *Box {
	return NewBox(
		BoxVtCenter{
			[]BoxEntry{
				BoxGlue{Dim{0, 0, -1}},
				BoxHzCenter{
					[]BoxEntry{
						BoxGlue{Dim{0, 0, -1}},
						BoxWidget{widget},
						BoxGlue{Dim{0, 0, -1}},
					},
				},
				BoxGlue{Dim{0, 0, -1}},
			},
		},
	)
}
