package config

import (
	"flag"
	"os"

	"github.com/aditya-K2/gomp/utils"
)

var (
	VERSION = "1.0.0"
)

func ParseFlags() {
	flag.StringVar(&ConfigPath, "c", ConfigPath,
		"Specify The Directory where to check for config.yml file.")
	flag.BoolVar(&showVersion, "v", showVersion,
		"Show the current Version of gomp")
	flag.Parse()
	if showVersion {
		utils.Print("WHITE", "gomp ")
		utils.Print("BLUE", VERSION+"\n")
		utils.Print("WHITE", "Report any issues at: ")
		utils.Print("BLUE", "https://github.com/aditya-K2/gomp/issues\n")
		os.Exit(0)
	}
}
