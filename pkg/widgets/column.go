package widgets

import (
	"github.com/Tom5521/gtk4tools/pkg/tools"
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

	c.Factory.ConnectSetup(tools.NewFactorySetup(func(listitem *gtk.ListItem) {
		if c.Setup == nil {
			return
		}
		c.Setup(listitem)
	}))
	c.Factory.ConnectBind(tools.NewFactoryBind(func(listitem *gtk.ListItem, pos int) {
		if c.Bind == nil {
			return
		}
		c.Bind(listitem, c.Model.At(pos))
	}))

	c.ColumnViewColumn = gtk.NewColumnViewColumn(title, &c.Factory.ListItemFactory)

	return c
}
