package widgets

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type FactorySetup func(*gtk.ListItem)

type ListSelectionMode int

const (
	SelectionNone ListSelectionMode = iota
	SelectionSingle
	SelectionMultiple
)
