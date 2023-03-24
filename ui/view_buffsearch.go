package ui

import (
	"fmt"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/utils"
	"github.com/aditya-K2/tview"
)

type BuffSearchView struct {
}

func (s BuffSearchView) Name() string {
	return "BuffSearchView"
}

func (s BuffSearchView) ShowChildrenContent() {
	UI := Ui
	CONN := client.Conn
	r, _ := UI.MainS.GetSelection()
	SetCurrentView(FView)
	if len(client.DirTree.Children[client.Matches[r].Index].Children) == 0 {
		if id, err := CONN.AddID(client.DirTree.Children[client.Matches[r].Index].AbsolutePath, -1); err != nil {
			SendNotification(fmt.Sprintf("Could Not add the Song %s to the Playlist",
				client.DirTree.Children[client.Matches[r].Index].AbsolutePath))
		} else {
			if err := CONN.PlayID(id); err != nil {
				SendNotification("Could not Play the Song")
			}
		}
	} else {
		PosStack.Push(client.Matches[r].Index)
		client.DirTree = &client.DirTree.Children[client.Matches[r].Index]
		FView.Update(UI.MainS)
	}
	UI.SearchBar.SetText("")
	// Resetting client.Matches
	client.Matches = nil
}

func (s BuffSearchView) ShowParentContent() {
	SendNotification("Not Allowed in this View")
	return
}

func (s BuffSearchView) AddToPlaylist() {
	UI := Ui
	CONN := client.Conn
	r, _ := UI.MainS.GetSelection()
	if err := CONN.Add(client.DirTree.Children[client.Matches[r].Index].AbsolutePath); err != nil {
		SendNotification(fmt.Sprintf("Could Not Add URI %s to the Playlist",
			client.DirTree.Children[client.Matches[r].Index].Path))
	} else {
		SetCurrentView(FView)
		SendNotification(fmt.Sprintf("URI Added %s to the Playlist",
			client.DirTree.Children[client.Matches[r].Index].Path))
		SetCurrentView(BuffSView)
	}
}

func (s BuffSearchView) Quit() {
	SetCurrentView(FView)
	Ui.SearchBar.SetText("")
	client.Matches = nil
}

func (f BuffSearchView) FocusBuffSearchView() {
	UI := Ui
	SetCurrentView(BuffSView)
	UI.App.SetFocus(UI.SearchBar)
}

func (f BuffSearchView) DeleteSongFromPlaylist() {}

func (s BuffSearchView) Update(inputTable *tview.Table) {
	m := client.Matches
	f := client.DirTree.Children
	inputTable.Clear()
	if m == nil || len(m) == 0 {
		FView.Update(inputTable)
	} else {
		for k, v := range m {
			if len(f[v.Index].Children) != 0 {
				inputTable.SetCell(k, 0,
					GetCell(
						utils.GetMatchedString(
							utils.Unique(v.MatchedIndexes), f[v.Index].Path, clr.MatchedFolder.String()),
						clr.Folder))
			} else {
				inputTable.SetCell(k, 0,
					GetCell(
						utils.GetMatchedString(
							utils.Unique(v.MatchedIndexes), f[v.Index].Title, clr.MatchedTitle.String()),
						clr.Track))
			}
			if k == 15 {
				break
			}
		}
	}
}
