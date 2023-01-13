package main

import (
	"os"
	"time"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/database"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/gomp/watchers"

	"github.com/fhs/gompd/v2/mpd"
)

func main() {
	config.ReadConfig()

	var mcerr error
	client.Conn, mcerr = mpd.Dial(
		utils.GetNetwork(config.Config.NetworkType,
			config.Config.Port,
			config.Config.NetworkAddress))

	if mcerr != nil {
		utils.Print("RED", "There was a Problem Connecting to the MPD Server\nTry Checking:\n")
		utils.Print("GREEN", "- if MPD Server is up and running.\n")
		utils.Print("GREEN", "- if you have provided correct port in the config.yml\n")
		os.Exit(-1)
	}

	Conn := client.Conn
	defer Conn.Close()

	cache.SetCacheDir(config.Config.CacheDir)

	watchers.Init()

	ui.Ui = ui.NewApplication()
	ui.Ui.ProgressBar.SetProgressFunc(watchers.ProgressFunction)

	if fileMap, err := Conn.ListAllInfo("/"); err != nil {
		utils.Print("RED", "Could Not Generate the File Map\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(err)
	} else {
		client.DirTree = client.GenerateDirectoryTree(fileMap)
	}

	// Generating the Directory Tree for File Navigation.

	if ArtistTree, err := client.GenerateArtistTree(); err != nil {
		utils.Print("RED", "Could Not Generate the ArtistTree\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(err)
	} else {
		// Used for Fuzzy Searching
		ArtistTreeContent := utils.ConvertToArray(ArtistTree)
		ui.SetArtistTreeContent(ArtistTreeContent)
	}

	ui.InitNotifier()

	// This Function Is Responsible for Changing the Focus it uses the Focus Map and Based on it Chooses
	// the Draw Function
	watchers.StartPlaylistWatcher()
	watchers.StartMPListener()
	watchers.StartRectWatcher()
	database.Publish()

	ui.SetCurrentView(&ui.PView)

	// Generating the Key Map Based on the Function Map Here Basically the Values will be flipped
	// In the config if togglePlayBack is mapped to [ T , P, SPACE ] then here Basically we will receive a map
	// for each event T, P, SPACE mapped to the same function togglePlayBack
	config.GenerateKeyMap(ui.FuncMap)

	go func() {
		for {
			ui.Ui.App.Draw()
			time.Sleep(time.Second)
		}
	}()

	if err := ui.Ui.App.Run(); err != nil {
		panic(err)
	}

	defer watchers.DBCheck()
}
