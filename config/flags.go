package config

import "flag"

var (
	VERSION = "1.0.0"
)

func ParseFlags() {
	flag.StringVar(&ConfigPath, "c", ConfigPath,
		"Specify The Directory where to check for config.yml file.")
	flag.BoolVar(&ShowVersion, "v", ShowVersion,
		"Show the current Version of gomp")
	flag.Parse()
}
