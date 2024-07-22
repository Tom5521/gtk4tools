package gtools

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

var _ Appender = (*gtk.Box)(nil)

// It is simply an interface of a glib.Objector that has the append method.
type Appender interface {
	glib.Objector
	Append(gtk.Widgetter)
}

var (
	_ ListItem = (*gtk.ListItem)(nil)
	_ ListItem = (*gtk.ColumnViewCell)(nil)
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
