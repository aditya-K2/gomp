package ui

import (
	"container/heap"
	"sync"
	"time"

	"github.com/aditya-K2/gomp/utils"

	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

var (
	maxNotifications                = 3
	EmptyNotification *notification = newNotifcation(
		"gomp.notify.Notifcation.Empty")
	nQueue       *notificationQueue
	qm, pm       sync.Mutex
	notAvailable = -1
	posArr       = positionArray{}
)

// Start Notification Service
func InitNotifier() {
	for _m := maxNotifications; _m != 0; _m-- {
		posArr = append(posArr, true)
	}
	nQueue = &notificationQueue{}
	heap.Init(nQueue)
	queueRoutine()
}

// notification Primitive
type notification struct {
	*tview.Box
	Text     string
	TimeRecv time.Time
	Position int
}

type notificationQueue []*notification

func (n notificationQueue) Len() int { return len(n) }
func (n notificationQueue) Less(i, j int) bool {
	ctime := time.Now()
	// Return the Oldest One.
	return ctime.Sub(n[i].TimeRecv) > ctime.Sub(n[j].TimeRecv)
}
func (n notificationQueue) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

func (n *notificationQueue) Push(x any) {
	*n = append(*n, x.(*notification))
}

func (n *notificationQueue) Pop() any {
	old := *n
	_n := len(old)
	x := old[_n-1]
	*n = old[0 : _n-1]
	return x
}

func (n notificationQueue) Empty() bool {
	return len(n) == 0
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

// this routine checks for available position and then sends
// the notification at the top of the queue to the notificationRoutine
func queueRoutine() {
	go func() {
		t := time.NewTicker(time.Millisecond * 200)
		for {
			select {
			case <-t.C:
				{
					for !posArr.Available() {
						continue
					}
					if !nQueue.Empty() {
						qm.Lock()
						_new := heap.Pop(nQueue).(*notification)
						qm.Unlock()
						notificationRoutine(_new)
					}
				}
			}
		}
	}()
}

func notificationRoutine(s *notification) {
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
		qm.Lock()
		heap.Push(nQueue, newNotifcation(text))
		qm.Unlock()
	}()
}
