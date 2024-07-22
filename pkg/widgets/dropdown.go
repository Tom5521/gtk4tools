package widgets

import (
	"slices"

	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type DropDown[T any] struct {
	*gtk.DropDown

	Items   []T
	Model   *gioutil.ListModel[T]
	Factory *gtk.SignalListItemFactory

	OnChanged func(index int)

	Setup FactorySetup
	Bind  ListBind[T]
}

func NewDropDown[T any](
	items []T,
	setup FactorySetup,
	bind ListBind[T],
) *DropDown[T] {
	d := &DropDown[T]{
		Items:   items,
		Model:   gioutil.NewListModel[T](),
		Factory: gtk.NewSignalListItemFactory(),

		Setup: setup,
		Bind:  bind,
	}

	d.RefreshModel()
	d.connectFactory()

	d.DropDown = gtk.NewDropDown(d.Model.ListModel, nil)
	d.SetFactory(&d.Factory.ListItemFactory)

	d.connectChanged()

	return d
}

func (d *DropDown[T]) Append(v T) {
	d.Items = append(d.Items, v)
	d.Model.Append(v)
}

func (d *DropDown[T]) Remove(index int) {
	d.Items = slices.Delete(d.Items, index, index+1)
	d.Model.Remove(index)
}

func (d *DropDown[T]) Splice(pos, rms int, values ...T) {
	spliceVar(&d.Items, pos, rms, values)
	d.Model.Splice(pos, rms, values...)
}

func (d *DropDown[T]) RefreshModel() {
	d.Model.Splice(0, d.Model.Len(), d.Items...)
}

// Private methods.

func (d *DropDown[T]) connectChanged() {
	d.ConnectAfter("notify::selected", func() {
		if d.OnChanged == nil {
			return
		}
		d.OnChanged(int(d.Selected()))
	})
}

func (d *DropDown[T]) connectFactory() {
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
		d.Bind(listitem, d.Model.At(pos))
	}))
}
