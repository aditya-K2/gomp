package views

import (
	"strings"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/ui/notify"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

type SearchView struct {
}

func (s SearchView) GetViewName() string {
	return "SearchView"
}
func (s SearchView) ShowChildrenContent() {
	UI := ui.Ui
	SearchContentSlice := client.SearchContentSlice
	if len(client.SearchContentSlice) <= 0 || client.SearchContentSlice == nil {
		notify.Send("No Search Results")
	} else {
		r, _ := UI.ExpandedView.GetSelection()
		client.AddToPlaylist(SearchContentSlice[r], true)
	}
}

func (s SearchView) ShowParentContent() {
	notify.Send("Not Allowed in this View")
	return
}

func (s SearchView) AddToPlaylist() {
	UI := ui.Ui
	SearchContentSlice := client.SearchContentSlice
	if len(client.SearchContentSlice) <= 0 || client.SearchContentSlice == nil {
		notify.Send("No Search Results")
	} else {
		r, _ := UI.ExpandedView.GetSelection()
		client.AddToPlaylist(SearchContentSlice[r], false)
	}
}

func (p SearchView) Quit() {
	ui.Ui.App.Stop()
}

func (s SearchView) FocusBuffSearchView()    {}
func (s SearchView) DeleteSongFromPlaylist() {}

func (s SearchView) Update(inputTable *tview.Table) {
	inputTable.Clear()
	c := client.SearchContentSlice
	_, _, width, _ := inputTable.GetInnerRect()
	for i, content := range c {
		switch content.(type) {
		case [3]string:
			{
				inputTable.SetCell(i, 0,
					GetCell(
						utils.GetFormattedString(content.([3]string)[0], width/3), tcell.ColorGreen, false))
				inputTable.SetCell(i, 1,
					GetCell(
						utils.GetFormattedString(content.([3]string)[1], width/3), tcell.ColorPurple, false))
				inputTable.SetCell(i, 2,
					GetCell(content.([3]string)[2], tcell.ColorYellow, false))
			}
		case [2]string:
			{
				inputTable.SetCell(i, 0,
					GetCell(
						utils.GetFormattedString(content.([2]string)[0], width/3), tcell.ColorYellow, false))
				inputTable.SetCell(i, 1,
					GetCell(
						utils.GetFormattedString(content.([2]string)[1], width/3), tcell.ColorPurple, false))
			}
		case string:
			{
				b := content.(string)
				if !strings.HasPrefix(b, client.WHITE_AND_BOLD) {
					inputTable.SetCell(i, 0,
						GetCell(content.(string), tcell.ColorPurple, false))
				} else {
					inputTable.SetCell(i, 0,
						GetCell(content.(string), tcell.ColorWhite, true).SetSelectable(false))
				}
			}
		}
	}
}
