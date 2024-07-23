package widgets

import (
	"slices"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
)

var _ Modeller[any] = (*ModelVar[any])(nil)

type ModelVar[T any] struct {
	items     *[]T
	listModel *gioutil.ListModel[T]
}

func NewModelVar[T any](items *[]T) *ModelVar[T] {
	m := &ModelVar[T]{
		items:     items,
		listModel: NewListModel[T](*items...),
	}

	return m
}

func (m *ModelVar[T]) ListModel() *gioutil.ListModel[T] {
	return m.listModel
}

func (m *ModelVar[T]) S() []T {
	return *m.items
}

func (m *ModelVar[T]) Ptr() *[]T {
	return m.items
}

func (m *ModelVar[T]) Len() int {
	return len(*m.items)
}

func (m *ModelVar[T]) SetItems(items []T) {
	*m.items = items
}

func (m *ModelVar[T]) Remove(i int) {
	*m.items = slices.Delete(*m.items, i, i+1)
	m.listModel.Remove(i)
}

func (m *ModelVar[T]) Append(items ...T) {
	*m.items = append(*m.items, items...)
	for _, i := range items {
		m.listModel.Append(i)
	}
}

func (m *ModelVar[T]) Splice(pos, rms int, appends ...T) {
	m.listModel.Splice(pos, rms, appends...)
	m.RefreshItems()
}

func (m *ModelVar[T]) RefreshModel() {
	m.listModel.Splice(0, m.listModel.Len(), *m.items...)
}

func (m *ModelVar[T]) RefreshItems() {
	*m.items = []T{}
	m.listModel.All()(func(t T) bool {
		*m.items = append(*m.items, t)
		return true
	})
}

func (m *ModelVar[T]) At(i int) T {
	return m.listModel.At(i)
}

func (m *ModelVar[T]) Range(f func(int, T) bool) {
	for i, v := range *m.items {
		if !f(i, v) {
			break
		}
	}
}
