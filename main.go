// Demo code for the TextView primitive.
package main

import (
	"github.com/rivo/tview"
)

func main(){
	flext := tview.NewFlex().SetDirection(tview.FlexRow).
			 AddItem(tview.NewFlex().
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 2, false), 0, 7, false).
			 AddItem(tview.NewBox().SetBorder(true).SetTitle("Lmao"), 0, 1, false)

	if err := tview.NewApplication().SetRoot(flext, true).Run(); err != nil {
		panic(err)
	}
}
