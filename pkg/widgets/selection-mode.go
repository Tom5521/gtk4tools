package widgets

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type ListSelectionMode int

const (
	SelectionNone ListSelectionMode = iota
	SelectionSingle
	SelectionMultiple
)

func (m ListSelectionMode) Modeller() gtk.SelectionModeller {
	var modeller gtk.SelectionModeller
	switch m {
	case SelectionSingle:
		modeller = gtk.NewSingleSelection(nil)
	case SelectionMultiple:
		modeller = gtk.NewMultiSelection(nil)
	default:
		modeller = gtk.NewNoSelection(nil)
	}

	return modeller
}
