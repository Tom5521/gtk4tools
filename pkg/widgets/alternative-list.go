package widgets

import (
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// REFERENCE:
// https://pkg.go.dev/github.com/diamondburned/gotk4/pkg/core/gioutil#ListModel

type AlternativeListSetup func(listitem *gtk.ListItem)
type AlternativeListBind func(listitem *gtk.ListItem, index int)
type AlternativeListLen func() int

type AlternativeList struct {
	*gtk.ListView

	Setup AlternativeListSetup
	Bind  AlternativeListBind
	Len   AlternativeListLen

	OnSelected         func(index int)
	OnMultipleSelected func(indexes []int)

	Model *gioutil.ListModel[int]

	SelectionModeller gtk.SelectionModeller
	SelectionMode     ListSelectionMode
	Factory           *gtk.SignalListItemFactory
}

func NewAlternativeList(
	smodel ListSelectionMode,
	lenfunc AlternativeListLen,
	setup AlternativeListSetup,
	bind AlternativeListBind,
) *AlternativeList {
	l := &AlternativeList{
		Model:         gioutil.NewListModelType[int]().New(),
		Factory:       gtk.NewSignalListItemFactory(),
		SelectionMode: smodel,

		Setup: setup,
		Bind:  bind,
		Len:   lenfunc,
	}

	l.reConnectFactory()
	l.RefreshModel()
	l.makeSelectionModeller(smodel)

	l.ListView = gtk.NewListView(l.SelectionModeller, &l.Factory.ListItemFactory)

	return l
}

// PUBLIC METHODS

// Returns the index of the selected item,
// or -1 if its index is null or the selection model does not allow it.
func (l *AlternativeList) Selected() int {
	model, ok := l.SelectionModeller.(*gtk.SingleSelection)
	if !ok {
		return -1
	}
	i := model.Selected()
	if model.Item(i) == nil {
		return -1
	}

	return int(i)
}

// This method iterates over each element in the list and returns the selected ones.
// It only does something if the ListModel is a MultipleSelection,
// otherwise it simply returns an empty list.
func (l *AlternativeList) MultipleSelected() []int {
	var out []int

	model, ok := l.SelectionModeller.(*gtk.MultiSelection)
	if !ok {
		return out
	}
	for i := range model.NItems() {
		item := model.Item(i)
		if item == nil {
			continue
		}
		if model.IsSelected(i) {
			out = append(out, int(i))
		}
	}

	return out
}

// The SetSelected method behaves as expected when applied to a SingleSelection SelectionModel.
// It selects the item at the specified index. However, when applied to a MultipleSelection,
// it behaves differently. Instead of selecting only the item at the specified index,
// it deselects all other items and selects only the one at that index.
//
// In the case of NoSelection, it simply does nothing.
func (l *AlternativeList) SetSelected(index int) {
	if index <= -1 {
		return
	}
	switch v := l.SelectionModeller.(type) {
	case *gtk.MultiSelection:
		v.SelectItem(uint(index), true)
	case *gtk.SingleSelection:
		v.SetSelected(uint(index))
	}
}

// This method only does something when the ListModel is of the MultiSelection type,
// and it requires you to pass the indexes that it will select.
// If any of the indexes cannot be converted to uint, it will simply iterate to the next element.
// However, if any of the indexes are not in the list,
// it will throw an error, i.e. it will crash.
//
// If an element is already selected, deselects it.
func (l *AlternativeList) SetMultipleSelections(indexes ...int) {
	model, ok := l.SelectionModeller.(*gtk.MultiSelection)
	if !ok {
		return
	}
	for _, i := range indexes {
		if i <= -1 {
			continue
		}
		if model.IsSelected(uint(i)) {
			model.UnselectItem(uint(i))
			continue
		}
		model.SelectItem(uint(i), false)
	}
}

func (l *AlternativeList) SetSelectionModeller(mode ListSelectionMode) {
	l.SelectionMode = mode
	l.makeSelectionModeller(mode)
	l.ListView.SetModel(l.SelectionModeller)
	l.ReConnectSelection()
}

func (l *AlternativeList) RefreshSelectionModeller() {
	l.SetSelectionModeller(l.SelectionMode)
}

// Can be used when modifying List.Setup and/or List.Bind to redraw
// the entire list following the new Setup and Bind.
func (l *AlternativeList) RefreshFactory() {
	l.reConnectFactory()
	l.RefreshModel()
}

// Refreshes absolutely everything. To be more specific here is the list of what it refreshes:
//
// - Refreshes List.Items
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

func (l *AlternativeList) ReConnectSelection() {
	l.SelectionModeller.ConnectSelectionChanged(func(_, _ uint) {
		switch l.SelectionModeller.(type) {
		case *gtk.SingleSelection:
			if l.OnSelected != nil {
				l.OnSelected(l.Selected())
			}
		case *gtk.MultiSelection:
			if l.OnMultipleSelected != nil {
				l.OnMultipleSelected(l.MultipleSelected())
			}
		}
	})
}

// PRIVATE METHODS

func (l *AlternativeList) makeSelectionModeller(mode ListSelectionMode) {
	l.SelectionMode = mode
	switch mode {
	case SelectionNone:
		l.SelectionModeller = gtk.NewNoSelection(l.Model.ListModel)
	case SelectionSingle:
		l.SelectionModeller = gtk.NewSingleSelection(l.Model.ListModel)
	case SelectionMultiple:
		l.SelectionModeller = gtk.NewMultiSelection(l.Model.ListModel)
	default:
		l.SelectionModeller = gtk.NewNoSelection(l.Model.ListModel)
	}
}

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
