package ui

import (
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

var (
	clr         = config.Config.Colors
	CurrentView View
	BuffSView   BuffSearchView
	SView       SearchView
	FView       FileView
	PView       PlaylistView
	MPView      MostPlayedView
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

func GetCell(text string, color config.Color) *tview.TableCell {
	return tview.NewTableCell(text).
		SetAlign(tview.AlignLeft).
		SetStyle(tcell.StyleDefault.
			Foreground(color.Color()).
			Background(tcell.ColorBlack).
			Bold(color.Bold).
			Italic(color.Italic))
}
