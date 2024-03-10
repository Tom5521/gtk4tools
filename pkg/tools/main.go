package tools

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

// It is simply an interface of a gtk.widgetter that has the append function.
type Appender interface {
	gtk.Widgetter
	Append(gtk.Widgetter)
}

func Append(parent Appender, widgets ...gtk.Widgetter) {
	for _, w := range widgets {
		parent.Append(w)
	}
}
