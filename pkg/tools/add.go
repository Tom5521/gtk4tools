package tools

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type Appender interface {
	Append(gtk.Widgetter)
}

type AddChilder interface {
	AddChild(gtk.Widgetter)
}

func Appends(parent Appender, childs ...gtk.Widgetter) {
	for _, ch := range childs {
		parent.Append(ch)
	}
}

func AddChilds(parent AddChilder, childs ...gtk.Widgetter) {
	for _, ch := range childs {
		parent.AddChild(ch)
	}
}
