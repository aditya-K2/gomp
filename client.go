package main

import (
	"strings"

	"github.com/rivo/tview"
)

func getFormattedString(s string, width int) string {
	if len(s) < width {
		s += strings.Repeat(" ", (width - len(s)))
	} else {
		s = s[:(width - 2)]
		s += "  "
	}
	return s
}

func togglePlayBack() error {
	status, err := CONN.Status()
	if status["state"] == "play" && err == nil {
		CONN.Pause(true)
	} else if status["state"] == "pause" && err == nil {
		CONN.Play(-1)
	}
	return err
}

func UpdatePlaylist(inputTable *tview.Table) {
	_playlistAttr, _ := CONN.PlaylistInfo(-1, -1)

	inputTable.Clear()
	for i, j := range _playlistAttr {
		_, _, w, _ := inputTable.GetInnerRect()
		if j["Title"] == "" || j["Artist"] == "" || j["Album"] == "" {
			inputTable.SetCell(i, 0, tview.NewTableCell(getFormattedString(j["file"], w/3)))
		} else {
			inputTable.SetCell(i, 0, tview.NewTableCell(getFormattedString("[green]"+j["Title"], w/3)))
			inputTable.SetCell(i, 1, tview.NewTableCell(getFormattedString("[magenta]"+j["Artist"], w/3)))
			inputTable.SetCell(i, 2, tview.NewTableCell("[yellow]"+j["Album"]))
		}
	}
}

func join(stringSlice []string) string {
	var _s string = stringSlice[0]
	for i := 1; i < len(stringSlice); i++ {
		if _s != "" {
			_s += ("/" + stringSlice[i])
		}
	}
	return _s
}

func Update(f []FileNode, inputTable *tview.Table) {
	inputTable.Clear()
	for i, j := range f {
		if len(j.children) == 0 {
			_songAttributes, err := CONN.ListAllInfo(j.absolutePath)
			if err == nil && _songAttributes[0]["Title"] != "" {
				_, _, w, _ := inputTable.GetInnerRect()
				inputTable.SetCell(i, 0,
					tview.NewTableCell("[green]"+getFormattedString(_songAttributes[0]["Title"], w/3)).
						SetAlign(tview.AlignLeft))

				inputTable.SetCell(i, 1,
					tview.NewTableCell("[magenta]"+getFormattedString(_songAttributes[0]["Artist"], w/3)).
						SetAlign(tview.AlignLeft))

				inputTable.SetCell(i, 2,
					tview.NewTableCell("[yellow]"+_songAttributes[0]["Album"]).
						SetAlign(tview.AlignLeft))

			} else if _songAttributes[0]["Title"] == "" {
				inputTable.SetCell(i, 0,
					tview.NewTableCell("[blue]"+j.path).
						SetAlign(tview.AlignLeft))
			}
		} else {
			inputTable.SetCell(i, 0,
				tview.NewTableCell("[yellow::b]"+j.path).
					SetAlign(tview.AlignLeft))
		}
	}
}
