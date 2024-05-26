package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type PointerList[T any] struct {
	*List[T]

	Items *[]T
}

func NewPointerList[T any](
	items *[]T,
	smodel ListSelectionMode,
	setup FactorySetup,
	bind ListBind[T],
) *PointerList[T] {
	l := &PointerList[T]{
		Items: items,
		List: &List[T]{
			Setup:         setup,
			Bind:          bind,
			SelectionMode: smodel,
			Factory:       gtk.NewSignalListItemFactory(),
			Model:         gioutil.NewListModel[T](),
		},
	}

	l.RefreshModel()
	l.reConnectFactory()
	l.makeSelectionModeller(smodel)
	l.reConnectSelection()

	l.ListView = gtk.NewListView(l.SelectionModeller, &l.Factory.ListItemFactory)

	return l
}

// Re-generate the list with the items provided.
func (l *PointerList[T]) SetItems(items *[]T) {
	l.Splice(0, l.Model.NItems(), *items...)
}

func (l *PointerList[T]) Remove(index int) {
	if index <= -1 {
		return
	}
	*l.Items = slices.Delete(*l.Items, index, index+1)
	l.Model.Remove(index)
}

func (l *PointerList[T]) Append(item T) {
	*l.Items = append(*l.Items, item)
	l.Model.Append(item)
}

func (l *PointerList[T]) Splice(pos, nRemovals int, additions ...T) {
	if pos <= -1 || nRemovals <= -1 {
		return
	}
	*l.Items = slices.Delete(*l.Items, pos, nRemovals)
	*l.Items = append(*l.Items, additions...)
	l.Model.Splice(pos, nRemovals, additions...)
}

func (l *PointerList[T]) RefreshItems() {
	*l.Items = []T{}
	for i := range l.Model.NItems() {
		*l.Items = append(*l.Items, l.Model.Item(i))
	}
}

// Internal functions

func (l *PointerList[T]) reConnectFactory() {
	l.Factory.ConnectSetup(func(listitem *gtk.ListItem) {
		if l.Setup == nil {
			return
		}
		l.Setup(listitem)
	})
	l.Factory.ConnectBind(func(listitem *gtk.ListItem) {
		if l.Bind == nil {
			return
		}
		l.Bind(listitem, (*l.Items)[listitem.Position()])
	})
}

func (l *PointerList[T]) RefreshModel() {
	if l.Model.NItems() == 0 {
		for _, i := range *l.Items {
			l.Model.Append(i)
		}
	}
	l.Model.Splice(0, l.Model.NItems(), *l.Items...)
}
