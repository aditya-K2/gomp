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
	if _, ok := ArtistM[selectedSuggestion]; ok {
		_content = append(_content, WHITE_AND_BOLD+"Artists :")
		_content = append(_content, selectedSuggestion)
		_content = append(_content, WHITE_AND_BOLD+"Artist Albums :")
		for albumName := range ArtistM[selectedSuggestion] {
			_content = append(_content, [2]string{albumName, selectedSuggestion})
		}
		_content = append(_content, WHITE_AND_BOLD+"Artist Tracks :")
		for albumName, trackList := range ArtistM[selectedSuggestion] {
			for track := range trackList {
				_content = append(_content, [3]string{track, selectedSuggestion, albumName})
			}
		}
	}
	if aMap := QueryAlbum(selectedSuggestion); len(aMap) != 0 {
		_content = append(_content, WHITE_AND_BOLD+"Albums :")
		for mSlice := range aMap {
			_content = append(_content, mSlice)
		}
		_content = append(_content, WHITE_AND_BOLD+"Album Tracks :")
		for a, pathSlice := range aMap {
			for _, path := range pathSlice {
				_content = append(_content, [3]string{path[0], a[1], a[0]})
			}
		}
	}
	if tMap := QueryTitle(selectedSuggestion); len(tMap) != 0 {
		_content = append(_content, WHITE_AND_BOLD+"Tracks :")
		for mSlice := range tMap {
			_content = append(_content, mSlice)
		}
	}
	return _content, nil
}

// Adds All tracks from a specified album to a playlist
func AddAlbum(alb string, artist string) error {
	clist := Conn.BeginCommandList()
	for _, fpath := range ArtistM[artist][alb] {
		clist.Add(fpath)
	}
	if err := clist.End(); err != nil {
		return errors.New("Could Not Add Album : " + alb)
	} else {
		return nil
	}
}

// Adds All tracks from a specified artist to a playlist
func AddArtist(artist string) error {
	clist := Conn.BeginCommandList()
	if val, ok := ArtistM[artist]; ok {
		for _, v := range val {
			for _, fpath := range v {
				clist.Add(fpath)
			}
		}
		if err := clist.End(); err != nil {
			return errors.New("Could Not Add Artist : " + artist)
		} else {
			return nil
		}
	} else {
		return errors.New("Could Not Add Artist : " + artist)
	}
}

// Adds Specified Track to the Playlist
func AddTitle(artist, alb, track string, addAndPlay bool) error {
	if addAndPlay {
		if id, err := Conn.AddID(ArtistM[artist][alb][track], -1); err != nil {
			return errors.New("Could Not Add Track : " + track)
		} else {
			if _err := Conn.PlayID(id); _err != nil {
				return errors.New("Could Not Play Track : " + track)
			}
		}
	} else {
		if err := Conn.Add(ArtistM[artist][alb][track]); err != nil {
			return errors.New("Could Not Add Track : " + track)
		}
	}
	return nil
}
