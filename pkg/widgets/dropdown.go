package widgets

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type DropDown struct {
	*gtk.DropDown
	Items []string
}

func NewDropDown(items []string) *DropDown {
	d := &DropDown{
		Items:    items,
		DropDown: gtk.NewDropDownFromStrings(items),
	}

	return d
}

func (d *DropDown) ConnectSelected(f func()) {
	d.ConnectAfter("notify::selected", f)
}
