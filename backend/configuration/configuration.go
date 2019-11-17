package configuration

import (
	"encoding/json"
	"flag"
	"os"
	"path"
	"time"
)

var config *Configuration

//Configuration Struct with Configuration data
type Configuration struct {
	Port                      int
	Logfile                   string
	LogLevel                  string
	PuzzleSessionWriteDelayMs time.Duration
	CredentialFile            string
	CassandraHost             string
}

func Get() *Configuration {
	if config == nil {
		configPath := flag.String("onedown-config",
			path.Join(os.Getenv("GOPATH"), "configuration", "configuration.json"),
			"Path to Configuration JSON file.")
		file, err := os.Open(*configPath)
		if err != nil {
			panic("Unable to open Configuration file!")
		}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(config)
		if err != nil {
			panic("Unable to open Configuration file!")
		}
	}
	return config
}
