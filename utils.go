package main
import(
	"strings"
	"strconv"
)

func strTime(e float64) string{
	a := int(e)
	var min, seconds string
	if (a/60 < 10){
		min = "0"
		min += strconv.Itoa(a/60)
	} else {
		min = strconv.Itoa(a/60)
	}
	if (a%60 < 10){
		seconds = "0"
		seconds += strconv.Itoa(a%60)
	} else{
		seconds = strconv.Itoa(a%60)
	}
	return min + ":" + seconds
}

func insertAt(inputString, stringTobeInserted string, index int) string{
    s := inputString[:index] + stringTobeInserted + inputString[index:]
	return s
}

func getText(width , percentage float64, eta string) string{
	q := "[#000000:#ffffff:bl]"
	var a string
	a += strings.Repeat(" ", int(width) - len(eta))
	a = insertAt(a, eta, int(width/2) - 10)
	a = insertAt(a, "[-:-:-]", int(width * percentage / 100))
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
