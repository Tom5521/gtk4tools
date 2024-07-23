package widget_test

import (
	"fmt"

	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/Tom5521/gtk4tools/pkg/widget"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleList() {
	items := []string{"1", "2", "3"}
	list := widget.NewList(
		items,
		gtools.SelectionMultiple,
		func(listitem gtools.ListItem) {
			listitem.SetChild(gtk.NewLabel(""))
		},
		func(listitem gtools.ListItem, obj string) {
			listitem.Child().(*gtk.Label).SetText(obj)
		},
	)

	list.ConnectMultipleSelected(func(indexes []int) {
		for _, i := range indexes {
			fmt.Printf("|%s|", list.At(i))
		}
		fmt.Println()
	})
	list.ConnectSelected(func(index int) {
		fmt.Println(list.At(index))
	})
}

func ExampleList_RefreshFactory() {
	items := []string{"1", "2", "3"}
	list := widget.NewList(
		items,
		gtools.SelectionMultiple,
		func(listitem gtools.ListItem) {
			listitem.SetChild(gtk.NewLabel(""))
		},
		func(listitem gtools.ListItem, obj string) {
			listitem.Child().(*gtk.Label).SetText(obj)
		},
	)

	list.Setup = func(li gtools.ListItem) {
		li.SetChild(gtk.NewText())
	}
	list.Bind = func(li gtools.ListItem, s string) {
		li.Child().(*gtk.Text).SetText(s)
	}
	list.RefreshFactory()
}
