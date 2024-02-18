package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Tom5521/gt4tools/pkg/boxes"
	t "github.com/Tom5521/gt4tools/pkg/tools"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	app := gtk.NewApplication("com.test.window", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() {
		activate(app)
	})
	if code := app.Run(os.Args); code > 0 {
		os.Exit(0)
	}
}

func activate(app *gtk.Application) {
	w := gtk.NewApplicationWindow(app)
	w.SetDefaultSize(500, 400)

	var buttons []*gtk.Button
	for i := range 30 {
		button := gtk.NewButtonWithLabel("Button " + strconv.Itoa(i))
		button.ConnectClicked(func() {
			fmt.Println(i)
		})
		buttons = append(buttons, button)
	}

	var labels []gtk.Widgetter
	for i := range 30 {
		labels = append(labels, gtk.NewLabel("Label "+strconv.Itoa(i)))
	}

	vbox := boxes.NewHbox(
		boxes.NewScrolledVbox(
			// Convert a slice of a specific type to a gtk.Widgetter slice.
			t.ToWidgetter(buttons)...,
		),
		boxes.NewScrolledVbox(
			labels...,
		),
	)
	vbox.SetSpacing(1)
	vbox.SetHomogeneous(true)

	w.SetChild(vbox)
	w.Show()
}
