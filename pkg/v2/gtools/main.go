package gtools

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// Converts the slice of the gtk.Widgetter types to a slice of gtk.Widgetter.
func ToWidgetter[T gtk.Widgetter](items []T) []gtk.Widgetter {
	var widgets []gtk.Widgetter
	for _, w := range items {
		widgets = append(widgets, w)
	}
	return widgets
}

func Append(parent Appender, widgets ...gtk.Widgetter) {
	for _, w := range widgets {
		parent.Append(w)
	}
}
