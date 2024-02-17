package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type selecter interface {
	gtk.SelectionModeller
	Selected() uint
	Item(uint) *glib.Object
}

type SelectionMode int

const (
	SelectionNone SelectionMode = iota
	SelectionSingle
	SelectionMultiple
)

type List struct {
	*gtk.ListView

	Items []string

	SelectionModeller gtk.SelectionModeller
	SelectionMode     SelectionMode
	Model             *gtk.StringList
	Factory           *gtk.SignalListItemFactory
}

func NewList(items []string, smodel SelectionMode, setup, bind func(listitem *gtk.ListItem)) *List {
	l := &List{
		Items:         items,
		SelectionMode: smodel,
		Model:         gtk.NewStringList(items),
		Factory:       gtk.NewSignalListItemFactory(),
	}
	l.Factory.ConnectSetup(setup)
	l.Factory.ConnectBind(bind)

	l.SetSelectionModeller(smodel)

	l.ListView = gtk.NewListView(l.SelectionModeller, &l.Factory.ListItemFactory)

	return l
}

func (l *List) SetSelectionModeller(mode SelectionMode) {
	switch mode {
	case SelectionNone:
		l.SelectionModeller = gtk.NewNoSelection(l.Model)
	case SelectionSingle:
		l.SelectionModeller = gtk.NewSingleSelection(l.Model)
	case SelectionMultiple:
		l.SelectionModeller = gtk.NewMultiSelection(l.Model)
	default:
		l.SelectionModeller = gtk.NewNoSelection(l.Model)
	}
}

func (l *List) SetItems(items []string) {
	l.Splice(0, int(l.Model.NItems()), items...)
}

func (l *List) Remove(index int) {
	if index <= -1 {
		return
	}
	l.Items = slices.Delete(l.Items, index, index+1)
	l.Model.Remove(uint(index))
}

func (l *List) Append(item string) {
	l.Items = append(l.Items, item)
	l.Model.Append(item)
}

func (l *List) Splice(pos, nRemovals int, additions ...string) {
	if pos <= -1 || nRemovals <= -1 {
		return
	}
	l.Items = slices.Delete(l.Items, pos, nRemovals)
	l.Items = append(l.Items, additions...)
	l.Model.Splice(uint(pos), uint(nRemovals), additions)
}

func (l *List) ConnectSelected(f func(index int)) {
	l.SelectionModeller.ConnectSelectionChanged(func(_, _ uint) {
		f(l.Selected())
	})
}

func (l *List) Selected() int {
	model, ok := l.SelectionModeller.(selecter)
	if !ok {
		return -1
	}
	i := model.Selected()
	if model.Item(i) == nil {
		return -1
	}
	return int(i)
}
