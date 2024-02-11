package wcreate

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

func SliceToWidgets[T gtk.Widgetter](structs []T) []gtk.Widgetter {
	var widgets []gtk.Widgetter
	for _, b := range structs {
		widgets = append(widgets, gtk.Widgetter(b))
	}
	return widgets
}
