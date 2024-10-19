package gtools_test

import (
	"testing"

	"github.com/Tom5521/gtk4tools/pkg/v2/gtools"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func TestInitGTK(t *testing.T) {
	if !gtk.InitCheck() {
		t.Log("Failed to init GTK")
		t.FailNow()
	}
}

func TestFetch(t *testing.T) {
	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<interface>
  <requires lib="gtk" version="4.0"/>
  <object class="GtkWindow" id="window">
    <property name="title">Hello World</property>
  </object>
</interface>`

	ui := new(
		struct {
			Window *gtk.Window `gtk:"window"`
		},
	)

	gtools.FetchObjects(ui, xml)

	if ui.Window == nil {
		t.Log("ui.Window is nil")
		t.FailNow()
	}
	if ui.Window.Title() != "Hello World" {
		t.Log("ui.Window.Title() isn't 'Hello World'")
		t.Fail()
	}
}

var (
	modeller gtools.Modeller[string]
	items    = []string{"1", "2", "3", "4"}
)

func testModeller(t *testing.T) {
	if len(items) != modeller.Len() {
		t.Log("Modeller len is greater than slices len")
		t.FailNow()
	}

	if items[1] != modeller.At(1) {
		t.Log("Modeller item 1 isn't equal to slice item 1")
		t.FailNow()
	}
}

func TestModel(t *testing.T) {
	modeller = gtools.NewModel(items...)
	testModeller(t)
}

func TestModelVar(t *testing.T) {
	modeller = gtools.NewModelVar(&items)
	testModeller(t)
}
