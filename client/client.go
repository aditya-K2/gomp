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
	if status, err := Conn.Status(); err != nil {
		return err
	} else {
		if status["state"] == "play" {
			Conn.Pause(true)
		} else if status["state"] == "pause" || status["state"] == "stop" {
			Conn.Play(-1)
		}
	}
	return nil
}

// The GenerateContentSlice returns a slice of the content to be displayed on the Search View.
func GenerateContentSlice(selectedSuggestion string) ([]interface{}, error) {
	var _content []interface{}

	if strings.TrimRight(selectedSuggestion, " ") == "" {
		return nil, EmptySearchErr
	}

	if artists := getArtists(selectedSuggestion); len(artists) != 0 {
		_artists := []interface{}{}
		_arTitles := []interface{}{}
		_arAlbums := []interface{}{}
		for _, artist := range artists {
			_artists = append(_artists, artist)
			if albums := getArtistAlbums(artist); len(albums) != 0 {
				for _, album := range albums {
					for _, title := range getAlbumTitles(album) {
						_arTitles = append(_arTitles, [3]string{title, artist, album})
					}
					_arAlbums = append(_arAlbums, [2]string{album, artist})
				}
			}
		}
		_content = append(_content, WHITE_AND_BOLD+"Artists :")
		_content = append(_content, _artists...)
		_content = append(_content, WHITE_AND_BOLD+"Artist Albums :")
		_content = append(_content, _arAlbums...)
		_content = append(_content, WHITE_AND_BOLD+"Artist Titles :")
		_content = append(_content, _arTitles...)
	}

	if albums := getAlbums(selectedSuggestion); len(albums) != 0 {
		_albums := []interface{}{}
		_alTitles := []interface{}{}
		for _, album := range albums {
			artist := getTag([]string{"artist", "album", album})[0]
			_albums = append(_albums, [2]string{album, artist})
			for _, album := range albums {
				if titles := getAlbumTitles(album); len(titles) != 0 {
					for _, title := range titles {
						_alTitles = append(_alTitles, [3]string{title, artist, album})
					}
				}
			}
		}
		_content = append(_content, WHITE_AND_BOLD+"Albums :")
		_content = append(_content, _albums...)
		_content = append(_content, WHITE_AND_BOLD+"Album Titles :")
		_content = append(_content, _alTitles...)
	}

	if titles := getTitles(selectedSuggestion); len(titles) != 0 {
		_titles := []interface{}{}
		for _, title := range titles {
			artist := getTag([]string{"artist", "title", title})[0]
			album := getTag([]string{"album", "title", title})[0]
			_titles = append(_titles, [3]string{title, artist, album})
		}
		_content = append(_content, WHITE_AND_BOLD+"Titles :")
		_content = append(_content, _titles...)
	}

	return _content, nil
}

func getTag(filter []string) []string {
	if s, err := Conn.List(filter...); err != nil {
		return make([]string, 0)
	} else {
		return s
	}
}

func getArtists(artist string) []string {
	return getTag([]string{"artist", "artist", artist})
}

func getArtistAlbums(artist string) []string {
	return getTag([]string{"album", "artist", artist})
}

func getArtistTitles(artist string) []string {
	return getTag([]string{"title", "artist", artist})
}

func getAlbums(album string) []string {
	return getTag([]string{"album", "album", album})
}

func getAlbumTitles(album string) []string {
	return getTag([]string{"title", "album", album})
}

func getTitles(title string) []string {
	return getTag([]string{"title", "title", title})
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
	uris := getTag([]string{"file", "album", album})
	return add(uris)
}

func AddArtist(artist string) error {
	uris := getTag([]string{"file", "artist", artist})
	return add(uris)
}

func AddTitle(title string, play bool) error {
	uri := getTag([]string{"file", "title", title})[0]
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

func GetContent() []string {
	var p []string
	for _, v := range getTag([]string{"artist"}) {
		p = append(p, v)
	}
	for _, v := range getTag([]string{"album"}) {
		p = append(p, v)
	}
	for _, v := range getTag([]string{"title"}) {
		p = append(p, v)
	}
	return p
}
