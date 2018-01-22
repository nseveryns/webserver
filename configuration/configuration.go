package configuration

import (
	"encoding/json"
	"log"
	"os"
)

//TODO Find out if these should be public
// if they should not, find way to decode json in private struct
type Configuration struct {
	Port        uint16
	Static      bool
	Directory   string
	ErrorPage   string
	LandingPage string
}

func LoadConfiguration(fileName string) (config Configuration) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Unable to open file.", err)
	}
	decoder := json.NewDecoder(file) // Load the structure from the file
	decoder.Decode(&config)
	return config
}
