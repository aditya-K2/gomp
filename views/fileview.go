package views

import (
	"fmt"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/notify"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

type FileView struct {
}

func (f FileView) GetViewName() string {
	return "FileView"
}

func (f FileView) ShowChildrenContent() {
	UI := ui.Ui
	CONN := client.Conn
	r, _ := UI.ExpandedView.GetSelection()
	SetCurrentView(FView)
	if len(client.DirTree.Children[r].Children) == 0 {
		if id, err := CONN.AddID(client.DirTree.Children[r].AbsolutePath, -1); err != nil {
			notify.Notify.Send(fmt.Sprintf("Could not Add Song %s",
				client.DirTree.Children[r].Path))
		} else {
			if err := CONN.PlayID(id); err != nil {
				notify.Notify.Send(fmt.Sprintf("Could Not Play Song %s",
					client.DirTree.Children[r].Path))
			}
		}
	} else {
		client.DirTree = &client.DirTree.Children[r]
		FView.Update(UI.ExpandedView)
		UI.ExpandedView.Select(0, 0)
	}
}

func (f FileView) ShowParentContent() {
	UI := ui.Ui
	if client.DirTree.Parent != nil {
		client.DirTree = client.DirTree.Parent
		FView.Update(UI.ExpandedView)
	}
}

func (f FileView) AddToPlaylist() {
	UI := ui.Ui
	CONN := client.Conn
	r, _ := UI.ExpandedView.GetSelection()
	if err := CONN.Add(client.DirTree.Children[r].AbsolutePath); err != nil {
		notify.Notify.Send(fmt.Sprintf("Could not add %s to the Playlist",
			client.DirTree.Children[r].Path))
	}
}

func (f FileView) Quit() {
	ui.Ui.App.Stop()
}

func (f FileView) FocusBuffSearchView() {
	UI := ui.Ui
	SetCurrentView(BuffSView)
	UI.App.SetFocus(UI.SearchBar)
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
						utils.GetFormattedString(j.Title, w/3), tcell.ColorGreen, false))
				inputTable.SetCell(i, 1,
					GetCell(
						utils.GetFormattedString(j.Artist, w/3), tcell.ColorPurple, false))
				inputTable.SetCell(i, 2,
					GetCell(j.Album, tcell.ColorYellow, false))
			} else if j.Title == "" {
				inputTable.SetCell(i, 0,
					GetCell(j.Path, tcell.ColorBlue, true))
			}
		} else {
			inputTable.SetCell(i, 0,
				GetCell(j.Path, tcell.ColorYellow, true))
		}
	}
}
