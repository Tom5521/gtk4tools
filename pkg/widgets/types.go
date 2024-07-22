package widgets

import (
	"github.com/Tom5521/gtk4tools/pkg/gtools"
)

type FactorySetup func(gtools.ListItem)

type ListSelectionMode int

const (
	SelectionNone ListSelectionMode = iota
	SelectionSingle
	SelectionMultiple
)
