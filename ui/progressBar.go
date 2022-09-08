package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/gdamore/tcell/v2"

	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
)

var (
	CurrentSong string = ""
	CONN        *mpd.Client
	RENDERER    interface{ Send(string) }
)

func SetConnection(c *mpd.Client) {
	CONN = c
}

func ConnectRenderer(r interface{ Send(string) }) {
	RENDERER = r
}

// ProgressBar is a two-lined Box. First line is the BarTitle
// Second being the actual progress done.
// Use SetProgressFunc to provide the callback which provides the Fields each time the ProgressBar will be Drawn.
// The progressFunc must return (BarTitle, BarTopTitle, BarText, percentage) respectively
type ProgressBar struct {
	*tview.Box
	BarTitle     string
	BarText      string
	BarTopTitle  string
	progressFunc func() (BarTitle string,
		BarTopTitle string,
		BarText string,
		percentage float64)
}

func (self *ProgressBar) SetProgressFunc(pfunc func() (string, string, string, float64)) *ProgressBar {
	self.progressFunc = pfunc
	return self
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		Box: tview.NewBox(),
	}
}

func GetProgressGlyph(width, percentage float64, btext string) string {
	q := "[black:white:b]"
	var a string
	a += strings.Repeat(" ", int(width)-len(btext))
	a = utils.InsertAt(a, btext, int(width/2)-10)
	a = utils.InsertAt(a, "[-:-:-]", int(width*percentage/100))
	q += a
	return q
}

func (self *ProgressBar) Draw(screen tcell.Screen) {
	var (
		OFFSET int = 1
	)
	self.Box.SetBorder(true)
	self.Box.SetBackgroundColor(tcell.ColorDefault)
	var percentage float64
	self.BarTitle, self.BarTopTitle, self.BarText, percentage = self.progressFunc()
	self.DrawForSubclass(screen, self.Box)
	self.Box.SetTitle(self.BarTopTitle)
	self.Box.SetTitleAlign(tview.AlignRight)
	x, y, _width, _ := self.Box.GetInnerRect()
	tview.Print(screen, self.BarTitle, x+OFFSET, y, _width, tview.AlignLeft, tcell.ColorWhite)
	tview.Print(screen,
		GetProgressGlyph(float64(_width-OFFSET-1),
			percentage,
			self.BarText),
		x, y+2, _width-OFFSET, tview.AlignRight, tcell.ColorWhite)
}

func progressFunction() (string, string, string, float64) {
	_currentAttributes, err := CONN.CurrentSong()
	var song, top, text string
	var percentage float64
	if err == nil {
		song = "[green::bi]" +
			_currentAttributes["Title"] + "[-:-:-] - " + "[blue::b]" +
			_currentAttributes["Artist"] + "\n"
		if len(_currentAttributes) == 0 && CurrentSong != "" {
			CurrentSong = ""
			RENDERER.Send("stop")
		} else if song != CurrentSong && len(_currentAttributes) != 0 {
			RENDERER.Send(_currentAttributes["file"])
			CurrentSong = song
		}
	} else {
		utils.Print("RED", "Error Retrieving Current Song\n")
		panic(err)
	}
	_status, err := CONN.Status()
	el, err1 := strconv.ParseFloat(_status["elapsed"], 8)
	du, err := strconv.ParseFloat(_status["duration"], 8)
	top = fmt.Sprintf("[[::i] %s [-:-:-]Shuffle: %s Repeat: %s Volume: %s ]",
		utils.FormatString(_status["state"]),
		utils.FormatString(_status["random"]),
		utils.FormatString(_status["repeat"]),
		_status["volume"])
	if du != 0 {
		percentage = el / du * 100
		if err == nil && err1 == nil {
			text = utils.StrTime(el) + "/" + utils.StrTime(du) +
				"(" + strconv.FormatFloat(percentage, 'f', 2, 32) + "%" + ")"
		} else {
			text = ""
		}
	} else {
		text = "   ---:---"
		percentage = 0
	}
	if percentage > 100 {
		percentage = 0
	}
	return song, top, text, percentage
}
