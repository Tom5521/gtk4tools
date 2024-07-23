package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Tom5521/gtk4tools/pkg/boxes"
	"github.com/Tom5521/gtk4tools/pkg/v2/gtools"
	"github.com/Tom5521/gtk4tools/pkg/v2/widget"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	app := gtk.NewApplication("com.test.window", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() {
		activate(app)
	})
	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	w := gtk.NewApplicationWindow(app)
	w.SetDefaultSize(500, 400)

	var buttons []*gtk.Button
	for i := range 30 {
		button := widget.NewButton("Button "+strconv.Itoa(i), func() {
			fmt.Println(i)
		})
		buttons = append(buttons, button)
	}

	var labels []gtk.Widgetter
	for i := range 30 {
		labels = append(labels, gtk.NewLabel("Label "+strconv.Itoa(i)))
	}

	items := func() []string {
		var out []string
		for i := range 100 {
			out = append(out, strconv.Itoa(i))
		}
		return out
	}()
	list := widget.NewList(
		items,
		gtools.SelectionMultiple,
		func(listitem gtools.ListItem) {
			listitem.SetChild(gtk.NewLabel(""))
		},
		func(listitem gtools.ListItem, item string) {
			label := listitem.Child().(*gtk.Label)
			label.SetText(item)
		},
	)
	list.SetVExpand(true)

	list.ConnectMultipleSelected(func(indexes []int) {
		for _, i := range indexes {
			fmt.Printf("|%s|", list.At(i))
		}
		fmt.Println()
	})
	list.ConnectSelected(func(index int) {
		fmt.Println(list.At(index))
	})

	vbox := boxes.NewCHbox(1,
		boxes.NewScrolledVbox(
			// Convert a slice of a specific type to a gtk.Widgetter slice.
			gtools.ToWidgetter(buttons)...,
		),
		boxes.NewScrolledVbox(
			labels...,
		),
		boxes.NewScrolledVbox(
			list,
		),
	)
	vbox.SetHomogeneous(true)

	go func() {
		for i := range 10 {
			time.Sleep(time.Second)
			list.Append("Hi" + strconv.Itoa(i))
		}
	}()

	w.SetChild(vbox)
	w.Present()
}
