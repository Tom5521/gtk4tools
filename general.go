package gt4tools

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

func AppendWidgets(parent *gtk.Box, widgets ...gtk.Widgetter) {
	for _, w := range widgets {
		parent.Append(w)
	}
}
