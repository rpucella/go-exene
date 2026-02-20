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
		exene.Bounds{exene.Dim{400, 400, 400}, exene.FixDim(50)},
		"Counter Demo",
		exene.WithTextAlign("center"),
		exene.WithFontSize(36),
	)
	/*
		WithStyle("fontSize", "24px").
		WithStyle("color", "white").
		WithStyle("backgroundColor", "#666666").
	*/
	label := exene.NewText(
		exene.FixBounds(400, 40),
		"Count = 0",
		exene.WithTextAlign("center"),
	)
	/*
	    WithStyle("fontSize", "24px")
	*/
	setLabel := func(newCount int) {
		count = newCount
		label.UpdateText(fmt.Sprintf("Count = %d", count))
	}
	increment := exene.NewButton(
		exene.FixBounds(120, 40),
		"Increment",
		func() { setLabel(count + 1) },
	)
	/*
		WithStyle("width", "100px")
	*/
	reset := exene.NewButton(
		exene.FixBounds(120, 40),
		"Reset",
		func() { setLabel(0) },
	)
	content := exene.NewBox(
		exene.BoxVtCenter{
			[]exene.BoxEntry{
				exene.BoxWidget{title},
				exene.BoxGlue{exene.Dim{20, 20, 50}},
				exene.BoxHzCenter{
					[]exene.BoxEntry{
						exene.BoxGlue{exene.Dim{0, 0, 100}},
						exene.BoxWidget{increment},
						exene.BoxGlue{exene.Dim{20, 20, 100}},
						exene.BoxWidget{reset},
						exene.BoxGlue{exene.Dim{0, 0, 100}},
					},
				},
				exene.BoxGlue{exene.Dim{20, 20, 50}},
				exene.BoxWidget{label},
			},
		},
	)
	main := exene.Center(
		exene.NewFrame(
			5,
			exene.RgbHex("808080"),
			exene.NewPadding(
				20, 
				content,
			),
		),
	)
	return main
}
