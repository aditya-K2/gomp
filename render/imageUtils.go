// # Image Rendering
//
// ## Additional Padding and Extra Image Width
//
// The Default Position of the Image without any configuration assumes that you have no font or terminal padding or margin so Image will
// be rendered at different places in different terminals, Also the TUIs calculates positions with the respect to rows and columns
// and the image is rendered at pixel positions so the exact position can't be defined [ the app tries its best by calculating
// the font width and then predicting the position but it is best that you define some extra padding and your own image width ratio
// in config.yml. Please Read more about it in the [sample_config](https://github.com/aditya-K2/gomp/blob/master/sample_config.yml)
//
// for e.g
//
// ```yml
// # Default Values. Might be different than sample_config.yml
// ADDITIONAL_PADDING_X : 11
// ADDITIONAL_PADDING_Y : 18
//
// IMAGE_WIDTH_EXTRA_X  : -0.7
// IMAGE_WIDTH_EXTRA_Y  : -2.6
// ```
// ![Cover Art Position](./assets/default.png)
//
// Let's say upon opening gomp for the first time and your image is rendered this way.
//
// Here the `Y` Position is too low hence we have to decrease the `ADDITIONAL_PADDING_Y` so that image will be rendered
// in a better position so we decrement the  `ADDITIONAL_PADDING_Y` by `9`
//
// Now the configuration becomes
// ```yml
// ADDITIONAL_PADDING_Y : 9
// ```
//
// and the image appears like this:
//
// ![Padding](./assets/padding.png)
//
// One might be happy the way things turn out but for the perfectionist like me this is not enough.
// You can notice that the Height of the image is a little bit more than the box height and hence the image is flowing outside the box. Now it's  time to change the `WIDTH_Y`.
//
// Width can be changed by defining the `IMAGE_WIDTH_EXTRA_Y` and `IMAGE_WIDTH_EXTRA_X` it act's a little differently think of it like a chunk which is either added or subtracted from the image's original width. We can look at the configuration and realize that the chunk `IMAGE_WIDTH_EXTRA_Y` when subtracted from the original `IMAGE_WIDTH` doesn't get us the proper result and is a little to low. We need to subtract a more bigger chunk, Hence we will increase the magnitude of `IMAGE_WIDTH_EXTRA_Y` or decrease `IMAGE_WIDTH_EXTRA_Y`
//
// Now the Configuration becomes:
// ```yml
// IMAGE_WIDTH_EXTRA_Y : - 3.2
// ```
// and the image appears like this:
//
// ![Width](./assets/width.png)
//
// Which looks perfect. 🎉
//
// Read More about Additional Padding and Image Width in the [sample_config](https://github.com/aditya-K2/gomp/blob/master/sample_config.yml)
//
// Please change the configuration according to your needs.
package render

import (
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/ui/notify"
	"github.com/aditya-K2/gomp/utils"
	"github.com/dhowden/tag"
	"github.com/nfnt/resize"
)

var (
	ExtractionErr    = errors.New("No Image Extracted")
	ImageNotFoundErr = errors.New("Image Not Found")
	ImageWriteErr    = errors.New("Error Writing Image to File")
)

func ExtractImage(path string, imagePath string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", ExtractionErr
	}
	defer f.Close()

	meta, err := tag.ReadFrom(f)
	if err != nil {
		return "", ExtractionErr
	}

	if picture := meta.Picture(); picture == nil {
		return "", ImageNotFoundErr
	} else {
		if imHd, err := os.Create(imagePath); err != nil {
			return "", ImageWriteErr
		} else {
			imHd.Write(picture.Data)
		}
	}

	return imagePath, nil
}

/*
GetImagePath This Function returns the path to the image that is to be
rendered it checks first for the image in the cache
else it adds the image to the cache and then extracts it and renders it.
*/
func GetImagePath(path string) string {
	a, err := client.Conn.ListInfo(path)
	var extractedImage string = config.Config.DefaultImagePath
	if err == nil && len(a) != 0 {
		if cache.Exists(a[0]["artist"], a[0]["album"]) {
			extractedImage = cache.GenerateName(a[0]["artist"], a[0]["album"])
		} else {
			imagePath := cache.GenerateName(a[0]["artist"], a[0]["album"])
			absPath := utils.CheckDirectoryFmt(config.Config.MusicDirectory) + path
			if _eimg, exErr := ExtractImage(absPath, imagePath); exErr != nil {
				if config.Config.GetCoverArtFromLastFm {
					downloadedImage, lFmErr := getImageFromLastFM(a[0]["artist"], a[0]["album"], imagePath)
					if lFmErr == nil {
						notify.Send("Image From LastFM")
						extractedImage = downloadedImage
					} else {
						notify.Send(exErr.Error())
					}
				}
			} else {
				notify.Send("Image Extracted Succesfully!")
				extractedImage = _eimg
			}
		}
	} else {
		notify.Send(fmt.Sprintf("Couldn't Get Attributes for %s", path))
	}
	return extractedImage
}

/* Gets the Image Struct from the provided path */
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
		uint(float32(ui.ImgW)*(fw+float32(config.Config.ExtraImageWidthX))),
		uint(float32(ui.ImgH)*(fh+float32(config.Config.ExtraImageWidthY))),
		img,
		resize.Bilinear,
	)
	return img, nil
}
