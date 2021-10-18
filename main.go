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
var InsidePlaylist bool = true

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

	searchBar.SetBorder(true).SetTitle("Search").SetTitleAlign(tview.AlignLeft)
	Navbar.SetBorder(true)
	Navbar.SetSelectable(true, false)
	Navbar.SetCell(0, 0, tview.NewTableCell("PlayList"))
	Navbar.SetCell(1, 0, tview.NewTableCell("Files"))
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

	fileMap, err := conn.GetFiles()
	dirTree := generateDirectoryTree(fileMap)

	UpdatePlaylist(*conn, expandedView)

	_v, _ := conn.Status()
	Volume, _ = strconv.ParseInt(_v["volume"], 10, 64)
	Random, _ = strconv.ParseBool(_v["random"])
	Repeat, _ = strconv.ParseBool(_v["repeat"])

	expandedView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Rune() {
		case 108: // L : Key
			{
				r, _ := expandedView.GetSelection()
				if !InsidePlaylist {
					if len(dirTree.children[r].children) == 0 {
						id, _ := conn.AddId(dirTree.children[r].absolutePath, -1)
						conn.PlayId(id)
					} else {
						Update(*conn, dirTree.children[r].children, expandedView)
						dirTree = &dirTree.children[r]
					}
				} else {
					conn.Play(r)
				}
				return nil
			}
		case 112: // P : Key
			{
				togglePlayBack(*conn)
				return nil
			}
		case 104: // H : Key
			{
				if !InsidePlaylist {
					if dirTree.parent != nil {
						Update(*conn, dirTree.parent.children, expandedView)
						dirTree = dirTree.parent
					}
				}
				return nil
			}
		case 110: // N : Key
			{
				conn.Next()
				return nil
			}
		case 99: // C : Key
			{
				conn.Clear()
				if InsidePlaylist {
					UpdatePlaylist(*conn, expandedView)
				}
				return nil
			}
		case 78: // Shift - N : Key
			{
				conn.Previous()
				return nil
			}
		case 97: // A : Key
			{
				if !InsidePlaylist {
					r, _ := expandedView.GetSelection()
					conn.Add(dirTree.children[r].absolutePath)
				}
				return nil
			}
		case 122: // Z : Key
			{
				err := conn.Random(!Random)
				if err == nil {
					Random = !Random
				}
				return nil
			}
		case 114: // R : Key
			{
				err := conn.Repeat(!Repeat)
				if err == nil {
					Repeat = !Repeat
				}
				return nil
			}
		case 45: // Minus : Key
			{
				if Volume <= 0 {
					Volume = 0
				} else {
					Volume -= 10
				}
				conn.SetVolume(int(Volume))
				return nil
			}
		case 61: // Plus : Key
			{
				if Volume >= 100 {
					Volume = 100
				} else {
					Volume += 10
				}
				conn.SetVolume(int(Volume))
				return nil
			}
		case 50: // 2 : Key
			{
				InsidePlaylist = false
				Navbar.Select(1, 0)
				Update(*conn, dirTree.children, expandedView)
				return nil
			}
		case 49: // 1 : Key
			{
				InsidePlaylist = true
				Navbar.Select(0, 0)
				UpdatePlaylist(*conn, expandedView)
				return nil
			}
		case 51: // 3 : Key
			{
				InsidePlaylist = false
				Navbar.Select(2, 0)
				return nil
			}
		default:
			{
				return e
			}
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
