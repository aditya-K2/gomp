// # Colors and Style
//
// You can change `colors` and `styles` for some of the aspects of `gomp`
//
// #### Let's say to you want to change Color of Artist from the default Purple to Red
//
// In your `config.yml`
// ```yml
// COLORS:
//
// artist:
//
//	foreground: Red
//
// # Another Example
// pbar_artist:
//
//	foreground: "#ff0000" # For Hex Values
//	bold: True # Changes the Style
//	italic: False
//
// ```
//
// ![Dec30(Fri)012241PM](https://user-images.githubusercontent.com/51816057/210048064-b2816095-10f2-4f0b-83ed-0e87d636b894.png)
// ![Dec30(Fri)012315PM](https://user-images.githubusercontent.com/51816057/210048069-8e91509a-17a5-46da-a65e-ff8f427dde17.png)
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/aditya-K2/utils"
	"github.com/gdamore/tcell/v2"
)

var (
	ColorError = func(s string) {
		_s := fmt.Sprintf("Wrong Color Provided: %s", s)
		utils.Print("RED", _s)
		os.Exit(-1)
	}
	DColors = map[string]tcell.Color{
		"Black":   tcell.ColorBlack,
		"Maroon":  tcell.ColorMaroon,
		"Green":   tcell.ColorGreen,
		"Olive":   tcell.ColorOlive,
		"Navy":    tcell.ColorNavy,
		"Purple":  tcell.ColorPurple,
		"Teal":    tcell.ColorTeal,
		"Silver":  tcell.ColorSilver,
		"Gray":    tcell.ColorGray,
		"Red":     tcell.ColorRed,
		"Lime":    tcell.ColorLime,
		"Yellow":  tcell.ColorYellow,
		"Blue":    tcell.ColorBlue,
		"Fuchsia": tcell.ColorFuchsia,
		"Aqua":    tcell.ColorAqua,
		"White":   tcell.ColorWhite,
	}
)

type Color struct {
	Fg     string `mapstructure:"foreground"`
	Bg     string `mapstructure:"background"`
	Bold   bool   `mapstructure:"bold"`
	Italic bool   `mapstructure:"italic"`
}

// ### Following Aspects can be changed:
type Colors struct {
	// - `album`
	Artist Color `mapstructure:"artist"`
	// - `artist`
	Album Color `mapstructure:"album"`
	// - `track`
	Track Color `mapstructure:"track"`
	// - `file`
	File Color `mapstructure:"file"`
	// - `folder`
	Folder Color `mapstructure:"folder"`
	// - `timestamp`
	Timestamp Color `mapstructure:"timestamp"`
	// - `matched_title`
	MatchedTitle Color `mapstructure:"matched_title"`
	// - `matched_folder`
	MatchedFolder Color `mapstructure:"matched_folder"`
	// - `pbar_artist`
	PBarArtist Color `mapstructure:"pbar_artist"`
	// - `pbar_track`
	PBarTrack Color `mapstructure:"pbar_track"`
	// - `autocomplete`
	Autocomplete Color `mapstructure:"autocomplete"`
	Null         Color
}

func (c Color) Foreground() tcell.Color {
	if strings.HasPrefix(c.Fg, "#") && len(c.Fg) == 7 {
		return tcell.GetColor(c.Fg)
	} else if val, ok := DColors[c.Fg]; ok {
		return val
	} else {
		ColorError(c.Fg)
		return tcell.ColorBlack
	}
}

func (c Color) Background() tcell.Color {
	if c.Bg == "" {
		return tcell.ColorBlack
	}
	if strings.HasPrefix(c.Bg, "#") && len(c.Bg) == 7 {
		return tcell.GetColor(c.Bg)
	} else if val, ok := DColors[c.Bg]; ok {
		return val
	} else {
		ColorError(c.Bg)
		return tcell.ColorBlack
	}
}

func (c Color) Style() tcell.Style {
	return tcell.StyleDefault.
		Foreground(c.Foreground()).
		Background(c.Background()).
		Bold(c.Bold).
		Italic(c.Italic)
}

func (c Color) String() string {
	style := ""
	if c.Bold {
		style += "b"
	}
	if c.Italic {
		style += "i"
	}
	checkColor := func(s string) string {
		var res string
		if _, ok := DColors[s]; ok {
			res = strings.ToLower(s)
		} else if strings.HasPrefix(s, "#") && len(s) == 7 {
			res = s
		} else {
			ColorError(s)
		}
		return res
	}
	foreground := checkColor(c.Fg)
	return fmt.Sprintf("[%s::%s]", foreground, style)
}

func NewColors() *Colors {
	return &Colors{
		Artist: Color{
			Fg:     "Purple",
			Bold:   false,
			Italic: false,
		},
		Album: Color{
			Fg:     "Yellow",
			Bold:   false,
			Italic: false,
		},
		Track: Color{
			Fg:     "Green",
			Bold:   false,
			Italic: false,
		},
		Timestamp: Color{
			Fg:     "Red",
			Bold:   false,
			Italic: true,
		},
		File: Color{
			Fg:     "Blue",
			Bold:   true,
			Italic: false,
		},
		Folder: Color{
			Fg:     "Yellow",
			Bold:   true,
			Italic: false,
		},
		MatchedFolder: Color{
			Fg:     "Blue",
			Bold:   true,
			Italic: true,
		},
		MatchedTitle: Color{
			Fg:     "Yellow",
			Bold:   true,
			Italic: true,
		},
		PBarArtist: Color{
			Fg:     "Blue",
			Bold:   true,
			Italic: false,
		},
		PBarTrack: Color{
			Fg:     "Green",
			Bold:   true,
			Italic: true,
		},
		Autocomplete: Color{
			Fg:     "White",
			Bg:     "Black",
			Bold:   false,
			Italic: false,
		},
		Null: Color{
			Fg:     "White",
			Bold:   true,
			Italic: false,
		},
	}
}
