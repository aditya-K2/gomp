package main

import (
	"github.com/spf13/viper"
	"gitlab.com/diamondburned/ueberzug-go"
)

/*
	Renderer is just a channel on which we will send the Path to the song whose
	Image is to be Rendered. This channel is passed to the openImage which in turn is called
	by the Start() function as a go routine.
*/
type Renderer struct {
	c chan string
}

/*
	Returns a new Renderer with a string channel
*/
func newRenderer() *Renderer {
	c := make(chan string)
	return &Renderer{
		c: c,
	}
}

/*
	Send Image Path to Renderer
*/
func (self *Renderer) Send(path string) {
	self.c <- path
}

/*

   Go Routine that will Be Called and will listen on the channel c
   for changes and on getting a string over the channel will open the Image and
   keep listening again. This will keep the image blocked ( i.e no need to use time.Sleep() etc. )
   and saves resources too.

*/
func openImage(path string, c chan string) {
	fw, fh := getFontWidth()
	var im *ueberzug.Image
	if path != "stop" {
		absPath := viper.GetString("MUSIC_DIRECTORY") + path
		extractedImage := extractImageFromFile(absPath)
		if extractedImage == viper.GetString("DEFAULT_IMAGE_PATH") && viper.GetString("GET_COVER_ART_FROM_LAST_FM") == "TRUE" {
			a, err := CONN.ListInfo(path)
			if err == nil && len(a) != 0 {
				downloadedImage, err := getImageFromLastFM(a[0]["artist"], a[0]["album"])
				if err == nil {
					extractedImage = downloadedImage
				}
			}
		}
		img2, _ := getImg(extractedImage)
		im, _ = ueberzug.NewImage(img2, int(float32(IMG_X)*fw)+viper.GetInt("ADDITIONAL_PADDING_X"), int(float32(IMG_Y)*fh)+viper.GetInt("ADDITIONAL_PADDING_Y"))
	}
	d := <-c
	if im != nil {
		im.Clear()
	}
	if d != "stop" {
		openImage(d, c)
	} else {
		openImage("stop", c)
	}
}

/*
	Initialises the Renderer and calls the go routine openImage and passes the channel
	as argument.
*/
func (self *Renderer) Start(path string) {
	go openImage(path, self.c)
}
