package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/fhs/gompd/mpd"
	"strconv"
)

// The progressBar is just a string which is separated by the color formatting String
// for e.g

// "[:#fbff00:]******************`innerText`[-:-:-]                "

// the above string shows represents the progress until [-:-:-]
// [-:-:-] this string represents resetting colors so the substring before it would be with a
// colored background. this is done by calculating the innerRect of the table and taking that length as
// 100% and then representing the rest of the information in relation to it

type progressBar struct{
	t *tview.Table
}

// This Function returns a progressBar with a table of two rows
// the First row will contain information about the current Song
// and the Second one will contain the progressBar
func newProgressBar(conn mpd.Client) *progressBar {
	p := progressBar{}

	a := tview.NewTable().
			SetCell(0, 0, tview.NewTableCell("")).
			SetCell(1, 0, tview.NewTableCell(""))

	a.SetBorder(true)

	a.SetDrawFunc(func( s tcell.Screen, x, y, width, height int) (int,int,int,int) {
		p.updateTitle(conn)
		p.updateProgress(conn)
		return p.t.GetInnerRect()
	})

	p = progressBar{a}

	return &p
}

func (s *progressBar) updateTitle(conn mpd.Client){
	_currentAttributes, err := conn.CurrentSong()
	if err == nil {
		s.t.GetCell(0,0).Text = _currentAttributes["Artist"] + " - " +_currentAttributes["Title"]
	}
}

func (s *progressBar) updateProgress(conn mpd.Client){
	_status,err := conn.Status()
	_, _, _width, _ := s.t.GetInnerRect()
	el, err1 := strconv.ParseFloat(_status["elapsed"], 8)
	du , err := strconv.ParseFloat(_status["duration"], 8)
	percentage := el / du * 100
	if err == nil && err1 == nil{
		s.t.GetCell(1,0).Text = getText(float64(_width), percentage, convertToStrings(el) + "/" + convertToStrings(du) + "(" + strconv.FormatFloat(percentage, 'f', 2, 32) + "%" + ")")
	}
}
