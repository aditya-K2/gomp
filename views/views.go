package views

import (
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

var (
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

func GetCell(text string, foreground tcell.Color, bold bool) *tview.TableCell {
	return tview.NewTableCell(text).
		SetAlign(tview.AlignLeft).
		SetStyle(tcell.StyleDefault.
			Foreground(foreground).
			Background(tcell.ColorBlack).
			Bold(bold))
}
