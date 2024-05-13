package main

import (
	"fmt"
	"os"

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
		Person{
			Name: "Jonh Doe 1",
			Age:  21,
		},
		Person{
			Name: "Carlos Gimenez",
			Age:  26,
		},
		Person{
			Name: "Caroline Simpson",
			Age:  20,
		},
	}

	list := widgets.NewList[Person](
		items,
		widgets.SelectionSingle,
		func(li *gtk.ListItem) {
			li.SetChild(gtk.NewLabel(""))
		},
		func(li *gtk.ListItem, p Person) {
			li.Child().(*gtk.Label).SetText(p.Name)
		},
	)
	list.OnSelected = func(index int) {
		fmt.Println("Index: ", index)
		fmt.Println("Value: ", items[index])
	}

	w.SetChild(list)
	w.Show()
}
