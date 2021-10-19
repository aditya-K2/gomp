package main

import (
	"github.com/fhs/gompd/mpd"
	"github.com/rivo/tview"
)

type Application struct {
	App          *tview.Application
	expandedView *tview.Table
	Navbar       *tview.Table
	searchBar    *tview.Table
	pBar         *progressBar
}

func newApplication(conn mpd.Client) *Application {

	var pBar *progressBar = newProgressBar(conn)
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
		AddItem(pBar.t, 5, 1, false)

	expandedView.SetBorderPadding(1, 1, 1, 1).SetBorder(true)
	expandedView.SetSelectable(true, false)

	App := tview.NewApplication()
	App.SetRoot(mainFlex, true).SetFocus(expandedView)

	return &Application{
		App:          App,
		expandedView: expandedView,
		Navbar:       Navbar,
		searchBar:    searchBar,
		pBar:         pBar,
	}

}
