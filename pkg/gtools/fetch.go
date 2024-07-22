package gtools

import (
	"reflect"

	"github.com/Tom5521/gtk4tools/internal/walk"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// What this function does is to obtain the tags of the provided structure
// and with the content of the tag call gtk.Builder.GetObject
// and assign it to its corresponding field,
// if the field has no tag it will be skipped and if it is another sub struct
// it will continue until it reaches the end.
func FetchObjects(str any, builder string) {
	b := gtk.NewBuilderFromString(builder)
	walk.Into(str, func(v reflect.Value, sf reflect.StructField) bool {
		tag, ok := sf.Tag.Lookup("gtk")
		if !ok {
			return true
		}
		obj := b.GetObject(tag).Cast()
		v.Set(reflect.ValueOf(obj))
		return true
	})
}
