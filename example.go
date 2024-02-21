package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Tom5521/gtk4tools/pkg/boxes"
	t "github.com/Tom5521/gtk4tools/pkg/tools"
	"github.com/Tom5521/gtk4tools/pkg/widgets"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

var VV string

func main() {
	fmt.Println(VV)
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

	items := func() []string {
		var out []string
		for i := range 100 {
			out = append(out, strconv.Itoa(i))
		}
		return out
	}()
	list := widgets.NewList(
		items,
		widgets.SelectionMultiple,
		func(listitem *gtk.ListItem) {
			listitem.SetChild(gtk.NewLabel(""))
		},
		func(listitem *gtk.ListItem) {
			label := listitem.Child().(*gtk.Label)
			obj := listitem.Item().Cast().(*gtk.StringObject)
			label.SetText(obj.String())
		},
	)

	list.SetVExpand(true)

	list.OnMultipleSelected = func(indexes []int) {
		for _, i := range indexes {
			fmt.Printf("|%s|", list.Items[i])
		}
		fmt.Println()
	}
	list.OnSelected = func(index int) {
		fmt.Println(list.Items[index])
	}

	vbox := boxes.NewHbox(
		boxes.NewScrolledVbox(
			// Convert a slice of a specific type to a gtk.Widgetter slice.
			t.ToWidgetter(buttons)...,
		),
		boxes.NewScrolledVbox(
			labels...,
		),
		boxes.NewScrolledVbox(
			list,
		),
	)
	vbox.SetSpacing(1)
	vbox.SetHomogeneous(true)

	w.SetChild(vbox)
	w.Show()
}
