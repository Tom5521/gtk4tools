package widget

import (
	"github.com/Tom5521/gtk4tools/pkg/v2/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ListVar[T any] struct {
	*gtk.ListView
	*TemplateView[T, *gtk.ListView]
}

func NewListVar[T any](
	items *[]T,
	smodel gtools.ListSelectionMode,
	setup gtools.FactorySetup,
	bind gtools.FactoryBind[T],
) *ListVar[T] {
	l := &ListVar[T]{
		TemplateView: NewTemplateView[T](
			smodel,
			gtk.NewListView(nil, nil),
			setup,
			bind,
			gtools.NewModelVar(items),
		),
	}

	l.ListView = l.Setter

	return l
}
