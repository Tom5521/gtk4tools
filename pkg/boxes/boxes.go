package boxes

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

var DefaultSpacing = 4

// Creates a new *gtk.Box with vertical orientation.
func NewVbox(widgets ...gtk.Widgetter) *gtk.Box {
	vbox := gtk.NewBox(gtk.OrientationVertical, DefaultSpacing)
	for _, w := range widgets {
		vbox.Append(w)
	}

	return vbox
}

// Creates a new *gtk.Box with horizontal orientation.
func NewHbox(widgets ...gtk.Widgetter) *gtk.Box {
	hbox := gtk.NewBox(gtk.OrientationHorizontal, DefaultSpacing)
	for _, w := range widgets {
		hbox.Append(w)
	}
	return hbox
}

// Creates a *gtk.Paned with vertical orientation.
func NewVPaned(top, bottom gtk.Widgetter) *gtk.Paned {
	paned := gtk.NewPaned(gtk.OrientationVertical)
	paned.SetStartChild(top)
	paned.SetEndChild(bottom)
	return paned
}

// Creates a *gtk.Paned with horizontal orientation
func NewHPaned(left, right gtk.Widgetter) *gtk.Paned {
	paned := gtk.NewPaned(gtk.OrientationHorizontal)
	paned.SetStartChild(left)
	paned.SetEndChild(right)
	return paned
}

// Creates a new *gtk.Frame
func NewFrame(label string, child gtk.Widgetter) *gtk.Frame {
	frame := gtk.NewFrame(label)
	frame.SetChild(child)

	return frame
}

// Creates a grid that adapts to the specified size, items are fixed,
// they have to be added manually to the grid if you want to append them.
func NewAdaptativeGrid(size int, widgets ...gtk.Widgetter) *gtk.Grid {
	grid := gtk.NewGrid()

	var rowCount, columnCount int
	for _, w := range widgets {
		visible, ok := w.ObjectProperty("visible").(bool)
		if ok {
			if !visible {
				continue
			}
		}
		grid.Attach(w, columnCount, rowCount, 1, 1)
		if columnCount >= size {
			rowCount++
			columnCount = 0
			continue
		}
		columnCount++
	}

	return grid
}

type NotebookTab struct {
	Label gtk.Widgetter
	Child gtk.Widgetter
}

// Create a new notebook with the specified NotebookTabs.
func NewNoteBook(tabs ...NotebookTab) *gtk.Notebook {
	n := gtk.NewNotebook()
	for _, t := range tabs {
		n.AppendPage(t.Child, t.Label)
	}

	return n
}
