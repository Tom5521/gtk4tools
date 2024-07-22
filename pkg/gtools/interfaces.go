package gtools

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ListItem interface {
	glib.Objector

	AccessibleDescription() string
	AccessibleLabel() string
	Activatable() bool
	Child() gtk.Widgetter
	Focusable() bool
	Item() *glib.Object
	Position() uint
	Selectable() bool
	Selected() bool

	SetAccessibleDescription(description string)
	SetAccessibleLabel(label string)
	SetActivatable(activatable bool)
	SetChild(child gtk.Widgetter)
	SetFocusable(focusable bool)
	SetSelectable(selectable bool)
}

var (
	_ ListItem = (*gtk.ListItem)(nil)
	_ ListItem = (*gtk.ColumnViewCell)(nil)
)
