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
)

type Application struct {
	App          *tview.Application
	ExpandedView *tview.Table
	Navbar       *tview.Table
	SearchBar    *tview.InputField
	ProgressBar  *ProgressBar
	Pages        *tview.Pages
}

func NewApplication() *Application {

	pBar := NewProgressBar()
	pBar.SetProgressFunc(progressFunction)
	expandedView := tview.NewTable()
	Navbar := tview.NewTable()
	searchBar := tview.NewInputField()
	searchBar.SetFieldBackgroundColor(tcell.ColorDefault)
	imagePreviewer := tview.NewBox()
	imagePreviewer.SetBorder(true)
	imagePreviewer.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		ImgX, ImgY, ImgW, ImgH = imagePreviewer.GetRect()
		return imagePreviewer.GetInnerRect()
	})

	expandedView.SetBackgroundColor(tcell.ColorDefault)
	Navbar.SetBackgroundColor(tcell.ColorDefault)
	searchBar.SetBackgroundColor(tcell.ColorDefault)
	imagePreviewer.SetBackgroundColor(tcell.ColorDefault)

	searchBar.SetTitle("Search").SetTitleAlign(tview.AlignLeft)
	searchBar.SetAutocompleteBackgroundColor(tcell.ColorBlack)
	searchBar.SetAutocompleteSelectBackgroundColor(tcell.ColorWhite)
	searchBar.SetAutocompleteMainTextColor(tcell.ColorDarkGray)
	searchBar.SetAutocompleteSelectedTextColor(tcell.ColorBlack)
	Navbar.SetBorder(true)
	Navbar.SetSelectable(true, false)
	Navbar.SetCell(0, 0, tview.NewTableCell("PlayList"))
	Navbar.SetCell(1, 0, tview.NewTableCell("Files"))
	Navbar.SetCell(2, 0, tview.NewTableCell("Most Played"))
	Navbar.SetCell(3, 0, tview.NewTableCell("Search"))

	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(Navbar, 0, 4, false).
		AddItem(imagePreviewer, 9, 3, false)

	sNavExpViewFlex := tview.NewFlex().
		AddItem(searchNavFlex, 17, 1, false).
		AddItem(expandedView, 0, 4, false)

	searchBar.SetBorder(true)
	searchBarFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBar, 3, 1, false).
		AddItem(sNavExpViewFlex, 0, 1, false)

	MainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBarFlex, 0, 8, false).
		AddItem(pBar, 5, 1, false)

	expandedView.SetBorderPadding(1, 1, 1, 1).SetBorder(true)
	expandedView.SetSelectable(true, false)

	rootPages := tview.NewPages()
	rootPages.AddPage("Main", MainFlex, true, true)

	App := tview.NewApplication()
	App.SetRoot(rootPages, true).SetFocus(expandedView)

	searchBar.SetDoneFunc(func(k tcell.Key) {
		switch k {
		case tcell.KeyEscape:
			{
				App.SetFocus(expandedView)
			}
		}
	})

	return &Application{
		App:          App,
		ExpandedView: expandedView,
		Navbar:       Navbar,
		SearchBar:    searchBar,
		ProgressBar:  pBar,
		Pages:        rootPages,
	}

}
