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
type Person struct {
	Name string
	Age  uint
}

func activate(app *gtk.Application) {
	w := gtk.NewApplicationWindow(app)

	items := []Person{
		Person{
			Name: "Jonh Doe 1",
			Age:  21,
		},
		Person{
			Name: "Carlos Gimenez",
			Age:  26,
		},
		Person{
			Name: "Caroline Simpson",
			Age:  20,
		},
	}

	var personNames []string
	for _, p := range items {
		personNames = append(personNames, p.Name)
	}

	model := gtk.NewStringList(personNames)
	selectionModel := gtk.NewSingleSelection(model)
	selectionModel.ConnectSelectionChanged(func(_, _ uint) {
		fmt.Println("Index: ", selectionModel.Selected())
		fmt.Println("Value: ", personNames[selectionModel.Selected()])
	})

	factory := gtk.NewSignalListItemFactory()
	factory.ConnectSetup(func(listitem *gtk.ListItem) {
		listitem.SetChild(gtk.NewLabel(""))
	})
	factory.ConnectBind(func(listitem *gtk.ListItem) {
		obj := listitem.Item().Cast().(*gtk.StringObject)
		listitem.Child().(*gtk.Label).SetText(obj.String())
	})

	list := gtk.NewListView(selectionModel, &factory.ListItemFactory)

	w.SetChild(list)
	w.Show()
}
```

### After

```go
type Person struct {
	Name string
	Age  uint
}

func activate(app *gtk.Application) {
	w := gtk.NewApplicationWindow(app)

	items := []Person{
		Person{
			Name: "Jonh Doe 1",
			Age:  21,
		},
		Person{
			Name: "Carlos Gimenez",
			Age:  26,
		},
		Person{
			Name: "Caroline Simpson",
			Age:  20,
		},
	}

	list := widgets.NewList[Person](
		items,
		widgets.SelectionSingle,
		func(li *gtk.ListItem) {
			li.SetChild(gtk.NewLabel(""))
		},
		func(li *gtk.ListItem, p Person) {
			li.Child().(*gtk.Label).SetText(p.Name)
		},
	)
	list.OnSelected = func(index int) {
		fmt.Println("Index: ", index)
		fmt.Println("Value: ", items[index])
	}

	w.SetChild(list)
	w.Show()
}
```

You can test it by running `go run -v ./examples/Boxes/main.go`

## Documentation

The documentation is [here](https://pkg.go.dev/github.com/Tom5521/gtk4tools)

It took some time to appear in pkg.go.dev

Note: it is better to clone the repository and use godoc on it,
since pkg.go.dev takes a long time to index new versions.

## License

[MIT](https://choosealicense.com/licenses/mit/)
