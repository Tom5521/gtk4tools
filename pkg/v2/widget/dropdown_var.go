package widget

import (
	"github.com/Tom5521/gtk4tools/pkg/v2/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type DropDownVar[T any] struct {
	gtk.Widgetter
	gtools.Factoryer[T]
	*gtools.ModelVar[T]

	OnSelected func(T)

	dropDown *gtk.DropDown
}

func NewDropDownVar[T any](
	items *[]T,
	setup gtools.FactorySetup,
	bind gtools.FactoryBind[T],
) *DropDownVar[T] {
	d := &DropDownVar[T]{
		ModelVar: gtools.NewModelVar(items),
	}

	d.Factoryer = gtools.NewFactory[T](
		d.ModelVar,
		setup,
		bind,
	)

	d.dropDown = gtk.NewDropDown(d.ModelVar.ListModel(), nil)
	d.dropDown.SetFactory(&d.Factoryer.Factory().ListItemFactory)
	d.dropDown.ConnectAfter("notify::selected", func() {
		if d.OnSelected != nil {
			d.OnSelected(d.At(int(d.dropDown.Selected())))
		}
	})

	d.Widgetter = d.dropDown
	return d
}
