package gtools

import "github.com/diamondburned/gotk4/pkg/core/gioutil"

func NewListModel[T any](items ...T) *gioutil.ListModel[T] {
	m := gioutil.NewListModel[T]()
	for _, i := range items {
		m.Append(i)
	}

	return m
}
