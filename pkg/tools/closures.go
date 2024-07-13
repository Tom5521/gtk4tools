package tools

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func NewFactorySetup(f func(listitem *gtk.ListItem)) func(obj *glib.Object) {
	return func(obj *glib.Object) {
		f(obj.Cast().(*gtk.ListItem))
	}
}

func NewFactoryBind(f func(listitem *gtk.ListItem, pos int)) func(obj *glib.Object) {
	return func(obj *glib.Object) {
		i := obj.Cast().(*gtk.ListItem)
		f(i, int(i.Position()))
	}
}
