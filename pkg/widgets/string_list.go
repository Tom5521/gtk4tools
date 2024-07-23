package widgets

// Deprecated: Replace it with CustomList[string].
// Since version 1.6.1 this struct is renamed from List to StringList.
type StringList struct {
	*List[string]
}

// Creates a new list that keeps the self.Items[] updated with that of the UI.
//
// Deprecated: Replace it with NewCustomList[string]
// As of version 1.6.0 this is simply a wrapper for CustomList[string],
// so it is recommended to simply change it to this type.
// There should be no further incompatibility with respect to this migration,
// so it should not be a problem.
func NewStringList(
	items []string,
	smodel ListSelectionMode,
	setup FactorySetup,
	bind FactoryBind[string],
) *StringList {
	return &StringList{
		List: NewList(
			items,
			smodel,
			setup,
			bind,
		),
	}
}
