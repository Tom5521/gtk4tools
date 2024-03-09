package boxes_test

import (
	"github.com/Tom5521/gtk4tools/pkg/boxes"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleNewNoteBook() {
	w := gtk.NewWindow()

	notebook := boxes.NewNoteBook(
		&boxes.NotebookTab{
			Child: boxes.NewVbox(
				gtk.NewLabel("Hello World"),
			),
			Label: gtk.NewLabel("Tab1"),
		},
		&boxes.NotebookTab{
			Child: boxes.NewVbox(
				gtk.NewLabel("Hello World in tab 2"),
			),
			Label: gtk.NewLabel("Tab2"),
		},
	)

	w.SetChild(notebook)
	w.Show()
}
