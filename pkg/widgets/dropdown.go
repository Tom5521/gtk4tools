//go:build dev

package widgets

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type DropDown[T any] struct {
	*gtk.DropDown

	Model[T]
}

func NewDropDown[T any](
	items []T,
	setup FactorySetup,
	bind FactoryBind[T],
) *DropDown[T] {
	d := &DropDown[T]{
		// Model: ,
	}

	// gtk.NewDropDown(model gio.ListModeller, expression gtk.Expressioner)

	return d
}
