package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type List struct {
	*gtk.ListView

	Items []string

	SelectionModel gtk.SelectionModeller
	Model          *gtk.StringList
	Factory        *gtk.SignalListItemFactory
}

func NewList(items []string, smodel gtk.SelectionModeller, setup, bind func(listitem *gtk.ListItem)) *List {
	l := &List{
		Items:          items,
		SelectionModel: smodel,
		Model:          gtk.NewStringList(items),
		Factory:        gtk.NewSignalListItemFactory(),
	}
	l.Factory.ConnectSetup(setup)
	l.Factory.ConnectBind(bind)

	return l
}

func (l *List) SetItems(items ...string) {
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
	model, ok := l.SelectionModel.(selecter)
	if !ok {
		return
	}

	f(model.Selected())
}
