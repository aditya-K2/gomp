package main

import (
	"github.com/fhs/gompd/mpd"
	// "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	// "fmt"
)

func togglePlayBack(connection mpd.Client) error {
	status, err := connection.Status()
	if(status["state"] == "play" && err == nil){
		connection.Pause(true)
	} else if(status["state"] == "pause" && err == nil) {
		connection.Play(-1)
	}
	return err
}

func join ( stringSlice [] string ) string{
	var _s string = stringSlice[0]
	for i:= 1; i<len(stringSlice);i++{
		if ( _s != ""){
			_s += ( "/" + stringSlice[i])
		}
	}
	return _s
}

func Update(conn mpd.Client, f []FileNode, currentDirectoryStructure [] string, inputTable *tview.Table) {
	ab := join(currentDirectoryStructure)
	for i,j := range(f){
		if len(j.children) == 0 {
			a, err := conn.ListAllInfo(ab)
			if err == nil {

				inputTable.SetCell(i, 0,
				tview.NewTableCell("[#fbff00]" + a[0]["Title"]).
				SetAlign(tview.AlignLeft))

				inputTable.SetCell(i, 1,
				tview.NewTableCell("[#000000]" + a[0]["Artist"]).
				SetAlign(tview.AlignLeft))

				inputTable.SetCell(i, 2,
				tview.NewTableCell("[#ff0030]" + a[0]["Album"]).
				SetAlign(tview.AlignLeft))

			}
		} else {

			inputTable.SetCell(i, 0,
			tview.NewTableCell("[#fbff00::b]" + j.path).
			SetAlign(tview.AlignLeft))

		}
	}
}

func addSong(conn mpd.Client, currentDirectoryStructure [] string, currentCellContent string){
	conn.Add(join(currentDirectoryStructure) + "/" + currentCellContent)
}
