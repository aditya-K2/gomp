package render

import (
	"image"
	"os"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/notify"
	"github.com/aditya-K2/gomp/ui"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/utils"
	"github.com/nfnt/resize"
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
func (r *Renderer) Send(path string) {
	r.c <- path
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
		img2, _ := GetImg(extractedImage)
		im, _ = ueberzug.NewImage(img2,
			int(float32(ui.ImgX)*fw)+viper.GetInt("ADDITIONAL_PADDING_X"),
			int(float32(ui.ImgY)*fh)+viper.GetInt("ADDITIONAL_PADDING_Y"))
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

// Start Initialises the Renderer and calls the go routine OpenImage and passes the channel
// as argument.
func (r *Renderer) Start(path string) {
	go OpenImage(path, r.c)
}

// GetImagePath This Function returns the path to the image that is to be
// rendered it checks first for the image in the cache
// else it adds the image to the cache and then extracts it and renders it.
func GetImagePath(path string) string {
	a, err := client.Conn.ListInfo(path)
	var extractedImage string
	if err == nil && len(a) != 0 {
		if cache.Exists(a[0]["artist"], a[0]["album"]) {
			extractedImage = cache.GenerateName(a[0]["artist"], a[0]["album"])
		} else {
			imagePath := cache.GenerateName(a[0]["artist"], a[0]["album"])
			absPath := utils.CheckDirectoryFmt(viper.GetString("MUSIC_DIRECTORY")) + path
			extractedImage = ExtractImageFromFile(absPath, imagePath)
			if extractedImage == viper.GetString("DEFAULT_IMAGE_PATH") &&
				viper.GetString("GET_COVER_ART_FROM_LAST_FM") == "TRUE" {
				downloadedImage, err := getImageFromLastFM(a[0]["artist"], a[0]["album"], imagePath)
				if err == nil {
					notify.Notify.Send("Image From LastFM")
					extractedImage = downloadedImage
				} else {
					notify.Notify.Send("Falling Back to Default Image.")
				}
			} else {
				notify.Notify.Send("Extracted Image Successfully")
			}
		}
	}
	return extractedImage
}

func GetImg(uri string) (image.Image, error) {
	f, err := os.Open(uri)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	fw, fh := utils.GetFontWidth()
	img = resize.Resize(
		uint(float32(ui.ImgW)*(fw+float32(viper.GetFloat64("IMAGE_WIDTH_EXTRA_X")))),
		uint(float32(ui.ImgH)*(fh+float32(viper.GetFloat64("IMAGE_WIDTH_EXTRA_Y")))),
		img,
		resize.Bilinear,
	)
	return img, nil
}
