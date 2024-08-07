package main

import (
	"fmt"
	"math/rand/v2"
	"os"

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

	list := widget.NewListVar[Person](
		&items,
		gtools.SelectionSingle,
		func(li gtools.ListItem) {
			li.SetChild(gtk.NewLabel(""))
		},
		func(li gtools.ListItem, p Person) {
			li.Child().(*gtk.Label).SetText(p.Name)
		},
	)

	list.ConnectSelected(func(index int) {
		fmt.Println("Index: ", index)
		fmt.Println("Value: ", items[index])
	})
	list.ConnectMultipleSelected(func(indexes []int) {
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
	})

	list.SetVExpand(true)

	button1 := widget.NewButton("Change Selection Model", func() {
		switch list.SelectionMode() {
		case gtools.SelectionNone:
			list.SetSelectionMode(gtools.SelectionSingle)
		case gtools.SelectionSingle:
			list.SetSelectionMode(gtools.SelectionMultiple)
		case gtools.SelectionMultiple:
			list.SetSelectionMode(gtools.SelectionNone)
		}
	})

	entry := gtk.NewEntry()
	button2 := widget.NewButton("Append item", func() {
		if entry.Text() == "" {
			return
		}
		list.Append(Person{
			Name: entry.Text(),
			Age:  uint(rand.IntN(100)),
		})
		entry.SetText("")
	})

	entryBox := boxes.NewHbox(
		button2,
		entry,
	)
	entryBox.SetHomogeneous(true)

	buttonsBox := boxes.NewHbox(
		button1,
		entryBox,
	)
	buttonsBox.SetHomogeneous(true)

	w.SetChild(boxes.NewVbox(
		list,
		buttonsBox,
	))
	w.Present()
}
