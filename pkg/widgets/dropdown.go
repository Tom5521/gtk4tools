package widgets

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type DropDown struct {
	*gtk.DropDown
	Items []string

	OnSelected func(selected uint)
}

func NewDropDown(items []string) *DropDown {
	d := &DropDown{
		Items:    items,
		DropDown: gtk.NewDropDownFromStrings(items),
	}
	d.reConnectOnSelected()
	return d
}

func (d *DropDown) Refresh() {
	d.reConnectOnSelected()
}

// Internal functions

func (d *DropDown) reConnectOnSelected() {
	d.ConnectAfter("notify::selected", func() {
		if d.OnSelected == nil {
			return
		}
		d.OnSelected(d.Selected())
	})
}
