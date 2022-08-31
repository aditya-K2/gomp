package views

import (
	"fmt"

	"github.com/aditya-K2/gomp/globals"
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
	UI := globals.Ui
	CONN := globals.Conn
	r, _ := UI.ExpandedView.GetSelection()
	SetCurrentView(FView)
	if len(globals.DirTree.Children[r].Children) == 0 {
		if id, err := CONN.AddId(globals.DirTree.Children[r].AbsolutePath, -1); err != nil {
			globals.Notify.Send(fmt.Sprintf("Could not Add Song %s",
				globals.DirTree.Children[r].Path))
		} else {
			if err := CONN.PlayId(id); err != nil {
				globals.Notify.Send(fmt.Sprintf("Could Not Play Song %s",
					globals.DirTree.Children[r].Path))
			}
		}
	} else {
		globals.DirTree = &globals.DirTree.Children[r]
		FView.Update(UI.ExpandedView)
		UI.ExpandedView.Select(0, 0)
	}
}

func (f FileView) ShowParentContent() {
	UI := globals.Ui
	if globals.DirTree.Parent != nil {
		globals.DirTree = globals.DirTree.Parent
		FView.Update(UI.ExpandedView)
	}
}

func (f FileView) AddToPlaylist() {
	UI := globals.Ui
	CONN := globals.Conn
	r, _ := UI.ExpandedView.GetSelection()
	if err := CONN.Add(globals.DirTree.Children[r].AbsolutePath); err != nil {
		globals.Notify.Send(fmt.Sprintf("Could not add %s to the Playlist",
			globals.DirTree.Children[r].Path))
	}
}

func (f FileView) Quit() {
	globals.Ui.App.Stop()
}

func (f FileView) FocusBuffSearchView() {
	UI := globals.Ui
	SetCurrentView(BuffSView)
	UI.App.SetFocus(UI.SearchBar)
}

func (f FileView) DeleteSongFromPlaylist() {}

func (f FileView) Update(inputTable *tview.Table) {
	inputTable.Clear()
	for i, j := range globals.DirTree.Children {
		if len(j.Children) == 0 {
			_songAttributes, err := globals.Conn.ListAllInfo(j.AbsolutePath)
			if err == nil && _songAttributes[0]["Title"] != "" {
				_, _, w, _ := inputTable.GetInnerRect()
				inputTable.SetCell(i, 0,
					GetCell(
						utils.GetFormattedString(_songAttributes[0]["Title"], w/3), tcell.ColorGreen, false))

				inputTable.SetCell(i, 1,
					GetCell(
						utils.GetFormattedString(_songAttributes[0]["Artist"], w/3), tcell.ColorPurple, false))
				inputTable.SetCell(i, 2,
					GetCell(_songAttributes[0]["Album"], tcell.ColorYellow, false))

			} else if _songAttributes[0]["Title"] == "" {
				inputTable.SetCell(i, 0,
					GetCell(j.Path, tcell.ColorBlue, true))
			}
		} else {
			inputTable.SetCell(i, 0,
				GetCell(j.Path, tcell.ColorYellow, true))
		}
	}
}
