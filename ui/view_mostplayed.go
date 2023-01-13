package ui

import (
	"fmt"
	"time"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
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
		SendNotification("No Search Results")
	} else {
		r, _ := Ui.MainS.GetSelection()
		if id, err := client.Conn.AddID(s.FSlice[r], -1); err != nil {
			SendNotification(fmt.Sprintf("Couldn't Add %s to playlist", s.FSlice[r]))
		} else {
			client.Conn.PlayID(id)
		}
	}
}

func (s MostPlayedView) ShowParentContent() {
	SendNotification("Not Allowed in this View")
	return
}

func (s MostPlayedView) AddToPlaylist() {
	if len(s.FSlice) <= 0 || s.FSlice == nil {
		SendNotification("No Search Results")
	} else {
		r, _ := Ui.MainS.GetSelection()
		addToPlaylist(s.FSlice[r], false)
	}
}

func (p MostPlayedView) Quit() {
	Ui.App.Stop()
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
						utils.GetFormattedString(attr[0]["Title"], width/n), clr.Track))
				inputTable.SetCell(i, 1,
					GetCell(
						utils.GetFormattedString(attr[0]["Artist"], width/n), clr.Artist))
				inputTable.SetCell(i, 2,
					GetCell(
						utils.GetFormattedString(attr[0]["Album"], width/n), clr.Album))
				inputTable.SetCell(i, 3,
					GetCell(s.FMap[path].String(), clr.Timestamp))
			} else {
				inputTable.SetCell(i, 0,
					GetCell(
						utils.GetFormattedString(attr[0]["file"], width/n), clr.File))
				inputTable.SetCell(i, 3,
					GetCell(s.FMap[path].String(), clr.Timestamp))
			}
		}
		i++
	}
}
