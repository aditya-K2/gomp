package ui

// The Focus Map Helps to keep track of which UI Element Currently Has the Focus It can be queried to get the Current
// UI Element with Focus and also can set UI Focus keep in mind that it isn't Focus Map that is Responsible to change
// the Focus that is Done through the Update Function of UI.ExpandedView */
var FocusMap map[string]bool

func GenerateFocusMap() {
	FocusMap = make(map[string]bool)
	FocusMap["Playlist"] = true
	FocusMap["FileBrowser"] = false
	FocusMap["SearchView"] = false
	FocusMap["BuffSearchView"] = false
}

func HasFocus(s string) bool {
	return FocusMap[s]
}

func SetFocus(s string) {
	for k := range FocusMap {
		FocusMap[k] = false
	}
	FocusMap[s] = true
}
