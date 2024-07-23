package widgets

import (
	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ModelFactory[T any, M Setter] struct {
	*Model[T]

	Setup FactorySetup
	Bind  FactoryBind[T]

	OnSelected         func(index int)
	OnMultipleSelected func(indexes []int)

	ItemFactory *gtk.SignalListItemFactory

	SelectionMode     ListSelectionMode
	SelectionModeller gtk.SelectionModeller

	// https://pkg.go.dev/github.com/diamondburned/gotk4/pkg/core/gioutil#ListModel

	Setter M
}

func NewModelFactory[T any, M Setter](
	smodel ListSelectionMode,
	sm M,
	setup FactorySetup,
	bind FactoryBind[T],
	items ...T,
) *ModelFactory[T, M] {
	m := &ModelFactory[T, M]{
		Setup:         setup,
		Bind:          bind,
		SelectionMode: smodel,

		Setter: sm,

		// Initialization.
		ItemFactory: gtk.NewSignalListItemFactory(),
		Model:       NewModel(items...),
	}

	m.reConnectFactory()
	m.makeSelectionModeller(smodel)
	m.reConnectSelection()

	m.InitSetter()
	return m
}

func (m *ModelFactory[T, M]) InitSetter() {
	m.Setter.SetModel(m.SelectionModeller)
	m.Setter.SetFactory(&m.ItemFactory.ListItemFactory)
}

func (m *ModelFactory[T, _]) SetSelectionModeller(mode ListSelectionMode) {
	m.SelectionMode = mode
	m.makeSelectionModeller(mode)
	m.Setter.SetModel(m.SelectionModeller)
	m.reConnectSelection()
}

// Requests to select a range of items.
func (m ModelFactory[T, _]) SelectRange(pos, nItems int, unSelectRest bool) {
	m.SelectionModeller.SelectRange(uint(pos), uint(nItems), unSelectRest)
}

// With SelectionNone and SelectionSingle it does nothing, and with SelectionMultiple,
// it does what it promises, i.e., select all items in the list.
func (m *ModelFactory[T, _]) SelectAll() {
	m.SelectionModeller.SelectAll()
}

// With SelectionSingle and SelectionMultiple it deselects the only element
// that can be selected, with SelectionNone it does nothing.
func (m *ModelFactory[T, _]) UnselectAll() {
	m.SelectionModeller.UnselectAll()
}

// Requests to unselect a range of items.
func (m *ModelFactory[T, _]) UnselectRange(pos, nItems int) {
	m.SelectionModeller.UnselectRange(uint(pos), uint(nItems))
}

// Requests to unselect an item.
func (m *ModelFactory[T, _]) Unselect(index int) {
	m.SelectionModeller.UnselectItem(uint(index))
}

// Returns the index of the selected item,
// or -1 if its index is null or the selection model does not allow it.
func (m *ModelFactory[T, _]) Selected() int {
	model, ok := m.SelectionModeller.(*gtk.SingleSelection)
	if !ok {
		return -1
	}
	i := model.Selected()
	if model.Item(i) == nil {
		return -1
	}

	return int(i)
}

// The Select method behaves as expected when applied to a SingleSelection SelectionModel.
// It selects the item at the specified index. However, when applied to a MultipleSelection,
// it behaves differently. Instead of selecting only the item at the specified index,
// it deselects all other items and selects only the one at that index.
//
// In the case of NoSelection, it simply does nothing.
func (l *ModelFactory[T, _]) Select(index int) {
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

// This method iterates over each element in the list and returns the selected ones.
// It only does something if the ListModel is a MultipleSelection,
// otherwise it simply returns an empty list.
func (m *ModelFactory[T, _]) MultipleSelected() []int {
	model, ok := m.SelectionModeller.(*gtk.MultiSelection)
	if !ok {
		return nil
	}
	var out []int
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

// This method only does something when the ListModel is of the MultiSelection type,
// and it requires you to pass the indexes that it will select.
// If any of the indexes cannot be converted to uint, it will simply iterate to the next element.
// However, if any of the indexes are not in the list,
// it will throw an error, i.e. it will crash.
//
// If an element is already selected, deselects it.
func (l *ModelFactory[T, _]) SelectMultiple(indexes ...int) {
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

func (m *ModelFactory[T, _]) RefreshSelectionModeller() {
	m.SetSelectionModeller(m.SelectionMode)
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
func (m *ModelFactory[T, _]) Refresh() {
	m.RefreshItems()
	m.RefreshSelectionModeller()
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

func (m *ModelFactory[T, _]) reConnectSelection() {
	m.SelectionModeller.ConnectSelectionChanged(func(_, _ uint) {
		switch m.SelectionModeller.(type) {
		case *gtk.SingleSelection:
			if m.OnSelected != nil {
				m.OnSelected(m.Selected())
			}
		case *gtk.MultiSelection:
			if m.OnMultipleSelected != nil {
				m.OnMultipleSelected(m.MultipleSelected())
			}
		}
	})
}

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
		m.Bind(listitem, m.ListModel.At(pos))
	}))
}

func (m *ModelFactory[T, _]) makeSelectionModeller(mode ListSelectionMode) {
	m.SelectionMode = mode
	switch mode {
	case SelectionNone:
		m.SelectionModeller = gtk.NewNoSelection(m.ListModel.ListModel)
	case SelectionSingle:
		m.SelectionModeller = gtk.NewSingleSelection(m.ListModel.ListModel)
	case SelectionMultiple:
		m.SelectionModeller = gtk.NewMultiSelection(m.ListModel.ListModel)
	default:
		m.SelectionModeller = gtk.NewNoSelection(m.ListModel.ListModel)
	}
}
