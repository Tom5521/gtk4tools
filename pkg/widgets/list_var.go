package widgets

import (
	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ListVar[T any] struct {
	*gtk.ListView
	*ModelFactory[T, *gtk.ListView]
}

func NewListVar[T any](
	items *[]T,
	smodel ListSelectionMode,
	setup FactorySetup,
	bind FactoryBind[T],
) *ListVar[T] {
	l := &ListVar[T]{
		ModelFactory: NewModelFactory(
			smodel,
			gtk.NewListView(nil, nil),
			setup,
			bind,
			NewModelVar(items),
		),
	}

	l.ListView = l.Setter

	return l
}

// Internal functions

func (l *ListVar[T]) reConnectFactory() {
	l.ItemFactory.ConnectSetup(gtools.NewFactorySetup(func(listitem gtools.ListItem) {
		if l.Setup == nil {
			return
		}
		l.Setup(listitem)
	}))
	l.ItemFactory.ConnectBind(gtools.NewFactoryBind(func(listitem gtools.ListItem, pos int) {
		if l.Bind == nil {
			return
		}
		l.Bind(listitem, l.At(pos))
	}))
}
