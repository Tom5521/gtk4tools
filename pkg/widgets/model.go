package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
)

type Model[T any] struct {
	Items     []T
	ListModel *gioutil.ListModel[T]
}

func NewModel[T any](items ...T) *Model[T] {
	return &Model[T]{items, NewListModel(items...)}
}

// Re-generate the list with the items provided.
func (m *Model[T]) SetItems(items []T) {
	m.Splice(0, m.ListModel.Len(), items...)
}

func (m *Model[T]) Remove(index int) {
	if index <= -1 {
		return
	}
	m.Items = slices.Delete(m.Items, index, index+1)
	m.ListModel.Remove(index)
}

func (m *Model[T]) Append(items ...T) {
	m.Items = append(m.Items, items...)
	for _, i := range items {
		m.ListModel.Append(i)
	}
}

func (m *Model[T]) Splice(pos, rms int, values ...T) {
	if pos <= -1 || rms <= -1 {
		return
	}
	m.ListModel.Splice(pos, rms, values...)
	m.RefreshItems()
}

func (m *Model[T]) RefreshModel() {
	if m.ListModel.Len() == 0 {
		for _, i := range m.Items {
			m.ListModel.Append(i)
		}
	}
	m.ListModel.Splice(0, m.ListModel.Len(), m.Items...)
}

// Regenerates the List.Items based on the model.
func (m *Model[T]) RefreshItems() {
	m.Items = []T{}
	m.ListModel.All()(func(t T) bool {
		m.Items = append(m.Items, t)
		return true
	})
}
