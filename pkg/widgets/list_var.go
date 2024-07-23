package widgets

import (
	"slices"

	"github.com/Tom5521/gtk4tools/pkg/gtools"
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
	bind FactoryBind[T],
) *ListVar[T] {
	l := &ListVar[T]{
		Items: items,
		List: &List[T]{
			ModelFactory: &ModelFactory[T, *gtk.ListView]{
				Setup:         setup,
				Bind:          bind,
				SelectionMode: smodel,
				ItemFactory:   gtk.NewSignalListItemFactory(),
				Setter:        gtk.NewListView(nil, nil),
				Model:         NewModel(*items...),
			},
		},
	}

	l.ListView = l.Setter

	l.reConnectFactory()
	l.makeSelectionModeller(smodel)
	l.reConnectSelection()

	l.InitSetter()

	return l
}

// Re-generate the list with the items provided.
func (l *ListVar[T]) SetItems(items *[]T) {
	l.Splice(0, l.ListModel.Len(), *items...)
}

func (l *ListVar[T]) Remove(index int) {
	if index <= -1 {
		return
	}
	*l.Items = slices.Delete(*l.Items, index, index+1)
	l.ListModel.Remove(index)
}

func (l *ListVar[T]) Append(item T) {
	*l.Items = append(*l.Items, item)
	l.ListModel.Append(item)
}

func (l *ListVar[T]) Splice(pos, nRemovals int, additions ...T) {
	if pos <= -1 || nRemovals <= -1 {
		return
	}
	l.ListModel.Splice(pos, nRemovals, additions...)
	l.RefreshItems()
}

func (l *ListVar[T]) RefreshItems() {
	*l.Items = []T{}
	l.ListModel.All()(func(v T) bool {
		*l.Items = append(*l.Items, v)
		return true
	})
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
		l.Bind(listitem, (*l.Items)[pos])
	}))
}

func (l *ListVar[T]) RefreshModel() {
	if l.ListModel.Len() == 0 {
		for _, i := range *l.Items {
			l.ListModel.Append(i)
		}
	}
	l.ListModel.Splice(0, l.ListModel.Len(), *l.Items...)
}
