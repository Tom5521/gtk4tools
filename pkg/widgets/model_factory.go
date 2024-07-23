package widgets

import (
	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ModelFactory[T any, M SetModelFactoryer] struct {
	Modeller[T]
	Selectioner

	Setup FactorySetup
	Bind  FactoryBind[T]

	ItemFactory *gtk.SignalListItemFactory

	Setter M
}

func NewModelFactory[T any, M SetModelFactoryer](
	smodel ListSelectionMode,
	sm M,
	setup FactorySetup,
	bind FactoryBind[T],
	modeller Modeller[T],
) *ModelFactory[T, M] {
	m := &ModelFactory[T, M]{
		Setup:    setup,
		Bind:     bind,
		Modeller: modeller,
		Setter:   sm,

		// Initialization.
		ItemFactory: gtk.NewSignalListItemFactory(),
	}

	m.Selectioner = NewSelection(smodel, m.Modeller.ListModel())
	m.reConnectFactory()

	m.InitSetter()
	return m
}

func (m *ModelFactory[T, M]) InitSetter() {
	m.Setter.SetModel(m.SelectionModeller())
	m.Setter.SetFactory(&m.ItemFactory.ListItemFactory)
}

func (m *ModelFactory[T, _]) SetSelectionMode(mode ListSelectionMode) {
	m.Selectioner.SetSelectionMode(mode)
	m.Setter.SetModel(m.SelectionModeller())
}

// Refreshes absolutely everything. To be more specific here is the list of what it refreshes:
//
// - Refreshes List.Items
//
// - Refreshes the factory
//
// - Refreshes the model
//
// I generally discourage its use and prefer to refresh things as
// they are modified manually and individually.
func (m *ModelFactory[T, _]) Refresh() {
	m.RefreshItems()
	m.reConnectFactory()
	m.RefreshModel()
}

// Can be used when modifying List.Setup and/or List.Bind to redraw
// the entire list following the new Setup and Bind.
func (m *ModelFactory[T, _]) RefreshFactory() {
	m.reConnectFactory()
	m.RefreshModel()
}

// Internal methods.

func (m *ModelFactory[T, _]) reConnectFactory() {
	m.ItemFactory.ConnectSetup(gtools.NewFactorySetup(func(listitem gtools.ListItem) {
		if m.Setup == nil {
			return
		}
		m.Setup(listitem)
	}))
	m.ItemFactory.ConnectBind(gtools.NewFactoryBind(func(listitem gtools.ListItem, pos int) {
		if m.Bind == nil {
			return
		}
		m.Bind(listitem, m.At(pos))
	}))
}
