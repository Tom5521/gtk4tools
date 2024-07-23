package widgets

import (
	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type (
	AlternativeListBind FactoryBind[int]
	AlternativeListLen  func() int
)

type AlternativeList struct {
	*gtk.ListView

	*ModelFactory[int, *gtk.ListView]

	Len AlternativeListLen
}

// It creates an alternative list, which is practically the same as fyne's,
// that only requires the slice length to work. The only real difference is
// that the list has to be refreshed manually with each modification of the base slice.
func NewAlternativeList(
	smodel ListSelectionMode,
	lenfunc AlternativeListLen,
	setup FactorySetup,
	bind AlternativeListBind,
) *AlternativeList {
	l := &AlternativeList{
		Len: lenfunc,
		ModelFactory: &ModelFactory[int, *gtk.ListView]{
			Setup: setup,
			Bind: func(li gtools.ListItem, t int) {
				bind(li, t)
			},
			SelectionMode: smodel,
			Setter:        gtk.NewListView(nil, nil),

			ItemFactory: gtk.NewSignalListItemFactory(),
			Model:       NewModel[int](),
		},
	}

	l.reConnectFactory()
	l.RefreshModel()
	l.makeSelectionModeller(smodel)
	l.reConnectSelection()

	l.InitSetter()

	l.ListView = l.Setter

	return l
}

// PUBLIC METHODS

// Refreshes absolutely everything. To be more specific here is the list of what it refreshes:
//
// - Refreshes the selection mode
//
// - Refreshes the factory
//
// - Refreshes the model
//
// I generally discourage its use and prefer to refresh things as
// they are modified manually and individually.
func (l *AlternativeList) Refresh() {
	l.RefreshSelectionModeller()
	l.reConnectFactory()
	l.RefreshModel()
}

func (l *AlternativeList) RefreshModel() {
	l.ListModel.Splice(0, l.ListModel.Len(), make([]int, l.Len())...)
}

// PRIVATE METHODS

func (l *AlternativeList) reConnectFactory() {
	l.ItemFactory.ConnectSetup(gtools.NewFactorySetup(func(listitem gtools.ListItem) {
		if l.Setup == nil {
			return
		}
		l.Setup(listitem)
	}))
	l.ItemFactory.ConnectBind(gtools.NewFactoryBind(func(listitem gtools.ListItem, pos int) {
		l.Bind(listitem, pos)
	}))
}
