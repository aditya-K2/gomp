package render

import (
	"github.com/aditya-K2/gomp/ui"
	"github.com/fhs/gompd/v2/mpd"

	"github.com/aditya-K2/gomp/utils"
	"github.com/spf13/viper"
	"gitlab.com/diamondburned/ueberzug-go"
)

var (
	Rendr *Renderer
)

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
	fw, fh := utils.GetFontWidth()
	var im *ueberzug.Image
	if path != "stop" {
		extractedImage := GetImagePath(path)
		if img2, err := GetImg(extractedImage); err == nil {
			im, _ = ueberzug.NewImage(img2,
				int(float32(ui.ImgX)*fw)+viper.GetInt("ADDITIONAL_PADDING_X"),
				int(float32(ui.ImgY)*fh)+viper.GetInt("ADDITIONAL_PADDING_Y"))
		} else {
			ui.Notify.Send("Error Rendering Image!")
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
