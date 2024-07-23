package widgets

import (
	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type (
	FactorySetup       func(gtools.ListItem)
	FactoryBind[T any] func(gtools.ListItem, T)
)

type Setter interface {
	gtk.Widgetter
	SetModel(gtk.SelectionModeller)
	SetFactory(*gtk.ListItemFactory)
}

type ListSelectionMode int

const (
	SelectionNone ListSelectionMode = iota
	SelectionSingle
	SelectionMultiple
)
