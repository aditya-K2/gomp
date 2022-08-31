package views

import (
	"fmt"

	"github.com/aditya-K2/gomp/globals"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

type BuffSearchView struct {
}

func (s BuffSearchView) GetViewName() string {
	return "BuffSearchView"
}

func (s BuffSearchView) ShowChildrenContent() {
	UI := globals.Ui
	CONN := globals.Conn
	r, _ := UI.ExpandedView.GetSelection()
	SetCurrentView(FView)
	if len(globals.DirTree.Children[r].Children) == 0 {
		if id, err := CONN.AddId(globals.DirTree.Children[globals.Matches[r].Index].AbsolutePath, -1); err != nil {
			globals.Notify.Send(fmt.Sprintf("Could Not add the Song %s to the Playlist",
				globals.DirTree.Children[globals.Matches[r].Index].AbsolutePath))
		} else {
			if err := CONN.PlayId(id); err != nil {
				globals.Notify.Send("Could not Play the Song")
			}
		}
	} else {
		globals.DirTree = &globals.DirTree.Children[globals.Matches[r].Index]
		FView.Update(UI.ExpandedView)
	}
	UI.SearchBar.SetText("")
	// Resetting globals.Matches
	globals.Matches = nil
}

func (s BuffSearchView) ShowParentContent() {
	globals.Notify.Send("Not Allowed in this View")
	return
}

func (s BuffSearchView) AddToPlaylist() {
	UI := globals.Ui
	CONN := globals.Conn
	r, _ := UI.ExpandedView.GetSelection()
	if err := CONN.Add(globals.DirTree.Children[globals.Matches[r].Index].AbsolutePath); err != nil {
		globals.Notify.Send(fmt.Sprintf("Could Not Add URI %s to the Playlist",
			globals.DirTree.Children[globals.Matches[r].Index].Path))
	} else {
		SetCurrentView(FView)
		globals.Notify.Send(fmt.Sprintf("URI Added %s to the Playlist",
			globals.DirTree.Children[globals.Matches[r].Index].Path))
		SetCurrentView(BuffSView)
	}
}

func (s BuffSearchView) Quit() {
	UI := globals.Ui
	SetCurrentView(FView)
	UI.SearchBar.SetText("")
	globals.Matches = nil
}

func (f BuffSearchView) FocusBuffSearchView() {
	UI := globals.Ui
	SetCurrentView(BuffSView)
	UI.App.SetFocus(UI.SearchBar)
}

func (f BuffSearchView) DeleteSongFromPlaylist() {}

func (s BuffSearchView) Update(inputTable *tview.Table) {
	m := globals.Matches
	f := globals.DirTree.Children
	inputTable.Clear()
	if m == nil || len(m) == 0 {
		FView.Update(inputTable)
	} else {
		for k, v := range m {
			if len(f[v.Index].Children) != 0 {
				inputTable.SetCell(k, 0,
					GetCell(
						utils.GetMatchedString(
							utils.Unique(v.MatchedIndexes), f[v.Index].Path, "[blue:-:bi]"),
						tcell.ColorYellow, true))
			} else {
				inputTable.SetCell(k, 0,
					GetCell(
						utils.GetMatchedString(
							utils.Unique(v.MatchedIndexes), f[v.Index].Title, "[yellow:-:bi]"),
						tcell.ColorGreen, true))
			}
			if k == 15 {
				break
			}
		}
	}
}
