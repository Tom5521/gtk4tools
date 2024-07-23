package widget

import (
	"github.com/Tom5521/gtk4tools/pkg/v2/gtools"
)

type TemplateView[T any, M gtools.SetModelFactoryer] struct {
	gtools.Modeller[T]
	gtools.Factoryer[T]
	gtools.Selectioner

	Setter M
}

func NewTemplateView[T any, M gtools.SetModelFactoryer](
	smodel gtools.ListSelectionMode,
	sm M,
	setup gtools.FactorySetup,
	bind gtools.FactoryBind[T],
	modeller gtools.Modeller[T],
) *TemplateView[T, M] {
	m := &TemplateView[T, M]{
		Modeller: modeller,
		Setter:   sm,
	}

	m.Selectioner = gtools.NewSelection(
		smodel,
		m.Modeller.ListModel(),
	)
	m.Factoryer = gtools.NewFactory(
		m.Modeller,
		setup,
		bind,
	)

	m.InitSetter()
	return m
}

func (m *TemplateView[T, M]) InitSetter() {
	m.Setter.SetModel(m.SelectionModeller())
	m.Setter.SetFactory(&m.Factoryer.Factory().ListItemFactory)
}

func (m *TemplateView[T, _]) SetSelectionMode(mode gtools.ListSelectionMode) {
	m.Selectioner.SetSelectionMode(mode)
	m.Setter.SetModel(m.SelectionModeller())
}

// Refreshes absolutely everything. To be more specific here is the list of what it refreshes:
//
// - Refreshes List.Items
//
// - Refreshes the model
//
// I generally discourage its use and prefer to refresh things as
// they are modified manually and individually.
func (m *TemplateView[T, _]) Refresh() {
	m.RefreshItems()
	m.RefreshModel()
}
