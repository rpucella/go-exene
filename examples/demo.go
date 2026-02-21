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
	dark := ex.RgbHex("333333")
	light := ex.RgbHex("eeeeee")
	title := ex.NewLabel(
		ex.FixBounds(400, 50),
		"Counter Demo",
		ex.WithAlign("center"),
		ex.WithFontSize(36),
	)
	label := ex.NewLabel(
		ex.FixBounds(400, 40),
		"Count = 0",
		ex.WithAlign("center"),
	)
	setLabel := func(newCount int) {
		count = newCount
		label.UpdateText(fmt.Sprintf("Count = %d", count))
	}
	increment := ex.NewFrame(
		2,
		dark,
		ex.NewButton(
			ex.FixBounds(120, 40),
			"Increment",
			func() { setLabel(count + 1) },
		),
	)
	reset := ex.NewFrame(
		2,
		dark,
		ex.NewButton(
			ex.FixBounds(120, 40),
			"Reset",
			func() { setLabel(0) },
		),
	)
	var content *ex.Box
	extraCount := 0
	addExtra := ex.NewFrame(
		2,
		dark,
		ex.NewButton(
			ex.FixBounds(120, 40),
			"Add extra",
			func() {
				extraCount += 1
				label := ex.NewLabel(
					ex.FixBounds(200, 20),
					fmt.Sprintf("Extra #%d", extraCount),
					ex.WithAlign("center"),
				)
				content.Insert(4, ex.NewWBox(label))
			},
		),
	)
	dropExtra := ex.NewFrame(
		2,
		dark,
		ex.NewButton(
			ex.FixBounds(120, 40),
			"Drop extra",
			func() {
				if extraCount > 0 {
					content.Delete(4)
					extraCount -= 1
				}
			},
		),
	)
	content = ex.NewBox(
		ex.NewVtCenter(
			ex.NewWBox(title),
			ex.NewGlue(ex.NewDim(20, 20, 50)),
			ex.NewHzCenter(
				ex.NewGlue(ex.NewDim(0, 0, 100)),
				ex.NewWBox(increment),
				ex.NewGlue(ex.NewDim(20, 20, 100)),
				ex.NewWBox(reset),
				ex.NewGlue(ex.NewDim(20, 20, 100)),
				ex.NewWBox(addExtra),
				ex.NewGlue(ex.NewDim(20, 20, 100)),
				ex.NewWBox(dropExtra),
				ex.NewGlue(ex.NewDim(0, 0, 100)),
			),
			ex.NewGlue(ex.NewDim(20, 20, 50)),
			ex.NewWBox(label),
		),
	)
	main := ex.Center(
		ex.NewBackground(
			light, 
			dark,
			ex.NewFrame(
				3,
				dark,
				ex.NewFrame(
					20,
					ex.Transparent,
					content,
				),
			),
		),
	)
	return main
}

