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
	Notify            *NotificationServer
	MaxNotifications                = 3
	EmptyNotification *Notification = NewNotification(
		"gomp.notify.Notifcation.Empty")
	NQueue       *NotificationQueue
	qm, pm       sync.Mutex
	NotAvailable = -1
	posArr       = PositionArray{}
)

func Init() {
	for _m := MaxNotifications; _m != 0; _m-- {
		posArr = append(posArr, true)
	}
	Notify = NewNotificationServer()
	NQueue = &NotificationQueue{}
	heap.Init(NQueue)
	Notify.QueueRoutine()
	Notify.Start()
}

/* Notification Primitive */
type Notification struct {
	*tview.Box
	Text     string
	TimeRecv time.Time
	Position int
}

type NotificationQueue []*Notification

func (n NotificationQueue) Len() int { return len(n) }
func (n NotificationQueue) Less(i, j int) bool {
	ctime := time.Now()
	return ctime.Sub(n[i].TimeRecv) < ctime.Sub(n[j].TimeRecv)
}
func (n NotificationQueue) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

func (n *NotificationQueue) Push(x any) {
	*n = append(*n, x.(*Notification))
}

func (n *NotificationQueue) Pop() any {
	old := *n
	_n := len(old)
	x := old[_n-1]
	*n = old[0 : _n-1]
	return x
}

func (n NotificationQueue) Empty() bool {
	return len(n) == 0
}

type PositionArray []bool

func (p *PositionArray) Available() bool {
	var t = false
	pm.Lock()
	for _, v := range *p {
		t = t || v
	}
	pm.Unlock()
	return t
}

func (p *PositionArray) GetNextPosition() int {
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
	return NotAvailable
}

func (p *PositionArray) Free(i int) {
	pm.Lock()
	v := *p
	v[i] = true
	pm.Unlock()
}

/* Get A Pointer to A Notification Struct */
func NewNotification(s string) *Notification {
	return &Notification{
		Box:      tview.NewBox(),
		Text:     s,
		TimeRecv: time.Now(),
	}
}

/* Draw Function for the Notification Primitive */
func (self *Notification) Draw(screen tcell.Screen) {
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

/* Notification Server : Not an actual Server*/
type NotificationServer struct {
	c chan *Notification
}

/* Get A Pointer to a NotificationServer Struct */
func NewNotificationServer() *NotificationServer {
	return &NotificationServer{
		c: make(chan *Notification, MaxNotifications),
	}
}

/* This Method Starts the go routine for the NotificationServer */
func (self *NotificationServer) Start() {
	go NotificationRoutine(self.c, EmptyNotification)
}

func (self NotificationServer) QueueRoutine() {
	go func() {
		for {
			if len(self.c) <= MaxNotifications {
				if !NQueue.Empty() {
					qm.Lock()
					n := heap.Pop(NQueue).(*Notification)
					self.c <- n
					qm.Unlock()
				}
			}
		}
	}()
}

/* The Notification Server is just a *Notification channel and the
NotificationRoutine the channel is used to receive the Notification Data
through the Send Function The Channel keeps listening for the Notification
when it receives a Notification it checks if it is Empty or not if it is an
empty Notification it calls the NotificationRoutine with the empty routine else
it will call the go routine that renders the Notification for the Notification
Interval and agains start listening for the notfications sort of works like
a queue*/
func NotificationRoutine(c chan *Notification, s *Notification) {
	if s != EmptyNotification {
		go func() {
			currentTime := time.Now().String()
			np := posArr.GetNextPosition()
			if np == NotAvailable {
				for !posArr.Available() {
				}
				np = posArr.GetNextPosition()
			}
			s.Position = np
			Ui.Pages.AddPage(currentTime, s, false, true)
			Ui.App.SetFocus(Ui.ExpandedView)
			time.Sleep(time.Second * 1)
			Ui.Pages.RemovePage(currentTime)
			posArr.Free(np)
			Ui.App.SetFocus(Ui.ExpandedView)
		}()
	}
	for !posArr.Available() {
		continue
	}
	_new := <-c
	NotificationRoutine(c, _new)
}

/* Sends the Notification to the Notification Server */
func (self NotificationServer) Send(text string) {
	go func() {
		qm.Lock()
		heap.Push(NQueue, NewNotification(text))
		qm.Unlock()
	}()
}
