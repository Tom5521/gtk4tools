package boxes

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type ScrolledBox struct {
	*gtk.ScrolledWindow

	Orientation gtk.Orientation
	Spacing     int
	Child       *gtk.Box
}

func (s *ScrolledBox) SetChild(child *gtk.Box) {
	s.Child = child
	s.ScrolledWindow.SetChild(child)
}

func (s *ScrolledBox) SetSpacing(spacing int) {
	s.Child.SetSpacing(spacing)
	s.Spacing = spacing
}

// Set the orientation of the child.
func (s *ScrolledBox) SetOrientation(orientation gtk.Orientation) {
	s.Orientation = orientation
	s.Child.SetOrientation(orientation)
}

func NewScrolledCVbox(spacing int, widgets ...gtk.Widgetter) *ScrolledBox {
	vbox := NewCVbox(spacing, widgets...)
	sbox := &ScrolledBox{
		Orientation:    vbox.Orientation(),
		ScrolledWindow: gtk.NewScrolledWindow(),
	}

	sbox.SetChild(vbox)
	return sbox
}

// Creates a vertical box that is scrollable in the X and Y axes, the orientation is that of the child.
func NewScrolledVbox(widgets ...gtk.Widgetter) *ScrolledBox {
	return NewScrolledCVbox(DefaultSpacing, widgets...)
}

func NewScrolledCHbox(spacing int, widgets ...gtk.Widgetter) *ScrolledBox {
	hbox := NewCHbox(spacing, widgets...)
	sbox := &ScrolledBox{
		Orientation:    hbox.Orientation(),
		ScrolledWindow: gtk.NewScrolledWindow(),
	}

	sbox.SetChild(hbox)
	return sbox
}

// Creates a horizontal box that is scrollable in the X and Y axes, the orientation is that of the child.
func NewScrolledHbox(widgets ...gtk.Widgetter) *ScrolledBox {
	return NewScrolledCHbox(DefaultSpacing, widgets...)
}

func NewScrolled(child *gtk.Box) *ScrolledBox {
	s := &ScrolledBox{
		Child:       child,
		Orientation: child.Orientation(),
		Spacing:     child.Spacing(),
	}

	return s
}
