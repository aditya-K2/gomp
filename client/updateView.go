package client

import (
	"strings"

	"github.com/aditya-K2/fuzzy"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
)

func UpdateBuffSearchView(inputTable *tview.Table, m fuzzy.Matches, f []FileNode) {
	inputTable.Clear()
	if m == nil || len(m) == 0 {
		Update(f, inputTable)
	} else {
		for k, v := range m {
			if len(f[v.Index].Children) != 0 {
				inputTable.SetCellSimple(k, 0, utils.GetMatchedString(f[v.Index].Path, "#0000ff", "yellow", v.MatchedIndexes))
			} else {
				inputTable.SetCellSimple(k, 0, utils.GetMatchedString(f[v.Index].Title, "#fbff00", "green", v.MatchedIndexes))
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
			inputTable.SetCell(i, 0, tview.NewTableCell(utils.GetFormattedString(j["file"], w/3)))
		} else {
			inputTable.SetCell(i, 0, tview.NewTableCell(utils.GetFormattedString("[green]"+j["Title"], w/3)))
			inputTable.SetCell(i, 1, tview.NewTableCell(utils.GetFormattedString("[magenta]"+j["Artist"], w/3)))
			inputTable.SetCell(i, 2, tview.NewTableCell("[yellow]"+j["Album"]))
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
				inputTable.SetCell(i, 0, tview.NewTableCell(utils.GetFormattedString("[green]"+content.([3]string)[0], width/3)))
				inputTable.SetCell(i, 1, tview.NewTableCell(utils.GetFormattedString("[magenta]"+content.([3]string)[1], width/3)))
				inputTable.SetCell(i, 2, tview.NewTableCell(utils.GetFormattedString("[yellow]"+content.([3]string)[2], width/3)))
			}
		case [2]string:
			{
				inputTable.SetCell(i, 0, tview.NewTableCell(utils.GetFormattedString("[green]"+content.([2]string)[0], width/3)))
				inputTable.SetCell(i, 1, tview.NewTableCell(utils.GetFormattedString("[magenta]"+content.([2]string)[1], width/3)))
			}
		case string:
			{
				b := content.(string)
				if !strings.HasPrefix(b, WHITE_AND_BOLD) {
					inputTable.SetCell(i, 0, tview.NewTableCell("[green]"+content.(string)))
				} else {
					inputTable.SetCell(i, 0, tview.NewTableCell(content.(string)).SetSelectable(false))
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
					tview.NewTableCell("[green]"+utils.GetFormattedString(_songAttributes[0]["Title"], w/3)).
						SetAlign(tview.AlignLeft))

				inputTable.SetCell(i, 1,
					tview.NewTableCell("[magenta]"+utils.GetFormattedString(_songAttributes[0]["Artist"], w/3)).
						SetAlign(tview.AlignLeft))

				inputTable.SetCell(i, 2,
					tview.NewTableCell("[yellow]"+_songAttributes[0]["Album"]).
						SetAlign(tview.AlignLeft))

			} else if _songAttributes[0]["Title"] == "" {
				inputTable.SetCell(i, 0,
					tview.NewTableCell("[blue]"+j.Path).
						SetAlign(tview.AlignLeft))
			}
		} else {
			inputTable.SetCell(i, 0,
				tview.NewTableCell("[yellow::b]"+j.Path).
					SetAlign(tview.AlignLeft))
		}
	}
}
