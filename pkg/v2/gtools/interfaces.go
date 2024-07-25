package gtools

import (
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
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

type Factoryer[T any] interface {
	SetSignalListItemFactory(*gtk.SignalListItemFactory)
	SignalListItemFactory() *gtk.SignalListItemFactory

	Modeller() Modeller[T]
	SetModeller(Modeller[T])

	SetSetup(FactorySetup)
	Setup() FactorySetup

	SetBind(FactoryBind[T])
	Bind() FactoryBind[T]
}

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

type (
	FactorySetup       func(ListItem)
	FactoryBind[T any] func(ListItem, T)
)

type SetModelFactoryer interface {
	gtk.Widgetter
	SetModel(gtk.SelectionModeller)
	SetFactory(*gtk.ListItemFactory)
}

type Selectioner interface {
	ListModeller() gio.ListModeller
	SetListModeller(gio.ListModeller)

	Autoselect() bool
	SetAutoselect(bool)

	CanUnselect() bool
	SetCanUnselect(bool)

	SelectionModeller() gtk.SelectionModeller
	SetSelectionModeller(gtk.SelectionModeller)

	SelectionMode() ListSelectionMode
	SetSelectionMode(ListSelectionMode)

	ConnectSelected(func(int))
	ConnectMultipleSelected(func([]int))

	SelectRange(int, int, bool)
	SelectAll()
	Select(int)
	SelectMultiple(...int)

	UnselectAll()
	UnselectRange(int, int)
	Unselect(int)

	Selected() int
	MultipleSelected() []int
}

type Modeller[T any] interface {
	ListModel() *gioutil.ListModel[T]
	S() []T
	Len() int
	SetItems([]T)
	Remove(int)
	Append(...T)
	Splice(int, int, ...T)
	RefreshModel()
	RefreshItems()
	At(int) T
	Range(func(int, T) bool)
}

type ListSelectionMode int

const (
	SelectionNone ListSelectionMode = iota
	SelectionSingle
	SelectionMultiple
)
