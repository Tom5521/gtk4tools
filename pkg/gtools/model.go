package gtools

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
)

var _ Modeller[any] = (*Model[any])(nil)

type Model[T any] struct {
	items []T
	// https://pkg.go.dev/github.com/diamondburned/gotk4/pkg/core/gioutil#ListModel
	listModel *gioutil.ListModel[T]
}

func NewModel[T any](items ...T) *Model[T] {
	return &Model[T]{items, NewListModel(items...)}
}

func (m *Model[T]) At(index int) T {
	return m.listModel.At(index)
}

func (m *Model[T]) S() []T {
	return m.items
}

func (m *Model[T]) Range(f func(i int, v T) bool) {
	for i, v := range m.items {
		if !f(i, v) {
			break
		}
	}
}

func (m *Model[T]) ListModel() *gioutil.ListModel[T] {
	return m.listModel
}

func (m *Model[T]) Len() int {
	return len(m.items)
}

// Re-generate the list with the items provided.
func (m *Model[T]) SetItems(items []T) {
	m.Splice(0, m.listModel.Len(), items...)
}

func (m *Model[T]) Remove(index int) {
	if index <= -1 {
		return
	}
	m.items = slices.Delete(m.items, index, index+1)
	m.listModel.Remove(index)
}

func (m *Model[T]) Append(items ...T) {
	m.items = append(m.items, items...)
	for _, i := range items {
		m.listModel.Append(i)
	}
}

func (m *Model[T]) Splice(pos, rms int, values ...T) {
	if pos <= -1 || rms <= -1 {
		return
	}
	m.listModel.Splice(pos, rms, values...)
	m.RefreshItems()
}

func (m *Model[T]) RefreshModel() {
	if m.listModel.Len() == 0 {
		for _, i := range m.items {
			m.listModel.Append(i)
		}
	}
	m.listModel.Splice(0, m.listModel.Len(), m.items...)
}

// Regenerates the List.Items based on the model.
func (m *Model[T]) RefreshItems() {
	m.items = []T{}
	m.listModel.All()(func(t T) bool {
		m.items = append(m.items, t)
		return true
	})
}
