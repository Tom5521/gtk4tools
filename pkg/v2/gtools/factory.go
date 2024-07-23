package gtools

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

var _ Factoryer[any] = (*Factory[any])(nil)

type Factory[T any] struct {
	factory *gtk.SignalListItemFactory

	model Modeller[T]

	setup FactorySetup
	bind  FactoryBind[T]
}

func NewFactory[T any](
	model Modeller[T],
	setup FactorySetup,
	bind FactoryBind[T],
) *Factory[T] {
	f := &Factory[T]{
		factory: gtk.NewSignalListItemFactory(),
		model:   model,
		setup:   setup,
		bind:    bind,
	}
	f.reconnectFactory()

	return f
}

func (f *Factory[T]) SetFactory(fac *gtk.SignalListItemFactory) {
	f.factory = fac
	f.reconnectFactory()
}

func (f *Factory[T]) Factory() *gtk.SignalListItemFactory {
	return f.factory
}

func (f *Factory[T]) Modeller() Modeller[T] {
	return f.model
}

func (f *Factory[T]) SetModeller(m Modeller[T]) {
	f.model = m
}

func (f *Factory[T]) Setup() FactorySetup {
	return f.setup
}

func (f *Factory[T]) SetSetup(s FactorySetup) {
	f.setup = s
}

func (f *Factory[T]) Bind() FactoryBind[T] {
	return f.bind
}

func (f *Factory[T]) SetBind(b FactoryBind[T]) {
	f.bind = b
}

func (f *Factory[T]) reconnectFactory() {
	f.factory.ConnectSetup(
		NewFactorySetup(func(listitem ListItem) {
			if f.setup != nil {
				f.setup(listitem)
			}
		}),
	)
	f.factory.ConnectBind(
		NewFactoryBind(func(listitem ListItem, pos int) {
			if f.bind != nil {
				f.bind(listitem, f.model.At(pos))
			}
		}),
	)
}
