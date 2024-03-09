# GTK4Tools

Basic abstractions and widgets for gotk4 library

These abstractions are mainly intended to simplify certain methods and make them
resemble the Fyne framework, which in my opinion are much more readable
than certain default GTK4 methods.

## Usage/Examples

You can use it in your project by importing it with

`go get github.com/Tom5521/gtk4tools@latest`

This library is to simplify gtk4 and avoid declaring unwanted variables that
take up possible names for more useful variables.

Here is a Before/After of applying the library.

### Before

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

 buttonsBox := gtk.NewBox(gtk.OrientationVertical,4)
 for _,b := range buttons{
  buttonsBox.Append(b)
 }
 buttonSbox := gtk.NewScrolledWindow()
 buttonSbox.SetChild(buttonsBox)

 labelsBox := gtk.NewBox(gtk.OrientationVertical,4)
 for _,l := range labels{
  labelsBox.Append(l)
 }
 labelSbox := gtk.NewScrolledWindow()
 labelSbox.SetChild(labelsBox)

 vbox := gtk.NewBox(gtk.OrientationVertical,1)
 vbox.SetHomogeneous(true)
 vbox.Append(buttonSbox)
 vbox.Append(labelSbox)

 w.SetChild(vbox)
 w.Show()
}
```

### After

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

 vbox := boxes.NewCHbox(1,
  boxes.NewScrolledVbox(
   // Convert a slice of a specific type to a gtk.Widgetter slice.
   t.ToWidgetter(buttons...)...,
  ),
  boxes.NewScrolledVbox(
   labels...,
  ),
 )
 vbox.SetHomogeneous(true)

 w.SetChild(vbox)
 w.Show()
}
```

You can test it by running `go run -v example.go`

## Documentation

The documentation is [here](https://pkg.go.dev/github.com/Tom5521/gtk4tools)

It took some time to appear in pkg.go.dev

## License

[MIT](https://choosealicense.com/licenses/mit/)
