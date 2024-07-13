package widgets

import (
	"slices"

	"github.com/Tom5521/gtk4tools/pkg/tools"
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ListVar[T any] struct {
	*List[T]

	Items *[]T
}

func NewListVar[T any](
	items *[]T,
	smodel ListSelectionMode,
	setup FactorySetup,
	bind ListBind[T],
) *ListVar[T] {
	l := &ListVar[T]{
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
func (l *ListVar[T]) SetItems(items *[]T) {
	l.Splice(0, l.Model.Len(), *items...)
}

func (l *ListVar[T]) Remove(index int) {
	if index <= -1 {
		return
	}
	*l.Items = slices.Delete(*l.Items, index, index+1)
	l.Model.Remove(index)
}

func (l *ListVar[T]) Append(item T) {
	*l.Items = append(*l.Items, item)
	l.Model.Append(item)
}

func (l *ListVar[T]) Splice(pos, nRemovals int, additions ...T) {
	if pos <= -1 || nRemovals <= -1 {
		return
	}
	l.Model.Splice(pos, nRemovals, additions...)
	l.RefreshItems()
}

func (l *ListVar[T]) RefreshItems() {
	*l.Items = []T{}
	l.Model.All()(func(v T) bool {
		*l.Items = append(*l.Items, v)
		return true
	})
}

// Internal functions

func (l *ListVar[T]) reConnectFactory() {
	l.Factory.ConnectSetup(tools.NewFactorySetup(func(listitem *gtk.ListItem) {
		if l.Setup == nil {
			return
		}
		l.Setup(listitem)
	}))
	l.Factory.ConnectBind(tools.NewFactoryBind(func(listitem *gtk.ListItem, pos int) {
		if l.Bind == nil {
			return
		}
		l.Bind(listitem, (*l.Items)[pos])
	}))
}

func (l *ListVar[T]) RefreshModel() {
	if l.Model.Len() == 0 {
		for _, i := range *l.Items {
			l.Model.Append(i)
		}
	}
	l.Model.Splice(0, l.Model.Len(), *l.Items...)
}
