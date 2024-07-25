package widget

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

func NewButton(label string, onClicked func()) *gtk.Button {
	b := gtk.NewButtonWithLabel(label)
	b.ConnectClicked(func() {
		if onClicked != nil {
			onClicked()
		}
	})

	return b
}
