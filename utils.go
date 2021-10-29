package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

func getWidth() *winsize {
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

func getFontWidth() (float32, float32) {
	g := getWidth()
	fw := (float32(g.Xpixel) / float32(g.Col))
	fh := (float32(g.Ypixel) / float32(g.Row))
	return fw, fh
}

func strTime(e float64) string {
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

func insertAt(inputString, stringTobeInserted string, index int) string {
	s := inputString[:index] + stringTobeInserted + inputString[index:]
	return s
}

func getText(width, percentage float64, eta string) string {
	q := "[#000000:#ffffff:b]"
	var a string
	a += strings.Repeat(" ", int(width)-len(eta))
	a = insertAt(a, eta, int(width/2)-10)
	a = insertAt(a, "[-:-:-]", int(width*percentage/100))
	q += a
	return q
}

func formatString(a interface{}) string {
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

func getMusicDirectory() string {
	a, _ := os.UserHomeDir()
	content, err := ioutil.ReadFile(a + "/.config/mpd/mpd.conf")
	if err != nil {
		fmt.Println("No Config File for mpd Found")
		panic(err)
	}
	ab := string(content)
	maps := strings.Split(ab, "\n")
	for _, j := range maps {
		if strings.Contains(j, "music_directory") {
			s := strings.SplitAfter(strings.ReplaceAll(j, " ", ""), "y")[1]
			d := ""
			for z, m := range s {
				if (z != 0) && (z != (len(s) - 1)) {
					d += string(m)
				}
			}
			return d
		}
	}
	return ""
}
