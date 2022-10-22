package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
)

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
