package ui

import (
	"errors"
	"fmt"

	"github.com/aditya-K2/gomp/config"
	"github.com/fhs/gompd/v2/mpd"

	"gitlab.com/diamondburned/ueberzug-go"
)

var (
	Rendr *Renderer
)

func getFontWidth() (int, int, error) {
	w, h, err := ueberzug.GetParentSize()
	if err != nil {
		return 0, 0, err
	}
	_, _, rw, rh := Ui.Pages.GetRect()
	if rw == 0 || rh == 0 {
		return 0, 0, errors.New("Unable to get row width and height")
	}
	fw := w / rw
	fh := h / rh
	return fw, fh, nil
}

// Renderer is just a channel on which we will send the Path to the song whose
// Image is to be Rendered. This channel is passed to the OpenImage which in turn is called
// by the Start() function as a go routine.
type Renderer struct {
	c chan string
}

// NewRenderer Returns a new Renderer with a string channel
func NewRenderer() *Renderer {
	c := make(chan string)
	return &Renderer{
		c: c,
	}
}

// Send Image Path to Renderer
// Start Initialises the Renderer and calls the go routine OpenImage and passes the channel
// as argument.
func (r *Renderer) Send(path string, start bool) {
	if start {
		go OpenImage(path, r.c)
	} else {
		r.c <- path
	}
}

// OpenImage Go Routine that will Be Called and will listen on the channel c
// for changes and on getting a string over the channel will open the Image and
// keep listening again. This will keep the image blocked ( i.e. no need to use time.Sleep() etc. )
// and saves resources too.
func OpenImage(path string, c chan string) {
	fw, fh, err := getFontWidth()
	if err != nil {
		Ui.App.Stop()
		fmt.Printf("Error Occured While getting font width: %v\n", err)
	}
	var im *ueberzug.Image
	if path != "stop" {
		extractedImage := GetImagePath(path)
		if img2, err := GetImg(extractedImage); err == nil {
			im, _ = ueberzug.NewImage(img2,
				int(ImgX*fw)+config.Config.AdditionalPaddingX,
				int(ImgY*fh)+config.Config.AdditionalPaddingY)
		} else {
			SendNotification("Error Rendering Image!")
		}
	}
	d := <-c
	if im != nil {
		im.Clear()
	}
	if d != "stop" {
		OpenImage(d, c)
	} else {
		OpenImage("stop", c)
	}
}

func DrawCover(c mpd.Attrs, start bool) {
	if len(c) == 0 {
		Rendr.Send("stop", start)
	} else {
		Rendr.Send(c["file"], start)
	}
}
