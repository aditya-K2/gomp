package ui

import (
	"github.com/aditya-K2/fuzzy"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

func NewSearchBar() *tview.InputField {
	searchbar := tview.NewInputField()

	searchbar.SetFieldBackgroundColor(tcell.ColorDefault)
	searchbar.SetBackgroundColor(tcell.ColorDefault)
	searchbar.SetTitle("Search").SetTitleAlign(tview.AlignLeft)
	searchbar.SetBorder(true)
	searchbar.SetAutocompleteStyles(
		config.Config.Colors.Autocomplete.Background(),
		tcell.StyleDefault,
		config.Config.Colors.Autocomplete.Style().Reverse(true))
	searchbar.SetAutocompleteMatchFieldWidth(true)
	searchbar.SetDoneFunc(func(k tcell.Key) {
		switch k {
		case tcell.KeyEscape:
			{
				Ui.App.SetFocus(Ui.MainS)
			}
		}
	})

	searchbar.SetAutocompleteFunc(func(c string) []string {
		if GetCurrentView() == nil {
			return []string{}
		}
		if GetCurrentView().Name() == "BuffSearchView" {
			return nil
		} else {
			if c != "" && c != " " && c != "  " {
				matches := fuzzy.Find(c, client.AllContent)
				var suggestions []string
				for i, match := range matches {
					if i == 10 {
						break
					}
					suggestions = append(suggestions, (match.Str))
				}
				return suggestions
			} else {
				return make([]string, 0)
			}
		}
	})

	searchbar.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyCtrlP {
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		}
		if e.Key() == tcell.KeyCtrlN {
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		}
		return e
	})

	searchbar.SetDoneFunc(func(e tcell.Key) {
		var err error
		if e == tcell.KeyEnter {
			Ui.MainS.Select(0, 0)
			if GetCurrentView().Name() == "BuffSearchView" {
				Ui.App.SetFocus(Ui.MainS)
			} else {
				SetCurrentView(SView)
				client.SearchContentSlice = nil
				client.SearchContentSlice, err = client.GenerateContentSlice(searchbar.GetText())
				if err != nil {
					SendNotification("Could Not Retrieve the Results")
				} else {
					searchbar.SetText("")
					Ui.App.SetFocus(Ui.MainS)
					Ui.Navbar.Select(2, 0)
				}
			}
		}
		if e == tcell.KeyEscape {
			if GetCurrentView().Name() == "BuffSearchView" {
				client.Matches = nil
				SetCurrentView(FView)
			}
			searchbar.SetText("")
			Ui.App.SetFocus(Ui.MainS)
		}
	})

	searchbar.SetChangedFunc(func(text string) {
		if GetCurrentView().Name() == "BuffSearchView" {
			var f client.FileNodes = client.DirTree.Children
			client.Matches = fuzzy.FindFrom(text, f)
			BuffSView.Update(Ui.MainS)
		}
	})

	return searchbar
}
