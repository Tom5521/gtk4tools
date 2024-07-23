package widgets

import (
	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type (
	FactorySetup       func(gtools.ListItem)
	FactoryBind[T any] func(gtools.ListItem, T)
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

	SetSelectionModeller(gtk.SelectionModeller)
	SelectionModeller() gtk.SelectionModeller

	SetSelectionMode(ListSelectionMode)
	SelectionMode() ListSelectionMode

	ConnectSelected(func(int))
	ConnectMultipleSelected(func([]int))

	SelectRange(int, int, bool)
	SelectAll()
	UnselectAll()
	UnselectRange(int, int)
	Unselect(int)
	Selected() int
	Select(int)
	MultipleSelected() []int
	SelectMultiple(...int)
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
