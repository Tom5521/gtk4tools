package widget

import (
	"github.com/Tom5521/gtk4tools/pkg/v2/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type List[T any] struct {
	*gtk.ListView

	*TemplateView[T, *gtk.ListView]
}

func NewList[T any](
	items []T,
	smodel gtools.ListSelectionMode,
	setup gtools.FactorySetup,
	bind gtools.FactoryBind[T],
) *List[T] {
	l := &List[T]{}
	l.TemplateView = NewTemplateView[T, *gtk.ListView](
		smodel,
		gtk.NewListView(nil, nil),
		setup,
		bind,
		gtools.NewModel(items...),
	)

	l.ListView = l.TemplateView.Setter

	return l
}
