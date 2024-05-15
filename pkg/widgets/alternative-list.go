package widgets

import (
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type AlternativeListSetup func(listitem *gtk.ListItem)
type AlternativeListBind func(listitem *gtk.ListItem, index int)
type AlternativeListLen func() int

type AlternativeList struct {
	*List[int]

	Len AlternativeListLen
}

// It creates an alternative list, which is practically the same as fyne's,
// that only requires the slice length to work. The only real difference is
// that the list has to be refreshed manually with each modification of the base slice.
func NewAlternativeList(
	smodel ListSelectionMode,
	lenfunc AlternativeListLen,
	setup AlternativeListSetup,
	bind AlternativeListBind,
) *AlternativeList {
	l := &AlternativeList{
		List: &List[int]{
			Model:         gioutil.NewListModel[int](),
			Factory:       gtk.NewSignalListItemFactory(),
			SelectionMode: smodel,
			Setup:         ListSetup(setup),
			Bind:          ListBind[int](bind),
		},
		Len: lenfunc,
	}

	l.reConnectFactory()
	l.RefreshModel()
	l.makeSelectionModeller(smodel)
	l.reConnectSelection()

	l.ListView = gtk.NewListView(l.SelectionModeller, &l.Factory.ListItemFactory)

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
	l.Model.Splice(0, l.Model.NItems(), make([]int, l.Len())...)
}

// PRIVATE METHODS

func (l *AlternativeList) reConnectFactory() {
	l.Factory.ConnectSetup(func(listitem *gtk.ListItem) {
		if l.Setup == nil {
			return
		}
		l.Setup(listitem)
	})
	l.Factory.ConnectBind(func(listitem *gtk.ListItem) {
		if l.Bind == nil {
			return
		}
		l.Bind(listitem, int(listitem.Position()))
	})
}
