package wcreate

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

func CreateButtons(btnTxt ...string) []*gtk.Button {
	var btns []*gtk.Button
	for _, t := range btnTxt {
		btns = append(btns, gtk.NewButtonWithLabel(t))
	}
	return btns
}

func CreateLabels(lbs ...string) []*gtk.Label {
	var labels []*gtk.Label
	for _, l := range lbs {
		labels = append(labels, gtk.NewLabel(l))
	}
	return labels
}
