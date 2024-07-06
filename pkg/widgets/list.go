package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ListBind[T any] func(*gtk.ListItem, T)

type List[T any] struct {
	*gtk.ListView

	Items []T

	Setup FactorySetup
	Bind  ListBind[T]

	OnSelected         func(index int)
	OnMultipleSelected func(indexes []int)

	Factory *gtk.SignalListItemFactory

	SelectionMode     ListSelectionMode
	SelectionModeller gtk.SelectionModeller
	// https://pkg.go.dev/github.com/diamondburned/gotk4/pkg/core/gioutil#ListModel
	Model *gioutil.ListModel[T]
}

func NewList[T any](
	items []T,
	smodel ListSelectionMode,
	setup FactorySetup,
	bind ListBind[T],
) *List[T] {
	l := &List[T]{
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

func (l *List[T]) SetSelectionModeller(mode ListSelectionMode) {
	l.SelectionMode = mode
	l.makeSelectionModeller(mode)
	l.ListView.SetModel(l.SelectionModeller)
	l.reConnectSelection()
}

// Re-generate the list with the items provided.
func (l *List[T]) SetItems(items []T) {
	l.Splice(0, l.Model.NItems(), items...)
}

func (l *List[T]) Remove(index int) {
	if index <= -1 {
		return
	}
	l.Items = slices.Delete(l.Items, index, index+1)
	l.Model.Remove(index)
}

func (l *List[T]) Append(item T) {
	l.Items = append(l.Items, item)
	l.Model.Append(item)
}

func (l *List[T]) Splice(pos, rms int, values ...T) {
	if pos <= -1 || rms <= -1 {
		return
	}
	spliceVar(&l.Items, pos, rms, values)
	l.Model.Splice(pos, rms, values...)
}

// Returns the index of the selected item,
// or -1 if its index is null or the selection model does not allow it.
func (l *List[T]) Selected() int {
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
func (l *List[T]) MultipleSelected() []int {
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

// The Select method behaves as expected when applied to a SingleSelection SelectionModel.
// It selects the item at the specified index. However, when applied to a MultipleSelection,
// it behaves differently. Instead of selecting only the item at the specified index,
// it deselects all other items and selects only the one at that index.
//
// In the case of NoSelection, it simply does nothing.
func (l *List[T]) Select(index int) {
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
func (l *List[T]) SelectMultiple(indexes ...int) {
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

// Requests to select a range of items.
func (l *List[T]) SelectRange(pos, nItems int, unSelectRest bool) {
	l.SelectionModeller.SelectRange(uint(pos), uint(nItems), unSelectRest)
}

// With SelectionNone and SelectionSingle it does nothing, and with SelectionMultiple,
// it does what it promises, i.e., select all items in the list.
func (l *List[T]) SelectAll() {
	l.SelectionModeller.SelectAll()
}

// With SelectionSingle and SelectionMultiple it deselects the only element
// that can be selected, with SelectionNone it does nothing.
func (l *List[T]) UnselectAll() {
	l.SelectionModeller.UnselectAll()
}

// Requests to unselect a range of items.
func (l *List[T]) UnselectRange(pos, nItems int) {
	l.SelectionModeller.UnselectRange(uint(pos), uint(nItems))
}

// Requests to unselect an item.
func (l *List[T]) Unselect(index int) {
	l.SelectionModeller.UnselectItem(uint(index))
}

// Regenerates the List.Items based on the model.
func (l *List[T]) RefreshItems() {
	l.Items = []T{}
	for i := range l.Model.NItems() {
		l.Items = append(l.Items, l.Model.Item(i))
	}
}

// Can be used when modifying List.Setup and/or List.Bind to redraw
// the entire list following the new Setup and Bind.
func (l *List[T]) RefreshFactory() {
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
func (l *List[T]) Refresh() {
	l.RefreshItems()
	l.RefreshSelectionModeller()
	l.reConnectFactory()
	l.RefreshModel()
}

func (l *List[T]) RefreshModel() {
	if l.Model.NItems() == 0 {
		for _, i := range l.Items {
			l.Model.Append(i)
		}
	}
	l.Model.Splice(0, l.Model.NItems(), l.Items...)
}

func (l *List[T]) RefreshSelectionModeller() {
	l.SetSelectionModeller(l.SelectionMode)
}

// Internal functions

func (l *List[T]) reConnectFactory() {
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

func (l *List[T]) reConnectSelection() {
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

func (l *List[T]) makeSelectionModeller(mode ListSelectionMode) {
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
