package views

import (
	"github.com/aditya-K2/tview"
)

var (
	CurrentView View
	BuffSView   BuffSearchView
	SView       SearchView
	FView       FileView
	PView       PlaylistView
)

type View interface {
	Update(inputTable *tview.Table)
	ShowChildrenContent()
	ShowParentContent()
	AddToPlaylist()
	Quit()
	FocusBuffSearchView()
	DeleteSongFromPlaylist()
	GetViewName() string
}

func SetCurrentView(v View) {
	CurrentView = v
}

func GetCurrentView() View {
	return CurrentView
}
