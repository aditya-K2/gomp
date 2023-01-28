package client

import (
	"errors"

	"github.com/aditya-K2/fuzzy"
	"github.com/fhs/gompd/v2/mpd"

	"strings"
)

var (
	Conn               *mpd.Client
	WHITE_AND_BOLD     string = "[white::b]"
	DirTree            *FileNode
	Matches            fuzzy.Matches
	SearchContentSlice []interface{}
	EmptySearchErr     = errors.New("empty Search String Provided")
)

func TogglePlayBack() error {
	status, err := Conn.Status()
	if status["state"] == "play" && err == nil {
		Conn.Pause(true)
	} else if status["state"] == "pause" && err == nil {
		Conn.Play(-1)
	}
	return err
}

// The GenerateContentSlice returns a slice of the content to be displayed on the Search View. The Slice is generated
// because the random nature of maps as they return values randomly hence the draw function keeps changing the order
// in which the results appear.
func GenerateContentSlice(selectedSuggestion string) ([]interface{}, error) {
	var _content []interface{}
	if strings.TrimRight(selectedSuggestion, " ") == "" {
		return nil, EmptySearchErr
	}
	if artists := getArtists(selectedSuggestion); len(artists) != 0 {
		_content = append(_content, WHITE_AND_BOLD+"Artists :")
		_arTitles := []interface{}{}
		_arAlbums := []interface{}{}
		for _, artist := range artists {
			_content = append(_content, artist)
			if albums := getArtistAlbums(artist); len(albums) != 0 {
				for _, album := range albums {
					for _, title := range getAlbumTitles(album) {
						_arTitles = append(_arTitles, [3]string{title, artist, album})
					}
					_arAlbums = append(_arAlbums, [2]string{album, artist})
				}
			}
		}
		_content = append(_content, WHITE_AND_BOLD+"Artist Albums :")
		_content = append(_content, _arAlbums)
		_content = append(_content, WHITE_AND_BOLD+"Artist Titles :")
		_content = append(_content, _arTitles)
	}
	if albums := getAlbums(selectedSuggestion); len(albums) != 0 {
		_content = append(_content, WHITE_AND_BOLD+"Albums :")
		_albums := []interface{}{}
		_alTitles := []interface{}{}
		for _, album := range albums {
			artist := GetTag([]string{"artist", "album", album})[0]
			for _, album := range albums {
				if titles := getAlbumTitles(album); len(titles) != 0 {
					for _, title := range titles {
						_alTitles = append(_alTitles, [3]string{title, artist, album})
					}
				}
				_albums = append(_albums, [2]string{album, artist})
			}
		}
		_content = append(_content, _albums...)
		_content = append(_content, WHITE_AND_BOLD+"Album Titles :")
		_content = append(_content, _alTitles...)
	}
	if titles := getTitles(selectedSuggestion); len(titles) != 0 {
		_content = append(_content, WHITE_AND_BOLD+"Titles :")
		_titles := []interface{}{}
		for _, title := range titles {
			artist := GetTag([]string{"artist", "title", title})[0]
			album := GetTag([]string{"album", "title", title})[0]
			_titles = append(_titles, [3]string{title, artist, album})
		}
	}

	return _content, nil
}

func GetTag(filter []string) []string {
	if s, err := Conn.List(filter...); err != nil {
		return make([]string, 0)
	} else {
		return s
	}
}

func getArtists(artist string) []string {
	return GetTag([]string{"artist", artist})
}

func getArtistAlbums(artist string) []string {
	return GetTag([]string{"album", "artist", artist})
}

func getArtistTitles(artist string) []string {
	return GetTag([]string{"title", "artist", artist})
}

func getAlbums(album string) []string {
	return GetTag([]string{"album", album})
}

func getAlbumTitles(album string) []string {
	return GetTag([]string{"title", "album", album})
}

func getTitles(title string) []string {
	return GetTag([]string{"title", title})
}

func add(uris []string) error {
	clist := Conn.BeginCommandList()
	for _, uri := range uris {
		clist.Add(uri)
	}
	if err := clist.End(); err != nil {
		return errors.New("Error Adding uris")
	} else {
		return nil
	}
}

func AddAlbum(album string) error {
	uris := GetTag([]string{"file", "album", album})
	return add(uris)
}

func AddArtist(artist string) error {
	uris := GetTag([]string{"file", "artist", artist})
	return add(uris)
}

func AddTitle(title string, play bool) error {
	uri := GetTag([]string{"file", "title", title})[0]
	if play {
		if id, err := Conn.AddID(uri, -1); err != nil {
			return errors.New("Could Not Add Track : " + title)
		} else {
			if _err := Conn.PlayID(id); _err != nil {
				return errors.New("Could Not Play Track : " + title)
			}
		}
	} else {
		if err := Conn.Add(uri); err != nil {
			return errors.New("Could Not Add Track : " + title)
		}
	}
	return nil
}
