package ui

import (
	"sync"
	"time"

	"github.com/aditya-K2/utils"

	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

var (
	maxNotifications                = 3
	EmptyNotification *notification = newNotifcation(
		"gomp.notify.Notifcation.Empty")
	pm           sync.Mutex
	notAvailable = -1
	posArr       = positionArray{}
	c            chan string
)

// Start Notification Service
func InitNotifier() {
	for _m := maxNotifications; _m != 0; _m-- {
		posArr = append(posArr, true)
	}
	c = make(chan string, maxNotifications)
	routine()
}

// notification Primitive
type notification struct {
	*tview.Box
	Text     string
	TimeRecv time.Time
	Position int
}

// Array for all available positions where the notification can be displayed.
type positionArray []bool

// Check If there is a position available.
func (p *positionArray) Available() bool {
	var t = false
	pm.Lock()
	for _, v := range *p {
		t = t || v
	}
	pm.Unlock()
	return t
}

func (p *positionArray) GetNextPosition() int {
	pm.Lock()
	v := *p
	for k := range v {
		if v[k] {
			v[k] = false
			pm.Unlock()
			return k
		}
	}
	pm.Unlock()
	return notAvailable
}

// Free a position
func (p *positionArray) Free(i int) {
	pm.Lock()
	v := *p
	v[i] = true
	pm.Unlock()
}

// Get A Pointer to A Notification Struct
func newNotifcation(s string) *notification {
	return &notification{
		Box:      tview.NewBox(),
		Text:     s,
		TimeRecv: time.Now(),
	}
}

// Draw Function for the Notification Primitive
func (self *notification) Draw(screen tcell.Screen) {
	termDetails := utils.GetWidth()
	pos := (self.Position*3 + self.Position + 1)

	var (
		COL          int = int(termDetails.Col)
		TEXTLENGTH   int = len(self.Text)
		HEIGHT       int = 3
		TextPosition int = 1
	)

	self.Box.SetBackgroundColor(tcell.ColorBlack)
	self.SetRect(COL-(TEXTLENGTH+7), pos, TEXTLENGTH+4, HEIGHT)
	self.DrawForSubclass(screen, self.Box)
	tview.Print(screen, self.Text,
		COL-(TEXTLENGTH+5), pos+TextPosition, TEXTLENGTH,
		tview.AlignCenter, tcell.ColorWhite)
}

// this routine checks for available position and sends notification if
// position is available.
func routine() {
	go func() {
		for {
			val := <-c
			// Wait until a new position isn't available
			for !posArr.Available() {
				continue
			}
			notify(newNotifcation(val))
		}
	}()
}

func notify(s *notification) {
	if s != EmptyNotification {
		go func() {
			currentTime := time.Now().String()
			np := posArr.GetNextPosition()
			// Ensure a position is available.
			if np == notAvailable {
				for !posArr.Available() {
				}
				np = posArr.GetNextPosition()
			}
			s.Position = np
			Ui.Pages.AddPage(currentTime, s, false, true)
			Ui.App.SetFocus(Ui.MainS)
			time.Sleep(time.Second * 1)
			Ui.Pages.RemovePage(currentTime)
			posArr.Free(np)
			Ui.App.SetFocus(Ui.MainS)
		}()
	}
}

func SendNotification(text string) {
	go func() {
		c <- text
	}()
}
