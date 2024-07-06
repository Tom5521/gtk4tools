package widgets

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type ColumnView[T any] struct {
	*gtk.ColumnView

	Columns []*Column[T]

	SelectionModeller gtk.SelectionModeller
}

func NewColumnView[T any](
	smode gtk.SelectionModeller,
	columns ...*Column[T],
) *ColumnView[T] {
	c := &ColumnView[T]{
		ColumnView: gtk.NewColumnView(smode),
	}
	c.SetColumns(columns...)

	return c
}

func (c *ColumnView[T]) SetColumns(columns ...*Column[T]) {
	cols := c.ColumnView.Columns()
	for i := range cols.NItems() {
		col := cols.Item(i).Cast().(*gtk.ColumnViewColumn)
		c.ColumnView.RemoveColumn(col)
	}
	c.Columns = columns
	for _, col := range c.Columns {
		c.ColumnView.AppendColumn(col.ColumnViewColumn)
	}
}
