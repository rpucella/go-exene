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
	dispatcher := exene.NewDispatcher()
	widget := MainWidget(dispatcher)
	app := exene.NewBrowserApp(dispatcher, widget, browser)
	app.Start()
}

func MainWidget(d exene.Dispatcher) exene.Widget {
	count := 0
	title := exene.NewText(d, "Sample Counting Example").
		WithStyle("fontSize", "24px").
		WithStyle("color", "white").
		WithStyle("backgroundColor", "#666666").
		WithStyle("padding", "16px")
	label := exene.NewText(d, "Count = 0").WithStyle("fontSize", "24px")
	setLabel := func(newCount int) {
		count = newCount
		label.UpdateLabel(fmt.Sprintf("Count = %d", count))
	}
	increment := exene.NewButton(d, "Increment", func() { setLabel(count + 1) }).
		WithStyle("padding", "8px").
		WithStyle("width", "100px")
	reset := exene.NewButton(d, "Reset", func() { setLabel(0) }).
		WithStyle("padding", "8px")
	gap := func(n int) exene.Widget {
		return exene.NewGap(d, fmt.Sprintf("%dpx", n))
	}
	buttons := exene.NewBox(d, "row-center", []exene.Widget{increment, gap(32), reset})
	main := exene.NewBox(d, "column-center", []exene.Widget{title, gap(48), buttons, gap(48), label}).
		WithStyle("padding", "32px")
	return main
}
