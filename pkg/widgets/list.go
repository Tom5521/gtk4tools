package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ListSelectionMode int

const (
	SelectionNone ListSelectionMode = iota
	SelectionSingle
	SelectionMultiple
)

type List struct {
	*gtk.ListView

	Items []string

	OnSelected         func(index int)
	OnMultipleSelected func(indexes []int)

	SelectionModeller gtk.SelectionModeller
	SelectionMode     ListSelectionMode
	Model             *gtk.StringList
	Factory           *gtk.SignalListItemFactory
}

// Creates a new list that keeps the self.Items[] updated with that of the UI.
func NewList(items []string, smodel ListSelectionMode, setup, bind func(listitem *gtk.ListItem)) *List {
	l := &List{
		Items:         items,
		SelectionMode: smodel,
		Model:         gtk.NewStringList(items),
		Factory:       gtk.NewSignalListItemFactory(),
	}
	l.Factory.ConnectSetup(setup)
	l.Factory.ConnectBind(bind)

	l.makeSelectionModeller(smodel)
	l.reConnectSelection()

	l.ListView = gtk.NewListView(l.SelectionModeller, &l.Factory.ListItemFactory)

	return l
}

func (l *List) SetSelectionModeller(mode ListSelectionMode) {
	l.SelectionMode = mode
	l.makeSelectionModeller(mode)
	l.ListView.SetModel(l.SelectionModeller)
	l.reConnectSelection()
}

func (l *List) RefreshSelectionModeller() {
	l.SetSelectionModeller(l.SelectionMode)
}

// Re-generate the list with the items provided.
func (l *List) SetItems(items []string) {
	l.Splice(0, int(l.Model.NItems()), items...)
}

func (l *List) Remove(index int) {
	if index <= -1 {
		return
	}
	l.Items = slices.Delete(l.Items, index, index+1)
	l.Model.Remove(uint(index))
}

func (l *List) Append(item string) {
	l.Items = append(l.Items, item)
	l.Model.Append(item)
}

func (l *List) Splice(pos, nRemovals int, additions ...string) {
	if pos <= -1 || nRemovals <= -1 {
		return
	}
	l.Items = slices.Delete(l.Items, pos, nRemovals)
	l.Items = append(l.Items, additions...)
	l.Model.Splice(uint(pos), uint(nRemovals), additions)
}

// Returns the index of the selected item,
// or -1 if its index is null or the selection model does not allow it.
func (l *List) Selected() int {
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
func (l *List) MultipleSelected() []int {
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
func (l *List) SetSelected(index int) {
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
func (l *List) SetMultipleSelections(indexes ...int) {
	model, ok := l.SelectionModeller.(*gtk.MultiSelection)
	if !ok {
		return
	}
	for _, i := range indexes {
		if i <= -1 {
			continue
		}
		if model.IsSelected(uint(i)) {
			continue
		}
		model.SelectItem(uint(i), false)
	}
}

func (l *List) RefreshItems() {
	l.Items = []string{}
	for i := range l.Model.NItems() {
		item := l.Model.Item(i)
		if item == nil {
			continue
		}
		l.Items = append(l.Items, item.Cast().(*gtk.StringObject).String())
	}
}

// Internal functions

func (l *List) reConnectSelection() {
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

func (l *List) makeSelectionModeller(mode ListSelectionMode) {
	l.SelectionMode = mode
	switch mode {
	case SelectionNone:
		l.SelectionModeller = gtk.NewNoSelection(l.Model)
	case SelectionSingle:
		l.SelectionModeller = gtk.NewSingleSelection(l.Model)
	case SelectionMultiple:
		l.SelectionModeller = gtk.NewMultiSelection(l.Model)
	default:
		l.SelectionModeller = gtk.NewNoSelection(l.Model)
	}
}
