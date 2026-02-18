package main

import (
	"fmt"
	"github.com/rpucella/go-exene"
)


// Ideally, drop this in a configuration file.
var browser []string = []string{
	"/Applications/Firefox.app/Contents/MacOS/firefox",
	"--new-tab",
}

func main() {
	widget := MainWidget()
	shell := exene.NewShell(widget)
	app := exene.NewBrowserApp(shell, browser)
	app.Start()
}

func MainWidget() exene.Widget {
	count := 0
	title := exene.NewText(
		exene.FixedBounds(400, 50),
		"Sample Counting Example",
	)
	/*
		WithStyle("fontSize", "24px").
		WithStyle("color", "white").
		WithStyle("backgroundColor", "#666666").
	*/
	label := exene.NewText(
		exene.FixedBounds(400, 40),
		"Count = 0",
	)
	/*
	    WithStyle("fontSize", "24px")
	*/
	setLabel := func(newCount int) {
		count = newCount
		label.UpdateText(fmt.Sprintf("Count = %d", count))
	}
	increment := exene.NewButton(
		exene.FixedBounds(120, 40),
		"Increment",
		func() { setLabel(count + 1) },
	)
	/*
		WithStyle("width", "100px")
	*/
	reset := exene.NewButton(
		exene.FixedBounds(120, 40),
		"Reset",
		func() { setLabel(0) },
	)
	main := exene.NewBox(
		exene.BoxVtCenter{
			[]exene.BoxEntry{
				exene.BoxWidget{title},
				exene.BoxWidget{
					exene.NewFrame(
						5, 
						exene.NewBox(
							exene.BoxHzCenter{
								[]exene.BoxEntry{
									exene.BoxWidget{increment},
									exene.BoxGlue{exene.FixedDim(20)},
									exene.BoxWidget{reset},
								},
							},
						),
					),
				},
				exene.BoxGlue{exene.FixedDim(20)},
				exene.BoxWidget{label},
			},
		},
	)
	return main
}
