package config

import (
	"errors"
)

// # Key Mappings
//
// # Following Keys can be used for Mappings
//
// | Keys            | Using them in Config  |
// |-----------------|-----------------------|
// | a-z             | a-z                   |
// | A-Z             | A-z                   |
// | {,},(,),[,],<,> | {,},(,),[,],<,>       |
// | Enter(Return)   | ENTER/RETURN          |
// | Tab             | TAB                   |
// | Space           | SPACE                 |
//
// See config/kMap.go for more information
//
// For mapping a key to some function use the following format:
//
// ```yml
// Function: [ firstMapping, secondMapping, thirdMapping]
// ```
// for e.g
//
// ```yml
// togglePlayBack : [ "p", "TAB", "[" ] # using the quotes is neccessary.
// ```
//
// Following functions are provided :
//
// |          Functions                 | Default Key Mapping |
// |------------------------------------|---------------------|
// |     showChildrenContent            |    <kbd>l</kbd>     |
// |     togglePlayBack                 |    <kbd>p</kbd>     |
// |     showParentContent              |    <kbd>h</kbd>     |
// |     nextSong                       |    <kbd>n</kbd>     |
// |     clearPlaylist                  |    <kbd>c</kbd>     |
// |     previousSong                   |    <kbd>N</kbd>     |
// |     addToPlaylist                  |    <kbd>a</kbd>     |
// |     toggleRandom                   |    <kbd>z</kbd>     |
// |     toggleRepeat                   |    <kbd>r</kbd>     |
// |     decreaseVolume                 |    <kbd>-</kbd>     |
// |     increaseVolume                 |    <kbd>+</kbd>     |
// |     navigateToFiles                |    <kbd>2</kbd>     |
// |     navigateToPlaylist             |    <kbd>1</kbd>     |
// |     navigateToMostPlayed           |    <kbd>3</kbd>     |
// |     navigateToSearch               |    <kbd>4</kbd>     |
// |     quit                           |    <kbd>q</kbd>     |
// |     stop                           |    <kbd>s</kbd>     |
// |     updateDB                       |    <kbd>u</kbd>     |
// |     deleteSongFromPlaylist         |    <kbd>d</kbd>     |
// |     FocusSearch                    |    <kbd>?</kbd>     |
// |     FocusBuffSearch                |    <kbd>/</kbd>     |
// |     SeekForward                    |    <kbd>f</kbd>     |
// |     SeekBackward                   |    <kbd>b</kbd>     |

var (
	SPECIAL_KEYS = map[string]int{
		"TAB":    9,
		"RETURN": 13,
		"ENTER":  13,
		"SPACE":  32,
		"[":      91,
		"]":      93,
		"(":      40,
		")":      41,
		"{":      123,
		"}":      125,
		"<":      60,
		">":      62,
		"?":      63,
		"/":      47,
		";":      59,
		":":      58,
		"'":      39,
		"\"":     34,
	}

	/*Generating Default KEY_MAP which will then later be changed by GenerateKeyMap*/

	KEY_MAP = map[int]string{
		108: "showChildrenContent",
		112: "togglePlayBack",
		104: "showParentContent",
		110: "nextSong",
		99:  "clearPlaylist",
		78:  "previousSong",
		97:  "addToPlaylist",
		122: "toggleRandom",
		114: "toggleRepeat",
		45:  "decreaseVolume",
		61:  "increaseVolume",
		50:  "navigateToFiles",
		49:  "navigateToPlaylist",
		51:  "navigateToMostPlayed",
		52:  "navigateToSearch",
		113: "quit",
		115: "stop",
		117: "updateDB",
		100: "deleteSongFromPlaylist",
		63:  "FocusSearch",
		47:  "FocusBuffSearch",
		98:  "SeekBackward",
		102: "SeekForward",
	}
)

func GetAsciiValue(s string) (int, error) {
	if len([]rune(s)) == 1 {
		char := []rune(s)[0]
		if (int(char) >= 65 && int(char) <= 90) || (int(char) >= 97 && int(char) <= 122) {
			return int(char), nil
		} else if val, ok := SPECIAL_KEYS[s]; ok {
			return val, nil
		} else {
			return -1, errors.New("Not Found in the range")
		}
	} else if val, ok := SPECIAL_KEYS[s]; ok {
		return val, nil
	} else {
		return -1, errors.New("Not Found")
	}
}
