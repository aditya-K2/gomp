package client

import (
	"strings"

	"github.com/aditya-K2/fuzzy"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

func GetCell(text string, foreground tcell.Color, bold bool) *tview.TableCell {
	return tview.NewTableCell(text).
		SetAlign(tview.AlignLeft).
		SetStyle(tcell.StyleDefault.
			Foreground(foreground).
			Background(tcell.ColorBlack).
			Bold(bold))
}

func UpdateBuffSearchView(inputTable *tview.Table, m fuzzy.Matches, f []FileNode) {
	inputTable.Clear()
	if m == nil || len(m) == 0 {
		Update(f, inputTable)
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

func UpdatePlaylist(inputTable *tview.Table) {
	_playlistAttr, _ := CONN.PlaylistInfo(-1, -1)

	inputTable.Clear()
	for i, j := range _playlistAttr {
		_, _, w, _ := inputTable.GetInnerRect()
		if j["Title"] == "" || j["Artist"] == "" || j["Album"] == "" {
			inputTable.SetCell(i, 0,
				GetCell(
					utils.GetFormattedString(j["file"], w/3), tcell.ColorBlue, true))

		} else {
			inputTable.SetCell(i, 0,
				GetCell(
					utils.GetFormattedString(j["Title"], w/3), tcell.ColorGreen, false))
			inputTable.SetCell(i, 1,
				GetCell(
					utils.GetFormattedString(j["Artist"], w/3), tcell.ColorPurple, false))
			inputTable.SetCell(i, 2,
				GetCell(j["Album"], tcell.ColorYellow, false))
		}
	}
}

//  UpdateSearchView as the name suggests Updates the Search View the idea is to basically keep a fourth option called
//  Search in the Navigation bar which will render things from a global ContentSlice at least in the context of the main
//  function this will also help in persisting the Search Results.
func UpdateSearchView(inputTable *tview.Table, c []interface{}) {
	inputTable.Clear()
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
				if !strings.HasPrefix(b, WHITE_AND_BOLD) {
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

func Update(f []FileNode, inputTable *tview.Table) {
	inputTable.Clear()
	for i, j := range f {
		if len(j.Children) == 0 {
			_songAttributes, err := CONN.ListAllInfo(j.AbsolutePath)
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
