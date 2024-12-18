package gtools

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

func NewFactorySetup(f func(listitem ListItem)) func(*glib.Object) {
	return func(obj *glib.Object) {
		f(obj.Cast().(ListItem))
	}
}

func NewFactoryBind(f func(listitem ListItem, pos int)) func(*glib.Object) {
	return func(obj *glib.Object) {
		i, ok := obj.Cast().(ListItem)
		if !ok {
			panic("Invalid interface conversion! glib.Objector -> ListItem")
		}
		f(i, int(i.Position()))
	}
}
