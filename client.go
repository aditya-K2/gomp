package main

import (
	"github.com/fhs/gompd/mpd"
	// "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	// "fmt"
)

var directoryMap map[string][]int = make(map[string][]int)

func togglePlayBack(connection mpd.Client) error {
	status, err := connection.Status()
	if status["state"] == "play" && err == nil {
		connection.Pause(true)
	} else if status["state"] == "pause" && err == nil {
		connection.Play(-1)
	}
	return err
}

func UpdatePlaylist(conn mpd.Client, t *tview.Table) {
	_playlistAttr, _ := conn.PlaylistInfo(-1, -1)

	t.Clear()
	for i, j := range _playlistAttr {
		if j["Title"] == "" || j["Artist"] == "" || j["Album"] == "" {
			t.SetCell(i, 0, tview.NewTableCell(j["file"]))
		} else {
			t.SetCell(i, 0, tview.NewTableCell(j["Title"]))
			t.SetCell(i, 1, tview.NewTableCell(j["Artist"]))
			t.SetCell(i, 2, tview.NewTableCell(j["Album"]))
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

func Update(conn mpd.Client, f []FileNode, inputTable *tview.Table) {
	inputTable.Clear()
	for i, j := range f {
		if len(j.children) == 0 {
			_songAttributes, err := conn.ListAllInfo(j.absolutePath)
			if err == nil && _songAttributes[0]["Title"] != "" {
				inputTable.SetCell(i, 0,
					tview.NewTableCell("[#fbff00]"+_songAttributes[0]["Title"]).
						SetAlign(tview.AlignLeft))

				inputTable.SetCell(i, 1,
					tview.NewTableCell("[#fbff00]"+_songAttributes[0]["Artist"]).
						SetAlign(tview.AlignLeft))

				inputTable.SetCell(i, 2,
					tview.NewTableCell("[#ff0030]"+_songAttributes[0]["Album"]).
						SetAlign(tview.AlignLeft))

			} else if _songAttributes[0]["Title"] == "" {
				inputTable.SetCell(i, 0,
					tview.NewTableCell("[blue]"+j.path).
						SetAlign(tview.AlignLeft))
			}
		} else {
			inputTable.SetCell(i, 0,
				tview.NewTableCell("[#fbff00::b]"+j.path).
					SetAlign(tview.AlignLeft))
		}
	}
}
