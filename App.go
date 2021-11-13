package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var IMG_X, IMG_Y, IMG_W, IMG_H int

type Application struct {
	App          *tview.Application
	expandedView *tview.Table
	Navbar       *tview.Table
	searchBar    *tview.Table
	pBar         *progressBar
	Pages        *tview.Pages
}

func newApplication(r *Renderer) *Application {

	var pBar *progressBar = newProgressBar(r)
	expandedView := tview.NewTable()
	Navbar := tview.NewTable()
	searchBar := tview.NewTable()
	imagePreviewer := tview.NewBox()
	imagePreviewer.SetBorder(true)
	imagePreviewer.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		IMG_X, IMG_Y, IMG_W, IMG_H = imagePreviewer.GetRect()
		return imagePreviewer.GetInnerRect()
	})

	searchBar.SetBorder(true).SetTitle("Search").SetTitleAlign(tview.AlignLeft)
	Navbar.SetBorder(true)
	Navbar.SetSelectable(true, false)
	Navbar.SetCell(0, 0, tview.NewTableCell("PlayList"))
	Navbar.SetCell(1, 0, tview.NewTableCell("Files"))
	Navbar.SetCell(2, 0, tview.NewTableCell("Most Played"))

	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBar, 3, 1, false).
		AddItem(Navbar, 0, 4, false).
		AddItem(imagePreviewer, 9, 3, false)

	sNavExpViewFlex := tview.NewFlex().
		AddItem(searchNavFlex, 17, 1, false).
		AddItem(expandedView, 0, 4, false)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sNavExpViewFlex, 0, 8, false).
		AddItem(pBar.t, 5, 1, false)

	expandedView.SetBorderPadding(1, 1, 1, 1).SetBorder(true)
	expandedView.SetSelectable(true, false)

	rootPages := tview.NewPages()
	rootPages.AddPage("Main", mainFlex, true, true)

	App := tview.NewApplication()
	App.SetRoot(rootPages, true).SetFocus(expandedView)

	return &Application{
		App:          App,
		expandedView: expandedView,
		Navbar:       Navbar,
		searchBar:    searchBar,
		pBar:         pBar,
		Pages:        rootPages,
	}

}
