package ui

var FocusMap map[string]bool

func GenerateFocusMap() {
	FocusMap = make(map[string]bool)
	FocusMap["Playlist"] = true
	FocusMap["FileBrowser"] = false
	FocusMap["SearchView"] = false
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
