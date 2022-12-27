package views

import (
	"fmt"
	"time"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/ui/notify"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

type MostPlayedView struct {
	FMap   map[string]time.Duration
	FSlice []string
}

func (s MostPlayedView) GetViewName() string {
	return "MostPlayedView"
}
func (s MostPlayedView) ShowChildrenContent() {
	if len(s.FSlice) <= 0 || s.FSlice == nil {
		notify.Send("No Search Results")
	} else {
		r, _ := ui.Ui.ExpandedView.GetSelection()
		if id, err := client.Conn.AddID(s.FSlice[r], -1); err != nil {
			notify.Send(fmt.Sprintf("Couldn't Add %s to playlist", s.FSlice[r]))
		} else {
			client.Conn.PlayID(id)
		}
	}
}

func (s MostPlayedView) ShowParentContent() {
	notify.Send("Not Allowed in this View")
	return
}

func (s MostPlayedView) AddToPlaylist() {
	UI := ui.Ui
	if len(s.FSlice) <= 0 || s.FSlice == nil {
		notify.Send("No Search Results")
	} else {
		r, _ := UI.ExpandedView.GetSelection()
		client.AddToPlaylist(s.FSlice[r], false)
	}
}

func (p MostPlayedView) Quit() {
	ui.Ui.App.Stop()
}

func (s MostPlayedView) FocusBuffSearchView()    {}
func (s MostPlayedView) DeleteSongFromPlaylist() {}

func (s MostPlayedView) Update(inputTable *tview.Table) {
	n := 4
	inputTable.Clear()
	c := s.FSlice
	_, _, width, _ := inputTable.GetInnerRect()
	i := 0
	for _, path := range c {
		attr, err := client.Conn.ListAllInfo(path)
		if err != nil {
			utils.Pop(i, c)
			continue
		}
		if len(attr) != 0 {
			if attr[0]["Title"] != "" {
				inputTable.SetCell(i, 0,
					GetCell(
						utils.GetFormattedString(attr[0]["Title"], width/n), tcell.ColorGreen, false))
				inputTable.SetCell(i, 1,
					GetCell(
						utils.GetFormattedString(attr[0]["Artist"], width/n), tcell.ColorPurple, false))
				inputTable.SetCell(i, 2,
					GetCell(
						utils.GetFormattedString(attr[0]["Album"], width/n), tcell.ColorYellow, false))
				inputTable.SetCell(i, 3,
					GetCell(s.FMap[path].String(), tcell.ColorGreen, false))
			} else {
				inputTable.SetCell(i, 0,
					GetCell(
						utils.GetFormattedString(attr[0]["Album"], width/n), tcell.ColorBlue, true))
				inputTable.SetCell(i, 3,
					GetCell(s.FMap[path].String(), tcell.ColorGreen, false))
			}
		}
		i++
	}
}
