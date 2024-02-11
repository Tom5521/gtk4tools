package widgets

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type List struct {
	*gtk.ListBox
	Items     []string
	RowWidget func(int, string) gtk.Widgetter

	// default: 50 miliseconds
	UpdaterSleepTime time.Duration

	// Channels

	killUpdater bool

	clean   sync.Mutex
	refresh sync.Mutex

	// Preferences

	PrintWarnings bool
}

func NewList(items []string, rowWidget func(index int, item string) gtk.Widgetter) *List {
	l := &List{
		ListBox:          gtk.NewListBox(),
		Items:            items,
		RowWidget:        rowWidget,
		UpdaterSleepTime: 50 * time.Millisecond,
	}
	l.InitAutoUpdater()
	return l
}

func (l *List) CleanUp() {
	l.KillAutoUpdater()
}

func (l *List) Refresh() {
	l.refresh.Lock()
	defer l.refresh.Unlock()
	l.cleanWidgetRows()
	for index, item := range l.Items {
		row := gtk.NewListBoxRow()
		w := l.RowWidget(index, item)
		row.SetChild(w)
		l.Append(row)
	}
}

func (l *List) RemoveItem(index int) {
	if !(index >= 0 && index < len(l.Items)) {
		l.printWarn("The index does not exist")
		return
	}
	l.Items = slices.Delete(l.Items, index, index+1)
	l.Refresh()
}

func (l *List) InitAutoUpdater() {
	l.killUpdater = false
	l.Refresh()
	go func() {
		for {
			switch l.killUpdater {
			case true:
				return
			default:
				oldItems := make([]string, len(l.Items))
				copy(oldItems, l.Items)
				time.Sleep(l.UpdaterSleepTime)
				if !slices.Equal(l.Items, oldItems) {
					l.Refresh()
				}
			}
		}
	}()
}

func (l *List) KillAutoUpdater() {
	l.killUpdater = true
}

func (l *List) cleanWidgetRows() {
	l.clean.Lock()
	defer l.clean.Unlock()
	for {
		row := l.RowAtIndex(0)
		if row == nil {
			break
		}
		l.Remove(row)
	}
}

// Prints warnings if PrintWarnings is true.
func (l *List) printWarn(warn ...any) {
	if l.PrintWarnings {
		fmt.Print("WARNING: ")
		fmt.Println(warn...)
	}
}
