package main

import (
	"os"

	_ "embed"

	"github.com/Tom5521/gtk4tools/pkg/gtools"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	app := gtk.NewApplication("com.test.window", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() {
		activate(app)
	})
	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

//go:embed window.ui
var uiFile string

type ui struct {
	Window *gtk.Window `gtk:"window"`
}

func activate(app *gtk.Application) {
	ui := new(ui)
	gtools.FetchObjects(ui, uiFile)
	app.AddWindow(ui.Window)

	ui.Window.Present()
}
