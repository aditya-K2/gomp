package ui

import (
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

var FuncMap map[string]func()

func NewMainS() *tview.Table {
	mains := tview.NewTable()

	mains.SetBackgroundColor(tcell.ColorDefault)
	mains.SetBorderPadding(1, 1, 1, 1).SetBorder(true)
	mains.SetSelectable(true, false)

	mains.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		GetCurrentView().Update(mains)
		return mains.GetInnerRect()
	})

	FuncMap = GenerateFuncMap(client.Conn)

	// Input Handler
	mains.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if val, ok := config.KEY_MAP[int(e.Rune())]; ok {
			FuncMap[val]()
			return nil
		} else {
			if GetCurrentView().GetViewName() == "PlaylistView" {
				if e.Rune() == 'j' || e.Rune() == 'k' {
					if len(PView.Playlist) == 0 {
						SendNotification("Empty Playlist")
						return nil
					}
				}
			} else if GetCurrentView().GetViewName() == "SearchView" {
				if e.Rune() == 'j' || e.Rune() == 'k' {
					if client.SearchContentSlice == nil || len(client.SearchContentSlice) == 0 {
						SendNotification("No Search Results")
						return nil
					}
				}
			}
			return e
		}
	})

	mains.SetDoneFunc(func(e tcell.Key) {
		if e == tcell.KeyEscape {
			if GetCurrentView().GetViewName() == "BuffSearchView" {
				SetCurrentView(FView)
				Ui.SearchBar.SetText("")
				client.Matches = nil
			}
		}
	})

	return mains
}
