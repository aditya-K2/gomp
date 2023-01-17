package client

import (
	"errors"
	"fmt"
)

type Artists map[string]Albums
type Albums map[string]Titles
type Titles map[string]string

var (
	ArtistM Artists
)

// Querys the Artist Map for a track and returns a TrackMap
// (i.e [3]string{artist, album, track} -> Path) which will help us
// to add tracks to the playlist
func QueryTitle(track string) map[[3]string]string {
	titleMap := make(map[[3]string]string)
	for _artist, albums := range ArtistM {
		for _album, titles := range albums {
			for title, fpath := range titles {
				if title == track {
					titleMap[[3]string{title, _artist, _album}] = fpath
				}
			}
		}
	}
	return titleMap
}

// Querys the Artist Map for an album and returns a AlbumMap
// (i.e [3]string{artist, album } ->[]Path of songs in the album)
// which will help us to add all album tracks to the playlist
func QueryAlbum(album string) map[[2]string][][2]string {
	albumMap := make(map[[2]string][][2]string)
	for _artist, albums := range ArtistM {
		for _album, titles := range albums {
			if _album == album {
				var pslice [][2]string
				for title, fpath := range titles {
					pslice = append(pslice, [2]string{title, fpath})
				}
				albumMap[[2]string{_album, _artist}] = pslice
			}
		}
	}
	return albumMap
}

// GenerateArtistMap Artist Tree is a map of Artist to their Album Map
// Album Tree is a map of the tracks in that particular album.
func GenerateArtistMap() error {
	ArtistM = make(Artists)
	if info, err := Conn.ListAllInfo("/"); err == nil {
		for _, i := range info {
			artist := i["Artist"]
			album := i["Album"]
			title := i["Title"]
			if artist == "" {
				artist = "Unknown Artists"
			}
			if album == "" {
				album = "Unknown Artists"
			}
			if title == "" {
				title = i["file"]
			}
			if _, ArtistExists := ArtistM[artist]; !ArtistExists {
				ArtistM[artist] = make(map[string]Titles)
			}
			if _, AlbumExists := ArtistM[artist][album]; !AlbumExists {
				ArtistM[artist][album] = make(map[string]string)
			}
			if _, TitleExists := ArtistM[artist][album][title]; !TitleExists {
				ArtistM[artist][album][title] = i["file"]
			}
		}
	} else {
		return errors.New("Could Not Generate Artist Tree!\n")
	}
	return nil
}

// For Debugging
func PrintArtistTree() {
	for k, v := range ArtistM {
		fmt.Println(k, " : ")
		for k1, v1 := range v {
			fmt.Println("\t|---", k1, " : ")
			for k2 := range v1 {
				fmt.Println("\t\t|---", k2)
			}
		}
	}
}
