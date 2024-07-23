package widget

import (
	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type DropDown[T any] struct {
	gtk.Widgetter
	gtools.Factoryer[T]
	*gtools.Model[T]

	OnSelected func(T)

	dropDown *gtk.DropDown
}

func NewDropDown[T any](
	items []T,
	setup gtools.FactorySetup,
	bind gtools.FactoryBind[T],
) *DropDown[T] {
	d := &DropDown[T]{
		Model: gtools.NewModel(items...),
	}

	d.Factoryer = gtools.NewFactory[T](
		d.Model,
		setup,
		bind,
	)

	d.dropDown = gtk.NewDropDown(d.Model.ListModel(), nil)
	d.dropDown.SetFactory(&d.Factoryer.Factory().ListItemFactory)
	d.dropDown.ConnectAfter("notify::selected", func() {
		if d.OnSelected != nil {
			d.OnSelected(d.At(int(d.dropDown.Selected())))
		}
	})

	d.Widgetter = d.dropDown
	return d
}
