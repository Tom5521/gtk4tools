package gt4tools

import (
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type List struct {
	Widget    *gtk.ListBox
	Items     []string
	RowWidget func(string) gtk.Widgetter

	// Channels

	killUpdater chan bool
	isCleaning  chan bool

	// Preferences

	PrintWarnings bool
}

func NewList(items []string, rowWidget func(item string) gtk.Widgetter) *List {
	l := &List{
		Widget:    gtk.NewListBox(),
		Items:     items,
		RowWidget: rowWidget,
	}

	l.InitChans()
	l.InitAutoUpdater()
	return l
}

// I only recommend calling this function when you no longer need to update or modify your list,
// any call after this (if you do not re-initialize the channels) will mean crashes.
func (l *List) TerminateAllChans() {
	close(l.killUpdater)
	close(l.isCleaning)
}

func (l *List) InitChans() {
	l.killUpdater = make(chan bool)
	l.killUpdater = make(chan bool)
}

func (l *List) Refresh() {
	l.cleanWidgetRows()
	for _, i := range l.Items {
		row := gtk.NewListBoxRow()
		w := l.RowWidget(i)
		row.SetChild(w)
		l.Widget.Append(row)
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

// ONLY TO BE USED IN EXCEPTIONAL AND VERY CONTROLLED CASES.
//
// This function should only be used when for X or Y reason you can't/want to refresh the list yourself.
//
// If for some reason the goroutine (or not) and another refresh are done at the same time, this can happen to the app:
//
// - Freezing
//
// - Crashes
//
// - Unexpected errors
//
// Never but NEVER call two goroutines to refresh because the above mentioned will happen without exception.
func (l *List) InitAutoUpdater() {
	l.Refresh()
	l.killUpdater = make(chan bool)
	go func() {
		for {
			oldItems := l.Items
			select {
			case <-l.killUpdater:
				return
			case <-time.After(50 * time.Millisecond):
				if !reflect.DeepEqual(oldItems, l.Items) {
					l.Refresh()
				}
			}
		}
	}()
}

func (l *List) KillAutoUpdater() {
	l.killUpdater <- true
}

func (l *List) cleanWidgetRows() {
	go func() {
		select {
		case <-l.isCleaning:
			l.printWarn("Another cleaning[refresh] is already underway")
			return
		default:
			l.isCleaning <- true
			for {
				row := l.Widget.RowAtIndex(0)
				if row == nil {
					break
				}
				l.Widget.Remove(row)
			}
			_ = <-l.isCleaning
			l.isCleaning <- false
		}
	}()
}

// Prints warnings if PrintWarnings is true
func (l *List) printWarn(warn ...any) {
	if l.PrintWarnings {
		fmt.Print("WARNING: ")
		fmt.Println(warn...)
	}
}
