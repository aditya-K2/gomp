package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

var (
	COLORS map[string]string = map[string]string{
		"RESET":  "\033[0m",
		"RED":    "\033[31m",
		"GREEN":  "\033[32m",
		"YELLOW": "\033[33m",
		"BLUE":   "\033[34m",
		"PURPLE": "\033[35m",
		"CYAN":   "\033[36m",
		"GRAY":   "\033[37m",
		"WHITE":  "\033[97m",
	}
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
	fw := float32(g.Xpixel) / float32(g.Col)
	fh := float32(g.Ypixel) / float32(g.Row)
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
	if source, rerr := os.ReadFile(sourceImage); rerr != nil {
		return rerr
	} else {
		if f, cerr := os.Create(destinationImage); cerr != nil {
			defer f.Close()
			return cerr
		} else {
			defer f.Close()
			if _, werr := f.Write(source); werr != nil {
				return werr
			}
		}
	}
	return nil
}

func GetFormattedString(s string, width int) string {
	if len(s) < width {
		s += strings.Repeat(" ", width-len(s))
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

func ExpandHomeDir(path string) string {
	HOME_DIR, _ := os.UserHomeDir()
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(HOME_DIR, path[1:])
	} else if path == "~" {
		return HOME_DIR
	} else {
		return path
	}
}

func GetMatchedString(a []int, s, color string) string {
	// The Matches are sorted so we just have to traverse the Matches and if the two adjacent matches are not consecutive
	// then we append the color string at the start + offset and the nulcol ( reset ) at end + offset + 1 and then reset
	// start and end to a[k+1] for e.g if matches := []int{1, 2, 4, 5, 6, 9} then the start will be 1 and end will be 1
	// now until we reach `4` the value of end will change to `2` that means when we reach `4` the s string will be
	// `O[yellow:-:-]ut[-:-:-]putString` after that until we reach the end will be changed and finally become `6` and the
	// s string will be `O[yellow:-:-]ut[-:-:-]p[yellow:-:-]utS[-:-:-]tring`
	// Please note that after around 45 simulatenously highlighted characters tview stops highlighting and the color
	// sequences are rendered hope no one has that big of search query.
	start := a[0]
	end := a[0]
	offset := 0
	nulcol := "[-:-:-]"
	for k := range a {
		if k < len(a)-1 && a[k+1]-a[k] == 1 {
			end = a[k+1]
		} else if k < len(a)-1 {
			s = InsertAt(s, color, start+offset)
			offset += len(color)
			s = InsertAt(s, nulcol, end+offset+1)
			offset += len(nulcol)
			start = a[k+1]
			end = a[k+1]
		} else if k == len(a)-1 {
			s = InsertAt(s, color, start+offset)
			offset += len(color)
			s = InsertAt(s, nulcol, end+offset+1)
			offset += len(nulcol)
		}
	}
	return s
}

func Unique(intSlice []int) []int {
	keys := make(map[int]bool)
	var list []int
	for _, entry := range intSlice {
		if _, exists := keys[entry]; !exists {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func GetNetwork(ntype, mport, naddress string) (string, string) {
	del := ""
	nt := ntype
	port := mport
	if nt == "tcp" {
		del = ":"
	} else if nt == "unix" && port != "" {
		port = ""
	}
	return nt, naddress + del + port
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func IsSame[K comparable](a []K, b []K) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if a[k] != b[k] {
			return false
		}
	}
	return true
}

func Pop[T comparable](index int, a []T) []T {
	return append(a[0:index], a[index+1:]...)
}

func Print(color, text string) {
	fmt.Print(COLORS[color] + text + COLORS["RESET"])
}

func DownloadImage(url string, imagePath string) (string, error) {
	var reader io.Reader
	if strings.HasPrefix(url, "http") {
		r, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer r.Body.Close()
		reader = r.Body
		v, err := io.ReadAll(reader)
		if err == nil {
			b, err := os.Create(imagePath)
			if err == nil {
				if _, err := b.Write(v); err == nil {
					return imagePath, nil
				} else {
					return "", errors.New("could Not Write Image")
				}
			} else {
				b.Close()
				return "", err
			}
		} else {
			return "", err
		}
	}
	return "", errors.New("image Not Received")
}
