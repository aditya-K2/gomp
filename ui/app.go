package ui

import (
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

var (
	ImgY int
	ImgW int
	ImgH int
	ImgX int
	Ui   *Application
)

type Application struct {
	App            *tview.Application
	MainS          *tview.Table
	Navbar         *tview.Table
	SearchBar      *tview.InputField
	ProgressBar    *ProgressBar
	Pages          *tview.Pages
	ImagePreviewer *tview.Box
}

func NewApplication(hideImage bool) *Application {

	pBar := NewProgressBar()
	mainS := NewMainS()
	searchbar := NewSearchBar()

	Navbar := tview.NewTable()
	imagePreviewer := tview.NewBox()

	imagePreviewer.SetBorder(true)

	Navbar.SetBackgroundColor(tcell.ColorDefault)
	imagePreviewer.SetBackgroundColor(tcell.ColorDefault)

	Navbar.SetBorder(true)
	Navbar.SetSelectable(true, false)

	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(Navbar, 0, 4, false)

	if !hideImage {
		searchNavFlex.AddItem(imagePreviewer, 9, 3, false)
	}

	Navbar.SetCell(0, 0, tview.NewTableCell("Queue"))
	Navbar.SetCell(1, 0, tview.NewTableCell("Files"))
	Navbar.SetCell(2, 0, tview.NewTableCell("Search"))

	sNavExpViewFlex := tview.NewFlex().
		AddItem(searchNavFlex, 17, 1, false).
		AddItem(mainS, 0, 4, false)

	searchBarFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchbar, 3, 1, false).
		AddItem(sNavExpViewFlex, 0, 1, false)

	MainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBarFlex, 0, 8, false).
		AddItem(pBar, 5, 1, false)

	rootPages := tview.NewPages()
	rootPages.AddPage("Main", MainFlex, true, true)

	App := tview.NewApplication()
	App.SetRoot(rootPages, true).SetFocus(mainS)

	return &Application{
		App:            App,
		MainS:          mainS,
		Navbar:         Navbar,
		SearchBar:      searchbar,
		ProgressBar:    pBar,
		Pages:          rootPages,
		ImagePreviewer: imagePreviewer,
	}

}
