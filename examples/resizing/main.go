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
		10,
		exene.NewText(
			exene.Bounds{exene.Dim{500, 700, 1000}, exene.FixDim(50)},
			"Resizing Demo",
		),
	)
	subtitle := exene.NewFrame(
		10,
		exene.NewText(
			exene.Bounds{exene.Dim{700, 700, -1}, exene.FixDim(50)},
			"(Subtitle)",
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
