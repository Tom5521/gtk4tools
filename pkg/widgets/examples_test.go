package widgets_test

import (
	"fmt"

	"github.com/Tom5521/gtk4tools/pkg/widgets"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleList() {
	items := []string{"1", "2", "3"}
	list := widgets.NewList(
		items,
		widgets.SelectionMultiple,
		func(listitem *gtk.ListItem) {
			listitem.SetChild(gtk.NewLabel(""))
		},
		func(listitem *gtk.ListItem, obj string) {
			listitem.Child().(*gtk.Label).SetText(obj)
		},
	)

	list.OnMultipleSelected = func(indexes []int) {
		for _, i := range indexes {
			fmt.Printf("|%s|", list.Items[i])
		}
		fmt.Println()
	}
	list.OnSelected = func(index int) {
		fmt.Println(list.Items[index])
	}
}

func ExampleList_RefreshFactory() {
	items := []string{"1", "2", "3"}
	list := widgets.NewList(
		items,
		widgets.SelectionMultiple,
		func(listitem *gtk.ListItem) {
			listitem.SetChild(gtk.NewLabel(""))
		},
		func(listitem *gtk.ListItem, obj string) {
			listitem.Child().(*gtk.Label).SetText(obj)
		},
	)

	list.Setup = func(li *gtk.ListItem) {
		li.SetChild(gtk.NewText())
	}
	list.Bind = func(li *gtk.ListItem, s string) {
		li.Child().(*gtk.Text).SetText(s)
	}
	list.RefreshFactory()
}
