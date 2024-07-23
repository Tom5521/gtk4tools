package widgets

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type List[T any] struct {
	*gtk.ListView

	*ModelFactory[T, *gtk.ListView]
}

func NewList[T any](
	items []T,
	smodel ListSelectionMode,
	setup FactorySetup,
	bind FactoryBind[T],
) *List[T] {
	l := &List[T]{}
	l.ModelFactory = NewModelFactory[T, *gtk.ListView](
		smodel,
		gtk.NewListView(nil, nil),
		setup,
		bind,
		items...,
	)

	l.ListView = l.ModelFactory.Setter

	return l
}
