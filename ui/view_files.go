package ui

import (
	"fmt"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/utils"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

var (
	PosStack utils.Stack[int]
)

type FileView struct {
}

func (f FileView) Name() string {
	return "FileView"
}

func (f FileView) ShowChildrenContent() {
	CONN := client.Conn
	r, _ := Ui.MainS.GetSelection()
	SetCurrentView(FView)
	if len(client.DirTree.Children[r].Children) == 0 {
		if id, err := CONN.AddID(client.DirTree.Children[r].AbsolutePath, -1); err != nil {
			SendNotification(fmt.Sprintf("Could not Add Song %s",
				client.DirTree.Children[r].Path))
		} else {
			if err := CONN.PlayID(id); err != nil {
				SendNotification(fmt.Sprintf("Could Not Play Song %s",
					client.DirTree.Children[r].Path))
			}
		}
	} else {
		PosStack.Push(r)
		client.DirTree = &client.DirTree.Children[r]
		FView.Update(Ui.MainS)
		Ui.MainS.Select(0, 0)
	}
}

func (f FileView) ShowParentContent() {
	if client.DirTree.Parent != nil {
		var last bool = false
		r := PosStack.Top()
		PosStack.Pop()
		Ui.MainS.Select(r, 0)
		client.DirTree = client.DirTree.Parent
		if r == len(client.DirTree.Children)-1 {
			last = true
		}
		Ui.App.QueueEvent(tcell.NewEventKey(tcell.KeyRune, 'j', tcell.ModNone))
		if !last {
			Ui.App.QueueEvent(tcell.NewEventKey(tcell.KeyRune, 'k', tcell.ModNone))
		}
		FView.Update(Ui.MainS)
	}
}

func (f FileView) AddToPlaylist() {
	CONN := client.Conn
	r, _ := Ui.MainS.GetSelection()
	if err := CONN.Add(client.DirTree.Children[r].AbsolutePath); err != nil {
		SendNotification(fmt.Sprintf("Could not add %s to the Playlist",
			client.DirTree.Children[r].Path))
	}
}

func (f FileView) Quit() {
	Ui.App.Stop()
}

func (f FileView) FocusBuffSearchView() {
	SetCurrentView(BuffSView)
	Ui.App.SetFocus(Ui.SearchBar)
}

func (f FileView) DeleteSongFromPlaylist() {}

func (f FileView) Update(inputTable *tview.Table) {
	inputTable.Clear()
	for i, j := range client.DirTree.Children {
		if len(j.Children) == 0 {
			if j.Title != "" {
				_, _, w, _ := inputTable.GetInnerRect()
				inputTable.SetCell(i, 0,
					GetCell(
						utils.GetFormattedString(j.Title, w/3), clr.Track))
				inputTable.SetCell(i, 1,
					GetCell(
						utils.GetFormattedString(j.Artist, w/3), clr.Artist))
				inputTable.SetCell(i, 2,
					GetCell(j.Album, clr.Album))
			} else if j.Title == "" {
				inputTable.SetCell(i, 0,
					GetCell(j.Path, clr.File))
			}
		} else {
			inputTable.SetCell(i, 0,
				GetCell(j.Path, clr.Folder))
		}
	}
}
