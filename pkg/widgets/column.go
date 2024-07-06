package widgets

import (
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
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

	c.Factory.ConnectSetup(func(listitem *gtk.ListItem) {
		if c.Setup != nil {
			c.Setup(listitem)
		}
	})
	c.Factory.ConnectBind(func(listitem *gtk.ListItem) {
		if c.Bind != nil {
			c.Bind(listitem, c.Model.Item(int(listitem.Position())))
		}
	})

	c.ColumnViewColumn = gtk.NewColumnViewColumn(title, &c.Factory.ListItemFactory)

	return c
}
