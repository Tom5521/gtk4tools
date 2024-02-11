package tools

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type SetParenter interface {
	SetParent(gtk.Widgetter)
}

type SetChilder interface {
	SetChild(gtk.Widgetter)
}

func SetParent(parent gtk.Widgetter, childs ...SetParenter) {
	for _, ch := range childs {
		ch.SetParent(parent)
	}
}

func SetChilds(parent SetChilder, childs ...gtk.Widgetter) {
	for _, ch := range childs {
		parent.SetChild(ch)
	}
}
