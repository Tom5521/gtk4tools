package t

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

// Converts the slice of the gtk.Widgetter types to a slice of gtk.Widgetter.
func ToWidgetter[T gtk.Widgetter](items []T) []gtk.Widgetter {
	var ws []gtk.Widgetter
	for _, w := range items {
		ws = append(ws, w)
	}
	return ws
}
