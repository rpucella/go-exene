package main

import (
	"fmt"
	ex "github.com/rpucella/go-exene"
)


// Ideally, drop this in a configuration file.
var browser []string = []string{
	"/Applications/Firefox.app/Contents/MacOS/firefox",
	"--new-tab",
}

func main() {
	widget := MainWidget()
	shell := ex.NewShell(widget)
	app := ex.NewBrowserApp(shell, browser)
	app.Start()
}

func MainWidget() ex.Widget {
	count := 0
	title := ex.NewText(
		ex.Bounds{ex.Dim{400, 400, 400}, ex.FixDim(50)},
		"Counter Demo",
		ex.WithTextAlign("center"),
		ex.WithFontSize(36),
	)
	label := ex.NewText(
		ex.FixBounds(400, 40),
		"Count = 0",
		ex.WithTextAlign("center"),
	)
	setLabel := func(newCount int) {
		count = newCount
		label.UpdateText(fmt.Sprintf("Count = %d", count))
	}
	increment := ex.NewButton(
		ex.FixBounds(120, 40),
		"Increment",
		func() { setLabel(count + 1) },
	)
	reset := ex.NewButton(
		ex.FixBounds(120, 40),
		"Reset",
		func() { setLabel(0) },
	)
	content := ex.NewBox(
		ex.BoxVtCenter{
			[]ex.BoxEntry{
				ex.BoxWidget{title},
				ex.BoxGlue{ex.Dim{20, 20, 50}},
				ex.BoxHzCenter{
					[]ex.BoxEntry{
						ex.BoxGlue{ex.Dim{0, 0, 100}},
						ex.BoxWidget{increment},
						ex.BoxGlue{ex.Dim{20, 20, 100}},
						ex.BoxWidget{reset},
						ex.BoxGlue{ex.Dim{0, 0, 100}},
					},
				},
				ex.BoxGlue{ex.Dim{20, 20, 50}},
				ex.BoxWidget{label},
			},
		},
	)
	main := ex.Center(
		ex.NewFrame(
			5,
			ex.RgbHex("808080"),
			ex.NewPadding(
				20, 
				content,
			),
		),
	)
	return main
}
