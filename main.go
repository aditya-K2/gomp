package main

import (
	"os"
	"time"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/watchers"
	"github.com/aditya-K2/utils"

	"github.com/fhs/gompd/v2/mpd"
)

func main() {
	config.ParseFlags()
	config.ReadConfig()

	var mcerr error
	network, addr := utils.GetNetwork(config.Config.NetworkType, config.Config.Port, config.Config.NetworkAddress)
	if config.Config.Password != "" {
		client.Conn, mcerr = mpd.DialAuthenticated(network, addr, config.Config.Password)
	} else {
		client.Conn, mcerr = mpd.Dial(network, addr)
	}

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

	ui.Ui = ui.NewApplication(config.Config.HideImage)
	ui.Ui.ProgressBar.SetProgressFunc(watchers.ProgressFunction)

	// Generating the Directory Tree for File Navigation.
	if fileMap, err := Conn.ListAllInfo("/"); err != nil {
		utils.Print("RED", "Could Not Generate the File Map\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(err)
	} else {
		client.DirTree = client.GenerateDirectoryTree(fileMap)
	}

	client.GetContent()

	ui.InitNotifier()

	watchers.StartPlaylistWatcher()

    ui.SetBorderRunes(config.Config.RoundedCorners)

	if !config.Config.HideImage {
		watchers.StartRectWatcher()
	}

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
}
