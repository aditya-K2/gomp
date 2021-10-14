// Demo code for the TextView primitive.
package main

import (
	"github.com/rivo/tview"
	"log"
	"github.com/fhs/gompd/mpd"
	"github.com/gdamore/tcell/v2"
	"fmt"
)

func getContent(row, column int) {
}

func main(){
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	expandedView := tview.NewTable()
	Navbar := tview.NewTable().SetBorder(true)
	searchBar := tview.NewBox().SetBorder(true)
	progressBar := tview.NewBox().SetBorder(true)

	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(Navbar, 0, 1, false).
				AddItem(searchBar, 0, 7, false)

	sNavExpViewFlex := tview.NewFlex().
				AddItem(searchNavFlex, 0, 1, false).
				AddItem(expandedView, 0, 4, false)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(sNavExpViewFlex, 0, 8, false).
				AddItem(progressBar, 0, 1, false)

	expandedView.SetBorderPadding(1,1,1,1).SetBorder(true)
	expandedView.SetSelectable(true, false)
	fmt.Println(expandedView.HasFocus())

	a, err := conn.GetFiles()
	aer := generateDirectoryTree(a);
	ec := []string{"NothingHappens", "Remember When.mp3"}

	Update(*conn, aer.children, ec, expandedView)
	expandedView.SetDoneFunc(func(key tcell.Key){
		if key == tcell.KeyLeft {
			r, c := expandedView.GetSelection()
			fmt.Println(join(ec) + expandedView.GetCell(r, c).Text)
		}
	}).SetSelectedFunc(func(row, column int){
		fmt.Println(join(ec) + expandedView.GetCell(row, column).Text)
	})
	if err := tview.NewApplication().SetRoot(mainFlex, true).SetFocus(expandedView).Run(); err != nil {
		panic(err)
	}
}

