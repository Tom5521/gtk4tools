# Gtk4Tools

Basic abstractions and widgets for gotk4 library

These abstractions are mainly intended to simplify certain methods and make them
resemble the Fyne framework, which in my opinion are much more readable
than certain default GTK4 methods.

## Usage/Examples

You can use it in your project by importing it with

`go get github.com/Tom5521/gtk4tools@latest`

Example:

```go
func activate(app *gtk.Application) {
 w := gtk.NewApplicationWindow(app)
 w.SetDefaultSize(500, 400)

 var buttons []*gtk.Button
 for i := range 30 {
  buttons = append(buttons, gtk.NewButtonWithLabel("Button "+strconv.Itoa(i)))
 }

 var labels []gtk.Widgetter
 for i := range 30 {
  labels = append(labels, gtk.NewLabel("Label "+strconv.Itoa(i)))
 }

 vbox := boxes.NewHbox(
  boxes.NewScrolledVbox(
   // Convert a slice of a specific type to a gtk.Widgetter slice.
   t.ToWidgetter(buttons)...,
  ),
  boxes.NewScrolledVbox(
   labels...,
  ),
 )
 vbox.SetSpacing(1)
 vbox.SetHomogeneous(true)

 w.SetChild(vbox)
 w.Show()
}
```

You can test it by running `go run -v example.go`

## License

[MIT](https://choosealicense.com/licenses/mit/)
