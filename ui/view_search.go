package ui

import (
	"fmt"
	"strings"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
)

type SearchView struct {
}

func addToPlaylist(a interface{}, addAndPlay bool) {
	switch a.(type) {
	case [3]string:
		{
			b := a.([3]string)
			if err := client.AddTitle(client.ArtistTree, b[1], b[2], b[0], addAndPlay); err != nil {
				SendNotification(err.Error())
			} else {
				SendNotification(fmt.Sprintf(
					"%s Added Succesfully!", b[0]))
			}
		}
	case [2]string:
		{
			b := a.([2]string)
			if err := client.AddAlbum(client.ArtistTree, b[0], b[1]); err != nil {
				SendNotification(err.Error())
			} else {
				SendNotification(fmt.Sprintf(
					"%s Added Succesfully!", b[0]))
			}
		}
	case string:
		{
			b := a.(string)
			if err := client.AddArtist(client.ArtistTree, b); err != nil {
				SendNotification(err.Error())
			} else {
				SendNotification(fmt.Sprintf(
					"%s Added Succesfully!", b))
			}
		}
	}
}

func (s SearchView) GetViewName() string {
	return "SearchView"
}
func (s SearchView) ShowChildrenContent() {
	SearchContentSlice := client.SearchContentSlice
	if len(client.SearchContentSlice) <= 0 || client.SearchContentSlice == nil {
		SendNotification("No Search Results")
	} else {
		r, _ := Ui.MainS.GetSelection()
		addToPlaylist(SearchContentSlice[r], true)
	}
}

func (s SearchView) ShowParentContent() {
	SendNotification("Not Allowed in this View")
	return
}

func (s SearchView) AddToPlaylist() {
	SearchContentSlice := client.SearchContentSlice
	if len(client.SearchContentSlice) <= 0 || client.SearchContentSlice == nil {
		SendNotification("No Search Results")
	} else {
		r, _ := Ui.MainS.GetSelection()
		addToPlaylist(SearchContentSlice[r], false)
	}
}

func (p SearchView) Quit() {
	Ui.App.Stop()
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
						utils.GetFormattedString(content.([3]string)[0], width/3), clr.Track))
				inputTable.SetCell(i, 1,
					GetCell(
						utils.GetFormattedString(content.([3]string)[1], width/3), clr.Artist))
				inputTable.SetCell(i, 2,
					GetCell(content.([3]string)[2], clr.Album))
			}
		case [2]string:
			{
				inputTable.SetCell(i, 0,
					GetCell(
						utils.GetFormattedString(content.([2]string)[0], width/3), clr.Album))
				inputTable.SetCell(i, 1,
					GetCell(
						utils.GetFormattedString(content.([2]string)[1], width/3), clr.Artist))
			}
		case string:
			{
				b := content.(string)
				if !strings.HasPrefix(b, client.WHITE_AND_BOLD) {
					inputTable.SetCell(i, 0,
						GetCell(content.(string), clr.Artist))
				} else {
					inputTable.SetCell(i, 0,
						GetCell(content.(string), clr.Null).SetSelectable(false))
				}
			}
		}
	}
}
