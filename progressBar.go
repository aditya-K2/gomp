package main

import (
	"fmt"
	"strconv"

	"github.com/aditya-K2/goMP/utils"

	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

var CurrentSong string

// The progressBar is just a string which is separated by the color formatting String
// for e.g
// "[:#fbff00:]******************`innerText`[-:-:-]                "
// the above string shows represents the progress until [-:-:-]
// [-:-:-] this string represents resetting colors so the substring before it would be with a
// colored background. this is done by calculating the innerRect of the table and taking that length as
// 100% and then representing the rest of the information in relation to it
type progressBar struct {
	t *tview.Table
}

// This Function returns a progressBar with a table of two rows
// the First row will contain information about the current Song
// and the Second one will contain the progressBar
func newProgressBar(r *Renderer) *progressBar {
	p := progressBar{}

	a := tview.NewTable().
		SetCell(0, 0, tview.NewTableCell("")).
		SetCell(1, 0, tview.NewTableCell("")).
		SetCell(2, 0, tview.NewTableCell(""))

	a.SetBorder(true)

	a.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		p.updateTitle(r)
		p.updateProgress()
		return p.t.GetInnerRect()
	})

	CurrentSong = ""

	p = progressBar{a}
	return &p
}

func (s *progressBar) updateTitle(r *Renderer) {
	_currentAttributes, err := CONN.CurrentSong()
	if err == nil {
		song := "[green::bi]" + _currentAttributes["Title"] + "[-:-:-] - " + "[blue::b]" + _currentAttributes["Artist"] + "\n"
		s.t.GetCell(0, 0).Text = song
		if len(_currentAttributes) == 0 && CurrentSong != "" {
			CurrentSong = ""
			r.Send("stop")
		} else if song != CurrentSong && len(_currentAttributes) != 0 {
			r.Send(_currentAttributes["file"])
			CurrentSong = song
		}
	}
}

func (s *progressBar) updateProgress() {
	_status, err := CONN.Status()
	_, _, _width, _ := s.t.GetInnerRect()
	el, err1 := strconv.ParseFloat(_status["elapsed"], 8)
	du, err := strconv.ParseFloat(_status["duration"], 8)
	if du != 0 {
		percentage := el / du * 100
		if err == nil && err1 == nil {
			s.t.SetTitle(fmt.Sprintf("[[::i] %s [-:-:-]Shuffle: %s Repeat: %s Volume: %s ]", utils.FormatString(_status["state"]), utils.FormatString(_status["random"]), utils.FormatString(_status["repeat"]), _status["volume"])).SetTitleAlign(tview.AlignRight)
			s.t.GetCell(2, 0).Text = utils.GetText(float64(_width), percentage, utils.StrTime(el)+"/"+utils.StrTime(du)+"("+strconv.FormatFloat(percentage, 'f', 2, 32)+"%"+")")
		} else {
			s.t.SetTitle(fmt.Sprintf("[[::i] %s [-:-:-]Shuffle: %s Repeat: %s]", utils.FormatString(_status["state"]), utils.FormatString(_status["random"]), utils.FormatString(_status["repeat"]))).SetTitleAlign(tview.AlignRight)
			s.t.GetCell(2, 0).Text = ""
		}
	} else {
		s.t.SetTitle(fmt.Sprintf("[[::i] %s [-:-:-]Shuffle: %s Repeat: %s Volume: %s ]", utils.FormatString(_status["state"]), utils.FormatString(_status["random"]), utils.FormatString(_status["repeat"]), _status["volume"])).SetTitleAlign(tview.AlignRight)
		s.t.GetCell(2, 0).Text = utils.GetText(float64(_width), 0, "   ---:---")
	}
}
