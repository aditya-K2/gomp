package utils

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func GetWidth() *winsize {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return ws
}

func GetFontWidth() (float32, float32) {
	g := GetWidth()
	fw := (float32(g.Xpixel) / float32(g.Col))
	fh := (float32(g.Ypixel) / float32(g.Row))
	return fw, fh
}

func StrTime(e float64) string {
	a := int(e)
	var min, seconds string
	if a/60 < 10 {
		min = "0"
		min += strconv.Itoa(a / 60)
	} else {
		min = strconv.Itoa(a / 60)
	}
	if a%60 < 10 {
		seconds = "0"
		seconds += strconv.Itoa(a % 60)
	} else {
		seconds = strconv.Itoa(a % 60)
	}
	return min + ":" + seconds
}

func InsertAt(inputString, stringTobeInserted string, index int) string {
	s := inputString[:index] + stringTobeInserted + inputString[index:]
	return s
}

func GetText(width, percentage float64, eta string) string {
	q := "[#000000:#ffffff:b]"
	var a string
	a += strings.Repeat(" ", int(width)-len(eta))
	a = InsertAt(a, eta, int(width/2)-10)
	a = InsertAt(a, "[-:-:-]", int(width*percentage/100))
	q += a
	return q
}

func ConvertToArray(ArtistTree map[string]map[string]map[string]string) []string {
	var p []string
	for k2, v := range ArtistTree {
		p = append(p, k2)
		for k1, v1 := range v {
			p = append(p, k1)
			for k := range v1 {
				p = append(p, k)
			}
		}
	}
	return p
}

func FormatString(a interface{}) string {
	if a == "play" {
		return "Playing"
	} else if a == "1" {
		return "On"
	} else if a == "0" {
		return "Off"
	} else if a == "stop" {
		return "Stopped"
	} else {
		return "Paused"
	}
}

func Copy(sourceImage, destinationImage string) error {
	source, err := ioutil.ReadFile(sourceImage)
	if err != nil {
		return err
	} else {
		err = ioutil.WriteFile(destinationImage, source, 0644)
		if err != nil {
			return err
		}
		return nil
	}
}

func Join(stringSlice []string) string {
	var _s string = stringSlice[0]
	for i := 1; i < len(stringSlice); i++ {
		if _s != "" {
			_s += ("/" + stringSlice[i])
		}
	}
	return _s
}

func GetFormattedString(s string, width int) string {
	if len(s) < width {
		s += strings.Repeat(" ", (width - len(s)))
	} else {
		s = s[:(width - 2)]
		s += "  "
	}
	return s
}

func CheckDirectoryFmt(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	} else {
		return path + "/"
	}
}

func GetMatchedString(s string, color string, nulcol string, matchedIndexes []int) string {
	// The indexes that we will receive from the matchedIndexes are always sorted so we have to just
	// add the color string at
	//           `indexValue + ( len(colorString) * k )`
	//           where k is the index of the indexValue in the matchedIndexes slice
	// and we will need to also reset the colors, For that we check if the next indexValue in the matchedIndexes for
	// the current indexValue is not the consecutive value ( v + 1 ) if yes ( is not consecutive ) then we add the reset
	// color string at the k + 1 index in the string.
	// for e.g.
	//    if we have the following matchedIndexes slice
	//                        []int{ 1, 3, 4, 6}
	// During the First Iteration matchedIndexes[k] = 1 and and matchedIndexes[k+1] are not consecutive so the nulcol
	// string will be added to the matchedIndexes[k] + 1 index of the string
	// During the Second Iteration as 3, 4 are consecutive the nulcol will be skipped.
	color = fmt.Sprintf("[%s:-:bi]", color)
	nulcol = fmt.Sprintf("[%s:-:b]", nulcol)
	nulc := 0
	for k := range matchedIndexes {
		s = InsertAt(s, color, matchedIndexes[k]+(len(color)*k)+nulc)
		if k < len(matchedIndexes)-1 && matchedIndexes[k]-matchedIndexes[k+1] != 1 {
			s = InsertAt(s, nulcol, (matchedIndexes[k]+1)+(len(color)*(k+1))+nulc)
			nulc += len(nulcol)
		}
		if k == len(matchedIndexes)-1 {
			s = InsertAt(s, nulcol, ((matchedIndexes[len(matchedIndexes)-1] + 1) +
				(len(matchedIndexes) * len(color)) +
				(len(nulcol) * (len(matchedIndexes) - 1))))
		}
	}
	// Adding the Nulcol at the Start
	s = nulcol + s
	return s
}
