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
	title := exene.NewText(
		exene.FixedBounds(200, 50),
		"Paragraphs Demo",
	)
	para1 := exene.NewText(
		exene.FixedBounds(400, 100),
		"Lorem ipsum dolor sit amet.",
	)
	para2 := exene.NewText(
		exene.FixedBounds(300, 50),
		"Whereas now I can achieve the same by simplifying using built-in min() method, without casting over and over the operands.",
	)
	return exene.NewBox(
		exene.BoxVtLeft{
			[]exene.BoxEntry{
				exene.BoxWidget{title},
				exene.BoxWidget{para1},
				exene.BoxWidget{para2},
			},
		},
	)
}
