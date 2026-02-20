package main

import (
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
	title := exene.NewFrame(
		5,
		exene.RgbHex("333333"),
		exene.NewText(
			exene.Bounds{exene.Dim{500, 700, 1000}, exene.FixDim(50)},
			"Resizing Demo",
			exene.WithFontSize(36),
			exene.WithTextAlign("center"),
			exene.WithBackgroundColor("#cccccc"),
		),
	)
	subtitle := exene.NewFrame(
		5,
		exene.RgbHex("333333"),
		exene.NewText(
			exene.Bounds{exene.Dim{700, 700, -1}, exene.FixDim(50)},
			"(Subtitle)",
			exene.WithTextAlign("center"),
		),
	)
	main := exene.NewBox(
		exene.BoxVtCenter{
			[]exene.BoxEntry{
				exene.BoxWidget{title},
				exene.BoxGlue{exene.Dim{100, 100, 500}},
				exene.BoxWidget{subtitle},
			},
		},
	)
	return main
}
