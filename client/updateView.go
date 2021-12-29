package client

import (
	"strings"

	"github.com/aditya-K2/fuzzy"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

func UpdateBuffSearchView(inputTable *tview.Table, m fuzzy.Matches, f []FileNode) {
	inputTable.Clear()
	if m == nil || len(m) == 0 {
		Update(f, inputTable)
	} else {
		for k, v := range m {
			if len(f[v.Index].Children) != 0 {
				inputTable.SetCell(k, 0, tview.NewTableCell(utils.GetMatchedString(utils.Unique(v.MatchedIndexes), f[v.Index].Path, "[#0000ff:-:bi]")).
					SetAlign(tview.AlignLeft).
					SetStyle(tcell.StyleDefault.
						Foreground(tcell.ColorYellow).
						Background(tcell.ColorBlack).
						Bold(true)))
			} else {
				inputTable.SetCell(k, 0, tview.NewTableCell(utils.GetMatchedString(utils.Unique(v.MatchedIndexes), f[v.Index].Title, "[#fbff00:-:bi]")).
					SetAlign(tview.AlignLeft).
					SetStyle(tcell.StyleDefault.
						Foreground(tcell.ColorGreen).
						Background(tcell.ColorBlack).
						Bold(true)))
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
			inputTable.SetCell(i, 0, tview.NewTableCell(utils.GetFormattedString(j["file"], w/3)).
				SetAlign(tview.AlignLeft).
				SetStyle(tcell.StyleDefault.
					Foreground(tcell.ColorBlue).
					Background(tcell.ColorBlack).
					Bold(true)))

		} else {
			inputTable.SetCell(i, 0, tview.NewTableCell(utils.GetFormattedString(j["Title"], w/3)).
				SetAlign(tview.AlignLeft).
				SetStyle(tcell.StyleDefault.
					Foreground(tcell.ColorGreen).
					Background(tcell.ColorBlack)))
			inputTable.SetCell(i, 1, tview.NewTableCell(utils.GetFormattedString(j["Artist"], w/3)).
				SetAlign(tview.AlignLeft).
				SetStyle(tcell.StyleDefault.
					Foreground(tcell.ColorPurple).
					Background(tcell.ColorBlack)))
			inputTable.SetCell(i, 2, tview.NewTableCell(j["Album"]).
				SetAlign(tview.AlignLeft).
				SetStyle(tcell.StyleDefault.
					Foreground(tcell.ColorYellow).
					Background(tcell.ColorBlack)))
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
				inputTable.SetCell(i, 0, tview.NewTableCell(utils.GetFormattedString(content.([3]string)[0], width/3)).
					SetAlign(tview.AlignLeft).
					SetStyle(tcell.StyleDefault.
						Foreground(tcell.ColorGreen).
						Background(tcell.ColorBlack)))
				inputTable.SetCell(i, 1, tview.NewTableCell(utils.GetFormattedString(content.([3]string)[1], width/3)).
					SetAlign(tview.AlignLeft).
					SetStyle(tcell.StyleDefault.
						Foreground(tcell.ColorPurple).
						Background(tcell.ColorBlack)))
				inputTable.SetCell(i, 2, tview.NewTableCell(utils.GetFormattedString(content.([3]string)[2], width/3)).
					SetAlign(tview.AlignLeft).
					SetStyle(tcell.StyleDefault.
						Foreground(tcell.ColorYellow).
						Background(tcell.ColorBlack)))
			}
		case [2]string:
			{
				inputTable.SetCell(i, 0, tview.NewTableCell(utils.GetFormattedString(content.([2]string)[0], width/3)).
					SetAlign(tview.AlignLeft).
					SetStyle(tcell.StyleDefault.
						Foreground(tcell.ColorYellow).
						Background(tcell.ColorBlack)))
				inputTable.SetCell(i, 1, tview.NewTableCell(utils.GetFormattedString(content.([2]string)[1], width/3)).
					SetAlign(tview.AlignLeft).
					SetStyle(tcell.StyleDefault.
						Foreground(tcell.ColorPurple).
						Background(tcell.ColorBlack)))
			}
		case string:
			{
				b := content.(string)
				if !strings.HasPrefix(b, WHITE_AND_BOLD) {
					inputTable.SetCell(i, 0, tview.NewTableCell(content.(string)).
						SetAlign(tview.AlignLeft).
						SetStyle(tcell.StyleDefault.
							Foreground(tcell.ColorPurple).
							Background(tcell.ColorBlack)))
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
					tview.NewTableCell(utils.GetFormattedString(_songAttributes[0]["Title"], w/3)).
						SetAlign(tview.AlignLeft).
						SetStyle(tcell.StyleDefault.
							Foreground(tcell.ColorGreen).
							Background(tcell.ColorBlack)))

				inputTable.SetCell(i, 1,
					tview.NewTableCell(utils.GetFormattedString(_songAttributes[0]["Artist"], w/3)).
						SetAlign(tview.AlignLeft).
						SetStyle(tcell.StyleDefault.
							Foreground(tcell.ColorPurple).
							Background(tcell.ColorBlack)))

				inputTable.SetCell(i, 2,
					tview.NewTableCell(_songAttributes[0]["Album"]).
						SetAlign(tview.AlignLeft).
						SetStyle(tcell.StyleDefault.
							Foreground(tcell.ColorYellow).
							Background(tcell.ColorBlack)))

			} else if _songAttributes[0]["Title"] == "" {
				inputTable.SetCell(i, 0,
					tview.NewTableCell(j.Path).
						SetAlign(tview.AlignLeft).
						SetStyle(tcell.StyleDefault.
							Foreground(tcell.ColorBlue).
							Background(tcell.ColorBlack).
							Bold(true)))
			}
		} else {
			inputTable.SetCell(i, 0,
				tview.NewTableCell(j.Path).
					SetAlign(tview.AlignLeft).
					SetStyle(tcell.StyleDefault.
						Foreground(tcell.ColorYellow).
						Background(tcell.ColorBlack).
						Bold(true)))
		}
	}
}
