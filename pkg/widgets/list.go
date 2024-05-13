package widgets

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ListSetup func(*gtk.ListItem)
type ListBind func(*gtk.ListItem, string)

type ListSelectionMode int

const (
	SelectionNone ListSelectionMode = iota
	SelectionSingle
	SelectionMultiple
)

// Deprecated: Replace it with CustomList[string].
type List struct {
	*CustomList[string]
}

// Creates a new list that keeps the self.Items[] updated with that of the UI.
//
// Deprecated: Replace it with NewCustomList[string]
// As of version 1.6.0 this is simply a wrapper for CustomList[string],
// so it is recommended to simply change it to this type.
// There should be no further incompatibility with respect to this migration,
// so it should not be a problem.
//
// Since version 1.6.1 this struct is renamed to StringList.
func NewList(
	items []string,
	smodel ListSelectionMode,
	setup ListSetup,
	bind ListBind,
) *List {
	return &List{
		CustomList: NewCustomList(
			items,
			smodel,
			setup,
			CustomListBind[string](bind),
		),
	}
}
