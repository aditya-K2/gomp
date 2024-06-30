package ui

import "github.com/aditya-K2/tview"

var (
	borders = map[bool]map[string]rune{
		true: {
			"TopLeft":          '╭',
			"TopRight":         '╮',
			"BottomRight":      '╯',
			"BottomLeft":       '╰',
			"Vertical":         '│',
			"Horizontal":       '─',
			"TopLeftFocus":     '╭',
			"TopRightFocus":    '╮',
			"BottomRightFocus": '╯',
			"BottomLeftFocus":  '╰',
			"VerticalFocus":    '│',
			"HorizontalFocus":  '─',
		},
		false: {
			"TopLeft":          tview.Borders.TopLeft,
			"TopRight":         tview.Borders.TopRight,
			"BottomRight":      tview.Borders.BottomRight,
			"BottomLeft":       tview.Borders.BottomLeft,
			"Vertical":         tview.Borders.Vertical,
			"Horizontal":       tview.Borders.Horizontal,
			"TopLeftFocus":     tview.Borders.TopLeftFocus,
			"TopRightFocus":    tview.Borders.TopRightFocus,
			"BottomRightFocus": tview.Borders.BottomRightFocus,
			"BottomLeftFocus":  tview.Borders.BottomLeftFocus,
			"VerticalFocus":    tview.Borders.VerticalFocus,
			"HorizontalFocus":  tview.Borders.HorizontalFocus,
		},
	}
)

func SetBorderRunes(b bool) {
	tview.Borders.TopLeft = borders[b]["TopLeft"]
	tview.Borders.TopRight = borders[b]["TopRight"]
	tview.Borders.BottomRight = borders[b]["BottomRight"]
	tview.Borders.BottomLeft = borders[b]["BottomLeft"]
	tview.Borders.Vertical = borders[b]["Vertical"]
	tview.Borders.Horizontal = borders[b]["Horizontal"]
	tview.Borders.TopLeftFocus = borders[b]["TopLeftFocus"]
	tview.Borders.TopRightFocus = borders[b]["TopRightFocus"]
	tview.Borders.BottomRightFocus = borders[b]["BottomRightFocus"]
	tview.Borders.BottomLeftFocus = borders[b]["BottomLeftFocus"]
	tview.Borders.VerticalFocus = borders[b]["VerticalFocus"]
	tview.Borders.HorizontalFocus = borders[b]["HorizontalFocus"]
}
