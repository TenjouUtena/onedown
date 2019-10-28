package configuration

import (
	"encoding/json"
	"flag"
	"net/url"
	"os"
	"path"
	"time"
)

var config *Configuration

//Configuration Struct with Configuration data
type Configuration struct {
	Port                    int
	Logfile                 string
	LogLevel                string
	PuzzleSessionWriteDelay time.Duration
	CredentialFile          string
	CassandraHost           url.URL
}

func Get() *Configuration {
	if config == nil {
		configPath := flag.String("onedown-config",
			path.Join(os.Getenv("GOPATH"), "Configuration", "Configuration.json"),
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
