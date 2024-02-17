package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

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
	l.Splice(0, l.Model.NItems(), items...)
}

func (l *List) Remove(index uint) {
	l.Items = slices.Delete(l.Items, int(index), int(index+1))
	l.Model.Remove(index)
}

func (l *List) Append(item string) {
	l.Items = append(l.Items, item)
	l.Model.Append(item)
}

func (l *List) Splice(pos, nRemovals uint, additions ...string) {
	l.Items = slices.Delete(l.Items, int(pos), int(nRemovals))
	l.Items = append(l.Items, additions...)
	l.Model.Splice(pos, nRemovals, additions)
}

func (l *List) ConnectSelected(f func(index uint)) {
	type selecter interface {
		gtk.SelectionModeller
		Selected() uint
	}
	model, ok := l.SelectionModeller.(selecter)
	if !ok {
		return
	}

	model.ConnectSelectionChanged(func(_, _ uint) {
		f(model.Selected())
	})
}
