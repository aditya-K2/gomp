package config

import (
	"errors"
)

// # Key Mappings
//
// ## Following Keys can be used for Mappings
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
// See [config/kMap.go](https://github.com/aditya-K2/gomp/blob/master/config/kMap.go) for more information
//
// For mapping a key to some function use the following format:
//
// ```yml
// Function: [ firstMapping, secondMapping, thirdMapping]
// ```
// for e.g
//
// ```yml
// togglePlayBack : [ "p", "SPACE", "[" ] # using the quotes is neccessary.
// ```

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

	//
	// Following functions are provided :
	//
	// |          Functions                 | Default Key Mapping |
	// |------------------------------------|---------------------|
	KEY_MAP = map[int]string{
		// |     showChildrenContent            |    <kbd>l</kbd>     |
		108: "showChildrenContent",
		// |     togglePlayBack                 |    <kbd>p</kbd>     |
		112: "togglePlayBack",
		// |     showParentContent              |    <kbd>h</kbd>     |
		104: "showParentContent",
		// |     nextSong                       |    <kbd>n</kbd>     |
		110: "nextSong",
		// |     clearPlaylist                  |    <kbd>c</kbd>     |
		99: "clearPlaylist",
		// |     previousSong                   |    <kbd>N</kbd>     |
		78: "previousSong",
		// |     addToPlaylist                  |    <kbd>a</kbd>     |
		97: "addToPlaylist",
		// |     toggleRandom                   |    <kbd>z</kbd>     |
		122: "toggleRandom",
		// |     toggleRepeat                   |    <kbd>r</kbd>     |
		114: "toggleRepeat",
		// |     decreaseVolume                 |    <kbd>-</kbd>     |
		45: "decreaseVolume",
		// |     increaseVolume                 |    <kbd>+</kbd>     |
		61: "increaseVolume",
		// |     navigateToPlaylist             |    <kbd>1</kbd>     |
		49: "navigateToPlaylist",
		// |     navigateToFiles                |    <kbd>2</kbd>     |
		50: "navigateToFiles",
		// |     navigateToSearch               |    <kbd>3</kbd>     |
		51: "navigateToSearch",
		// |     quit                           |    <kbd>q</kbd>     |
		113: "quit",
		// |     stop                           |    <kbd>s</kbd>     |
		115: "stop",
		// |     updateDB                       |    <kbd>u</kbd>     |
		117: "updateDB",
		// |     deleteSongFromPlaylist         |    <kbd>d</kbd>     |
		100: "deleteSongFromPlaylist",
		// |     FocusSearch                    |    <kbd>?</kbd>     |
		63: "FocusSearch",
		// |     FocusBuffSearch                |    <kbd>/</kbd>     |
		47: "FocusBuffSearch",
		// |     SeekBackward                   |    <kbd>b</kbd>     |
		98: "SeekBackward",
		// |     SeekForward                    |    <kbd>f</kbd>     |
		102: "SeekForward",
		//
		//
		// -------
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
