package widgets

import (
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Column[T any] struct {
	*gtk.ColumnViewColumn
	Title string

	Factory *gtk.SignalListItemFactory

	Model *gioutil.ListModel[T]

	Setup FactorySetup
	Bind  ListBind[T]
}

func NewColumn[T any](
	title string,
	model *gioutil.ListModel[T],
	setup FactorySetup,
	bind ListBind[T],
) *Column[T] {
	c := &Column[T]{
		Title: title,
		Model: model,
		Setup: setup,
		Bind:  bind,

		// Init widgets.
		Factory: gtk.NewSignalListItemFactory(),
	}

	c.Factory.ConnectSetup(func(obj *glib.Object) {
		if c.Setup != nil {
			listitem := obj.Cast().(*gtk.ListItem)
			c.Setup(listitem)
		}
	})
	c.Factory.ConnectBind(func(obj *glib.Object) {
		if c.Bind != nil {
			listitem := obj.Cast().(*gtk.ListItem)
			item := c.Model.Item(listitem.Position())
			c.Bind(listitem, item.Cast().(T))
		}
	})

	c.ColumnViewColumn = gtk.NewColumnViewColumn(title, &c.Factory.ListItemFactory)

	return c
}
