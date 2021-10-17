// Demo code for the TextView primitive.
package main

import (
	"log"
	"strconv"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var Volume int64
var Random bool
var Repeat bool

func main() {

	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	var pBar *progressBar = newProgressBar(*conn)
	expandedView := tview.NewTable()
	Navbar := tview.NewTable()
	searchBar := tview.NewTable()

	searchBar.SetBorder(true)
	Navbar.SetBorder(true)
	Navbar.SetSelectable(true, false)
	Navbar.SetCell(0, 0, tview.NewTableCell("Files"))
	Navbar.SetCell(1, 0, tview.NewTableCell("Playlist"))
	Navbar.SetCell(2, 0, tview.NewTableCell("Most Played"))

	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBar, 0, 1, false).
		AddItem(Navbar, 0, 7, false)

	sNavExpViewFlex := tview.NewFlex().
		AddItem(searchNavFlex, 0, 1, false).
		AddItem(expandedView, 0, 4, false)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sNavExpViewFlex, 0, 8, false).
		AddItem(pBar.t, 0, 1, false)

	App := tview.NewApplication().SetRoot(mainFlex, true).SetFocus(expandedView)

	expandedView.SetBorderPadding(1, 1, 1, 1).SetBorder(true)
	expandedView.SetSelectable(true, false)

	a, err := conn.GetFiles()
	aer := generateDirectoryTree(a)

	Update(*conn, aer.children, expandedView)

	Navbar.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			App.SetFocus(searchBar)
		} else if key == tcell.KeyBacktab {
			App.SetFocus(expandedView)
		}
	})
	expandedView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			App.SetFocus(Navbar)
		} else if key == tcell.KeyBacktab {
			App.SetFocus(searchBar)
		}
	})
	searchBar.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			App.SetFocus(expandedView)
		} else if key == tcell.KeyBacktab {
			App.SetFocus(Navbar)
		}
	})

	v, _ := conn.Status()
	Volume, _ = strconv.ParseInt(v["volume"], 10, 64)
	Random, _ = strconv.ParseBool(v["random"])
	Repeat, _ = strconv.ParseBool(v["repeat"])

	expandedView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Rune() == 108 {
			r, _ := expandedView.GetSelection()
			if len(aer.children[r].children) == 0 {
				id, _ := conn.AddId(aer.children[r].absolutePath, -1)
				conn.PlayId(id)
			} else {
				Update(*conn, aer.children[r].children, expandedView)
				aer = &aer.children[r]
			}
			return nil
		} else if e.Rune() == 112 {
			togglePlayBack(*conn)
			return nil
		} else if e.Rune() == 104 {
			if aer.parent != nil {
				Update(*conn, aer.parent.children, expandedView)
				aer = aer.parent
			}
			return nil
		} else if e.Rune() == 110 {
			conn.Next()
			return nil
		} else if e.Rune() == 99 {
			conn.Clear()
			return nil
		} else if e.Rune() == 78 {
			conn.Previous()
			return nil
		} else if e.Rune() == 97 {
			r, _ := expandedView.GetSelection()
			conn.Add(aer.children[r].absolutePath)
			return nil
		} else if e.Rune() == 122 {
			err := conn.Random(!Random)
			if err == nil {
				Random = !Random
			}
			return nil
		} else if e.Rune() == 114 {
			err := conn.Repeat(!Repeat)
			if err == nil {
				Repeat = !Repeat
			}
			return nil
		} else if e.Rune() == 45 {
			if Volume <= 0 {
				Volume = 0
			} else {
				Volume -= 10
			}
			conn.SetVolume(int(Volume))
			return nil
		} else if e.Rune() == 61 {
			if Volume >= 100 {
				Volume = 100
			} else {
				Volume += 10
			}
			conn.SetVolume(int(Volume))
			return nil
		} else {
			// fmt.Println(e.Rune())
			return e
		}
	})

	go func() {
		for {
			App.Draw()
			time.Sleep(time.Second)
		}
	}()

	if err := App.Run(); err != nil {
		panic(err)
	}
}
