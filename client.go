package main

import (
	"github.com/fhs/gompd/mpd"
)

func togglePlayBack(connection mpd.Client) error {
	status, err := connection.Status()
	if(status["state"] == "play" && err == nil){
		connection.Pause(true)
	} else if(status["state"] == "pause" && err == nil) {
		connection.Play(-1)
	}
	return err
}
