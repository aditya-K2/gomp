package client

import (
	"errors"
	"fmt"

	"github.com/fhs/gompd/mpd"

	"strings"
)

var (
	CONN               *mpd.Client
	ArtistTree         map[string]map[string]map[string]string
	NotificationServer interface {
		Send(string)
	}
	WHITE_AND_BOLD string = "[white::b]"
)

func SetNotificationServer(n interface{ Send(string) }) {
	NotificationServer = n
}

func SetConnection(c *mpd.Client) {
	CONN = c
}

func TogglePlayBack() error {
	status, err := CONN.Status()
	if status["state"] == "play" && err == nil {
		CONN.Pause(true)
	} else if status["state"] == "pause" && err == nil {
		CONN.Play(-1)
	}
	return err
}

// The GenerateContentSlice returns a slice of the content to be displayed on the Search View. The Slice is generated
// because the random nature of maps as they return values randomly hence the draw function keeps changing the order
// in which the results appear.
func GenerateContentSlice(selectedSuggestion string) ([]interface{}, error) {
	var ContentSlice []interface{}
	if strings.TrimRight(selectedSuggestion, " ") == "" {
		NotificationServer.Send("Empty Search!")
		return nil, errors.New("empty Search String Provided")
	}
	if _, ok := ArtistTree[selectedSuggestion]; ok {
		ContentSlice = append(ContentSlice, WHITE_AND_BOLD+"Artists :")
		ContentSlice = append(ContentSlice, selectedSuggestion)
		ContentSlice = append(ContentSlice, WHITE_AND_BOLD+"Artist Albums :")
		for albumName := range ArtistTree[selectedSuggestion] {
			ContentSlice = append(ContentSlice, [2]string{albumName, selectedSuggestion})
		}
		ContentSlice = append(ContentSlice, WHITE_AND_BOLD+"Artist Tracks :")
		for albumName, trackList := range ArtistTree[selectedSuggestion] {
			for track := range trackList {
				ContentSlice = append(ContentSlice, [3]string{track, selectedSuggestion, albumName})
			}
		}
	}
	if aMap := QueryArtistTreeForAlbums(ArtistTree, selectedSuggestion); len(aMap) != 0 {
		ContentSlice = append(ContentSlice, WHITE_AND_BOLD+"Albums :")
		for mSlice := range aMap {
			ContentSlice = append(ContentSlice, mSlice)
		}
		ContentSlice = append(ContentSlice, WHITE_AND_BOLD+"Album Tracks :")
		for a, pathSlice := range aMap {
			for _, path := range pathSlice {
				ContentSlice = append(ContentSlice, [3]string{path[0], a[1], a[0]})
			}
		}
	}
	if tMap := QueryArtistTreeForTracks(ArtistTree, selectedSuggestion); len(tMap) != 0 {
		ContentSlice = append(ContentSlice, WHITE_AND_BOLD+"Tracks :")
		for mSlice := range tMap {
			ContentSlice = append(ContentSlice, mSlice)
		}
	}
	return ContentSlice, nil
}

// GenerateArtistTree Artist Tree is a map of Artist to their Album Map
// Album Tree is a map of the tracks in that particular album.
func GenerateArtistTree() (map[string]map[string]map[string]string, error) {
	ArtistTree = make(map[string]map[string]map[string]string)
	AllInfo, err := CONN.ListAllInfo("/")
	if err == nil {
		for _, i := range AllInfo {
			if _, ArtistExists := ArtistTree[i["Artist"]]; !ArtistExists {
				ArtistTree[i["Artist"]] = make(map[string]map[string]string)
			}
			if _, AlbumExists := ArtistTree[i["Artist"]][i["Album"]]; !AlbumExists {
				ArtistTree[i["Artist"]][i["Album"]] = make(map[string]string)
			}
			if _, TitleExists := ArtistTree[i["Artist"]][i["Album"]][i["Title"]]; !TitleExists {
				ArtistTree[i["Artist"]][i["Album"]][i["Title"]] = i["file"]
			}
		}
		return ArtistTree, nil
	} else {
		return nil, errors.New("Could Not Generate Artist Tree")
	}
}

func PrintArtistTree(a map[string]map[string]map[string]string) {
	for k, v := range a {
		fmt.Println(k, " : ")
		for k1, v1 := range v {
			fmt.Println("\t|---", k1, " : ")
			for k2 := range v1 {
				fmt.Println("\t\t|---", k2)
			}
		}
	}
}

//  Adds All tracks from a specified album to a playlist
func AddAlbum(a map[string]map[string]map[string]string, alb string, artist string) {
	for _, v := range a[artist][alb] {
		err := CONN.Add(v)
		if err != nil {
			NotificationServer.Send("Could Not Add Song : " + v)
		}
	}
	NotificationServer.Send("Album Added : " + alb)
}

//  Adds All tracks from a specified artist to a playlist
func AddArtist(a map[string]map[string]map[string]string, artist string) {
	if val, ok := a[artist]; ok {
		for _, v := range val {
			for _, path := range v {
				err := CONN.Add(path)
				if err != nil {
					NotificationServer.Send("Could Not Add Song : " + path)
				}
			}
		}
		NotificationServer.Send("Artist Added : " + artist)
	}
}

//  Adds Specified Track to the Playlist
func AddTitle(a map[string]map[string]map[string]string, artist, alb, track string, addAndPlay bool) {
	if addAndPlay {
		id, err := CONN.AddId(a[artist][alb][track], -1)
		CONN.PlayId(id)
		if err != nil {
			NotificationServer.Send("Could Not Add Track : " + track)
		}
	} else {
		err := CONN.Add(a[artist][alb][track])
		if err != nil {
			NotificationServer.Send("Could Not Add Track : " + track)
		}
	}
	NotificationServer.Send("Track Added : " + track)
}

/* Querys the Artist Tree for a track and returns a TrackMap (i.e [3]string{artist, album, track} -> Path) which will help us
to add tracks to the playlist */
func QueryArtistTreeForTracks(a map[string]map[string]map[string]string, track string) map[[3]string]string {
	TrackMap := make(map[[3]string]string)
	for artistName, albumMap := range a {
		for albumName, trackList := range albumMap {
			for trackName, path := range trackList {
				if trackName == track {
					TrackMap[[3]string{trackName, artistName, albumName}] = path
				}
			}
		}
	}
	return TrackMap
}

/* Querys the Artist Tree for an album and returns a AlbumMap (i.e [3]string{artist, album } ->[]Path of songs in the album)
which will help us to add all album tracks to the playlist */
func QueryArtistTreeForAlbums(a map[string]map[string]map[string]string, album string) map[[2]string][][2]string {
	AlbumMap := make(map[[2]string][][2]string)
	for artistName, albumMap := range a {
		for albumName, trackList := range albumMap {
			if albumName == album {
				var pathSlice [][2]string
				for trackName, path := range trackList {
					pathSlice = append(pathSlice, [2]string{trackName, path})
				}
				AlbumMap[[2]string{albumName, artistName}] = pathSlice
			}
		}
	}
	return AlbumMap
}

func AddToPlaylist(a interface{}, addAndPlay bool) {
	switch a.(type) {
	case [3]string:
		{
			b := a.([3]string)
			AddTitle(ArtistTree, b[1], b[2], b[0], addAndPlay)
		}
	case [2]string:
		{
			b := a.([2]string)
			AddAlbum(ArtistTree, b[0], b[1])
		}
	case string:
		{
			b := a.(string)
			AddArtist(ArtistTree, b)
		}
	}
}
