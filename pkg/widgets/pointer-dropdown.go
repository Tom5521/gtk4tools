package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type PointerDropDown[T any] struct {
	*DropDown[T]

	Items *[]T
}

func NewPointerDropDown[T any](
	items *[]T,
	setup FactorySetup,
	bind ListBind[T],
) *PointerDropDown[T] {
	d := &PointerDropDown[T]{
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

func (d *PointerDropDown[T]) Append(v T) {
	*d.Items = append(*d.Items, v)
	d.Model.Append(v)
}

func (d *PointerDropDown[T]) Remove(index int) {
	*d.Items = slices.Delete(*d.Items, index, index+1)
	d.Model.Remove(index)
}

func (d *PointerDropDown[T]) Splice(pos, rms int, values ...T) {
	*d.Items = splice(*d.Items, pos, rms, values)
	d.Model.Splice(pos, rms, values...)
}

func (d *PointerDropDown[T]) RefreshModel() {
	if d.Model.NItems() == 0 {
		for _, i := range *d.Items {
			d.Model.Append(i)
		}
		return
	}
	d.Model.Splice(0, 0, *d.Items...)
}

func (d *PointerDropDown[T]) connectFactory() {
	d.Factory.ConnectSetup(func(listitem *gtk.ListItem) {
		if d.Setup == nil {
			return
		}
		d.Setup(listitem)
	})
	d.Factory.ConnectBind(func(listitem *gtk.ListItem) {
		if d.Bind == nil {
			return
		}
		d.Bind(listitem, (*d.Items)[listitem.Position()])
	})
}
