package widgets

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type List struct {
	*gtk.ListView

	SelectionModel gtk.SelectionModeller
	Model          *gtk.StringList
	Factory        *gtk.SignalListItemFactory
}

func NewList(items []string, smodel gtk.SelectionModeller, setup, bind func(listitem *gtk.ListItem)) *List {
	l := &List{
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
	l.Model.Remove(index)
}

func (l *List) Append(item string) {
	l.Model.Append(item)
}

func (l *List) Splice(pos, nRemovals uint, additions ...string) {
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
