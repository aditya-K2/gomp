package notify

import (
	"time"

	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/utils"

	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

var (
	Notify *NotificationServer
)

/* Notification Primitive */
type Notification struct {
	*tview.Box
	Text string
}

/* Get A Pointer to A Notification Struct */
func NewNotification(s string) *Notification {
	return &Notification{
		Box:  tview.NewBox(),
		Text: s,
	}
}

/* Draw Function for the Notification Primitive */
func (self *Notification) Draw(screen tcell.Screen) {
	termDetails := utils.GetWidth()

	var (
		COL          int = int(termDetails.Col)
		TEXTLENGTH   int = len(self.Text)
		HEIGHT       int = 3
		TEXTPOSITION int = 2
	)

	self.Box.SetBackgroundColor(tcell.ColorBlack)
	self.SetRect(COL-(TEXTLENGTH+7), 1, TEXTLENGTH+4, HEIGHT)
	self.DrawForSubclass(screen, self.Box)
	tview.Print(screen, self.Text,
		COL-(TEXTLENGTH+5), TEXTPOSITION, TEXTLENGTH,
		tview.AlignCenter, tcell.ColorWhite)
}

/* Notification Server : Not an actual Server*/
type NotificationServer struct {
	c chan string
}

/* Get A Pointer to a NotificationServer Struct */
func NewNotificationServer() *NotificationServer {
	return &NotificationServer{
		c: make(chan string),
	}
}

/* This Method Starts the go routine for the NotificationServer */
func (self *NotificationServer) Start() {
	go NotificationRoutine(self.c, "EMPTY NOTIFICATION")
}

/* The Notification Server is just a string channel and the NotificationRoutine
the channel is used to receive the Notification Data through the Send Function
The Channel keeps listening for the Notification when it receives a Notification it checks if it
is Empty or not if it is an empty Notification it calls the NotificationRoutine with the empty routine else
it will call the go routine that renders the Notification for the Notification Interval and agains start listening
for the notfications sort of works like a que */
func NotificationRoutine(c chan string, s string) {
	if s != "EMPTY NOTIFICATION" {
		go func() {
			currentTime := time.Now().String()
			ui.Ui.Pages.AddPage(currentTime, NewNotification(s), false, true)
			ui.Ui.App.SetFocus(ui.Ui.ExpandedView)
			time.Sleep(time.Second * 1)
			ui.Ui.Pages.RemovePage(currentTime)
			ui.Ui.App.SetFocus(ui.Ui.ExpandedView)
		}()
	}
	NewNotification := <-c
	if NewNotification == "EMPTY NOTIFICATION" {
		NotificationRoutine(c, "EMPTY NOTIFICATION")
	} else {
		NotificationRoutine(c, NewNotification)
	}
}

/* Sends the Notification to the Notification Server */
func (self NotificationServer) Send(text string) {
	self.c <- text
}
