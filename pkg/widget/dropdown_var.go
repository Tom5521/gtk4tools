//go:build dev

package widget

import (
	"slices"

	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type DropDownVar[T any] struct {
	*DropDown[T]

	Items *[]T
}

func NewDropDownVar[T any](
	items *[]T,
	setup FactorySetup,
	bind FactoryBind[T],
) *DropDownVar[T] {
	d := &DropDownVar[T]{
		DropDown: &DropDown[T]{
			Model:   gioutil.NewListModel[T](),
			Factory: gtk.NewSignalListItemFactory(),

			Setup: setup,
			Bind:  bind,
		},
		Items: items,
	}

	d.RefreshModel()
	d.connectFactory()

	d.DropDown.DropDown = gtk.NewDropDown(d.Model.ListModel, nil)
	d.SetFactory(&d.Factory.ListItemFactory)

	d.connectChanged()

	return d
}

func (d *DropDownVar[T]) Append(v T) {
	*d.Items = append(*d.Items, v)
	d.Model.Append(v)
}

func (d *DropDownVar[T]) Remove(index int) {
	*d.Items = slices.Delete(*d.Items, index, index+1)
	d.Model.Remove(index)
}

func (d *DropDownVar[T]) Splice(pos, rms int, values ...T) {
	spliceVar(d.Items, pos, rms, values)
	d.Model.Splice(pos, rms, values...)
}

func (d *DropDownVar[T]) RefreshModel() {
	if d.Model.NItems() == 0 {
		for _, i := range *d.Items {
			d.Model.Append(i)
		}
		return
	}
	d.Model.Splice(0, 0, *d.Items...)
}

func (d *DropDownVar[T]) connectFactory() {
	d.Factory.ConnectSetup(gtools.NewFactorySetup(func(listitem gtools.ListItem) {
		if d.Setup == nil {
			return
		}
		d.Setup(listitem)
	}))
	d.Factory.ConnectBind(gtools.NewFactoryBind(func(listitem gtools.ListItem, pos int) {
		if d.Bind == nil {
			return
		}
		d.Bind(listitem, (*d.Items)[pos])
	}))
}
