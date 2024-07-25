package gtools_test

import (
	"fmt"

	"github.com/Tom5521/gtk4tools/pkg/v2/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ExampleFetchObjects() {
	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<interface>
  <requires lib="gtk" version="4.0"/>
  <object class="GtkWindow" id="window">
    <property name="title">Hello World</property>
  </object>
</interface>`

	type ui struct {
		Window *gtk.Window `gtk:"window"`
	}

	w := new(ui)
	gtools.FetchObjects(w, xml)

	fmt.Println(w.Window.Title()) // Hello World
}
