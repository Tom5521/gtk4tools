package widgets

import (
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

var _ Selectioner = (*Selection)(nil)

type Selection struct {
	selectionMode     ListSelectionMode
	selectionModeller gtk.SelectionModeller
	listModeller      gio.ListModeller

	onSelected         func(index int)
	onMultipleSelected func(indexes []int)
}

func NewSelection(mode ListSelectionMode, model gio.ListModeller) *Selection {
	s := &Selection{
		listModeller: model,
	}
	s.SetSelectionMode(mode)

	return s
}

// If the assigned ListModeller is not compatible with this setting it will simply do nothing.
// This setting is only compatible with SelectionSingle.
func (s *Selection) Autoselect() bool {
	if m, ok := s.selectionModeller.(*gtk.SingleSelection); ok {
		return m.Autoselect()
	}
	return false
}

// If the assigned ListModeller is not compatible with this setting it will simply do nothing.
// This setting is only compatible with SelectionSingle.
func (s *Selection) SetAutoselect(b bool) {
	if m, ok := s.selectionModeller.(*gtk.SingleSelection); ok {
		m.SetAutoselect(b)
	}
}

func (s *Selection) CanUnselect() bool {
	if m, ok := s.selectionModeller.(*gtk.SingleSelection); ok {
		return m.CanUnselect()
	}

	return false
}

func (s *Selection) SetCanUnselect(b bool) {
	if m, ok := s.selectionModeller.(*gtk.SingleSelection); ok {
		m.SetAutoselect(b)
	}
}

func (s *Selection) ListModeller() gio.ListModeller {
	return s.listModeller
}

func (s *Selection) SetListModeller(m gio.ListModeller) {
	type setModeller interface {
		SetModel(gio.ListModeller)
	}

	s.listModeller = m
	s.selectionModeller.(setModeller).SetModel(m)
}

func (s *Selection) SetSelectionModeller(m gtk.SelectionModeller) {
	s.selectionModeller = m
	switch m.(type) {
	case *gtk.NoSelection:
		s.selectionMode = SelectionNone
	case *gtk.SingleSelection:
		s.selectionMode = SelectionSingle
	case *gtk.MultiSelection:
		s.selectionMode = SelectionMultiple
	default:
		s.selectionMode = SelectionNone
	}
	s.reConnectSelection()
}

func (s *Selection) SelectionModeller() gtk.SelectionModeller {
	return s.selectionModeller
}

func (s *Selection) ConnectSelected(f func(int)) {
	s.onSelected = f
}

func (s *Selection) ConnectMultipleSelected(f func([]int)) {
	s.onMultipleSelected = f
}

// Requests to select a range of items.
func (s *Selection) SelectRange(pos, nItems int, unSelectRest bool) {
	s.selectionModeller.SelectRange(uint(pos), uint(nItems), unSelectRest)
}

// With SelectionNone and SelectionSingle it does nothing, and with SelectionMultiple,
// it does what it promises, i.e., select all items in the list.
func (s *Selection) SelectAll() {
	s.selectionModeller.SelectAll()
}

// With SelectionSingle and SelectionMultiple it deselects the only element
// that can be selected, with SelectionNone it does nothing.
func (s *Selection) UnselectAll() {
	s.selectionModeller.UnselectAll()
}

// Requests to unselect a range of items.
func (s *Selection) UnselectRange(pos, nItems int) {
	s.selectionModeller.UnselectRange(uint(pos), uint(nItems))
}

// Requests to unselect an item.
func (s *Selection) Unselect(index int) {
	s.selectionModeller.UnselectItem(uint(index))
}

// Returns the index of the selected item,
// or -1 if its index is null or the selection model does not allow it.
func (s *Selection) Selected() int {
	model, ok := s.selectionModeller.(*gtk.SingleSelection)
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
func (s *Selection) Select(index int) {
	if index <= -1 {
		return
	}
	switch v := s.selectionModeller.(type) {
	case *gtk.MultiSelection:
		v.SelectItem(uint(index), true)
	case *gtk.SingleSelection:
		v.SetSelected(uint(index))
	}
}

// This method iterates over each element in the list and returns the selected ones.
// It only does something if the ListModel is a MultipleSelection,
// otherwise it simply returns an empty list.
func (s *Selection) MultipleSelected() []int {
	model, ok := s.selectionModeller.(*gtk.MultiSelection)
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
func (s *Selection) SelectMultiple(indexes ...int) {
	model, ok := s.selectionModeller.(*gtk.MultiSelection)
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

func (s *Selection) SetSelectionMode(mode ListSelectionMode) {
	s.selectionMode = mode
	switch mode {
	case SelectionNone:
		s.selectionModeller = gtk.NewNoSelection(s.listModeller)
	case SelectionSingle:
		s.selectionModeller = gtk.NewSingleSelection(s.listModeller)
	case SelectionMultiple:
		s.selectionModeller = gtk.NewMultiSelection(s.listModeller)
	default:
		s.selectionModeller = gtk.NewNoSelection(s.listModeller)
	}
	s.reConnectSelection()
}

func (s *Selection) SelectionMode() ListSelectionMode {
	return s.selectionMode
}

func (m *Selection) reConnectSelection() {
	m.selectionModeller.ConnectSelectionChanged(func(_, _ uint) {
		switch m.selectionModeller.(type) {
		case *gtk.SingleSelection:
			if m.onSelected != nil {
				m.onSelected(m.Selected())
			}
		case *gtk.MultiSelection:
			if m.onMultipleSelected != nil {
				m.onMultipleSelected(m.MultipleSelected())
			}
		}
	})
}
