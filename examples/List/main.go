package main

import (
	"fmt"
	"os"

	"github.com/Tom5521/gtk4tools/pkg/boxes"
	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/Tom5521/gtk4tools/pkg/widgets"
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

type Person struct {
	Name string
	Age  uint
}

func activate(app *gtk.Application) {
	w := gtk.NewApplicationWindow(app)

	items := []Person{
		{
			Name: "Jonh Doe 1",
			Age:  21,
		},
		{
			Name: "Carlos Gimenez",
			Age:  26,
		},
		{
			Name: "Caroline Simpson",
			Age:  20,
		},
	}

	list := widgets.NewList[Person](
		items,
		widgets.SelectionSingle,
		func(li gtools.ListItem) {
			li.SetChild(gtk.NewLabel(""))
		},
		func(li gtools.ListItem, p Person) {
			li.Child().(*gtk.Label).SetText(p.Name)
		},
	)

	list.OnSelected = func(index int) {
		fmt.Println("Index: ", index)
		fmt.Println("Value: ", items[index])
	}
	list.OnMultipleSelected = func(indexes []int) {
		fmt.Println("Indexes: ", indexes)
		fmt.Print("Values: ")
		var values []Person
		for ci, v := range items {
			for _, i := range indexes {
				if i == ci {
					values = append(values, v)
				}
			}
		}
		fmt.Println(values)
	}

	list.SetVExpand(true)

	button := widgets.NewButton("Change Selection Model", func() {
		switch list.SelectionMode {
		case widgets.SelectionNone:
			list.SetSelectionModeller(widgets.SelectionSingle)
		case widgets.SelectionSingle:
			list.SetSelectionModeller(widgets.SelectionMultiple)
		case widgets.SelectionMultiple:
			list.SetSelectionModeller(widgets.SelectionNone)
		}
	})

	w.SetChild(boxes.NewVbox(
		list,
		button,
	))
	w.Present()
}
