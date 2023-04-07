package ui

import (
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/tview"
)

var (
	clr         = config.Config.Colors
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
	Name() string
}

func SetCurrentView(v View) {
	CurrentView = v
}

func GetCurrentView() View {
	return CurrentView
}

func GetCell(text string, color config.Color) *tview.TableCell {
	return tview.NewTableCell(text).
		SetAlign(tview.AlignLeft).
		SetStyle(color.Style())
}
