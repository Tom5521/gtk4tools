package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type CustomListBind[T any] func(*gtk.ListItem, T)

type CustomList[T any] struct {
	*gtk.ListView

	Items []T

	Setup ListSetup
	Bind  CustomListBind[T]

	OnSelected         func(index int)
	OnMultipleSelected func(indexes []int)

	Factory *gtk.SignalListItemFactory

	SelectionMode     ListSelectionMode
	SelectionModeller gtk.SelectionModeller
	Model             *gioutil.ListModel[T]
}

func NewCustomList[T any](
	items []T,
	smodel ListSelectionMode,
	setup ListSetup,
	bind CustomListBind[T],
) *CustomList[T] {
	l := &CustomList[T]{
		Items:         items,
		Setup:         setup,
		Bind:          bind,
		SelectionMode: smodel,

		// Initialization
		Factory: gtk.NewSignalListItemFactory(),
		Model:   gioutil.NewListModel[T](),
	}

	l.RefreshModel()
	l.reConnectFactory()
	l.makeSelectionModeller(smodel)
	l.reConnectSelection()

	l.ListView = gtk.NewListView(l.SelectionModeller, &l.Factory.ListItemFactory)

	return l
}
func (l *CustomList[T]) SetSelectionModeller(mode ListSelectionMode) {
	l.SelectionMode = mode
	l.makeSelectionModeller(mode)
	l.ListView.SetModel(l.SelectionModeller)
	l.reConnectSelection()
}

func (l *CustomList[T]) RefreshSelectionModeller() {
	l.SetSelectionModeller(l.SelectionMode)
}

// Re-generate the list with the items provided.
func (l *CustomList[T]) SetItems(items []T) {
	l.Splice(0, l.Model.NItems(), items...)
}

func (l *CustomList[T]) Remove(index int) {
	if index <= -1 {
		return
	}
	l.Items = slices.Delete(l.Items, index, index+1)
	l.Model.Remove(index)
}

func (l *CustomList[T]) Append(item T) {
	l.Items = append(l.Items, item)
	l.Model.Append(item)
}

func (l *CustomList[T]) Splice(pos, nRemovals int, additions ...T) {
	if pos <= -1 || nRemovals <= -1 {
		return
	}
	l.Items = slices.Delete(l.Items, pos, nRemovals)
	l.Items = append(l.Items, additions...)
	l.Model.Splice(pos, nRemovals, additions...)
}

// Returns the index of the selected item,
// or -1 if its index is null or the selection model does not allow it.
func (l *CustomList[T]) Selected() int {
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
func (l *CustomList[T]) MultipleSelected() []int {
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
func (l *CustomList[T]) SetSelected(index int) {
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
func (l *CustomList[T]) SetMultipleSelections(indexes ...int) {
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

func (l *CustomList[T]) RefreshItems() {
	l.Items = []T{}
	for i := range l.Model.NItems() {
		l.Items = append(l.Items, l.Model.Item(i))
	}
}

// Can be used when modifying List.Setup and/or List.Bind to redraw
// the entire list following the new Setup and Bind.
func (l *CustomList[T]) RefreshFactory() {
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
func (l *CustomList[T]) Refresh() {
	l.RefreshItems()
	l.RefreshSelectionModeller()
	l.reConnectFactory()
	l.RefreshModel()
}

// Internal functions

func (l *CustomList[T]) reConnectFactory() {
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
		l.Bind(listitem, l.Items[listitem.Position()])
	})
}

func (l *CustomList[T]) reConnectSelection() {
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

func (l *CustomList[T]) makeSelectionModeller(mode ListSelectionMode) {
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

func (l *CustomList[T]) RefreshModel() {
	if l.Model.NItems() == 0 {
		for _, i := range l.Items {
			l.Model.Append(i)
		}
	}
	l.Model.Splice(0, l.Model.NItems(), l.Items...)
}
