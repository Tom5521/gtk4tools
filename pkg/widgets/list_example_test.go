package widgets

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleList() {
	items := []string{"1", "2", "3"}
	list := NewList(
		items,
		SelectionMultiple,
		func(listitem *gtk.ListItem) {
			listitem.SetChild(gtk.NewLabel(""))
		},
		func(listitem *gtk.ListItem) {
			label := listitem.Child().(*gtk.Label)
			obj := listitem.Item().Cast().(*gtk.StringObject)
			label.SetText(obj.String())
		},
	)

	list.OnMultipleSelected = func(indexes []int) {
		fmt.Println(indexes)
	}
	list.OnSelected = func(index int) {
		fmt.Println(index)
	}
}
