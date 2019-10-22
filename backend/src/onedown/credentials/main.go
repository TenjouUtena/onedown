package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

func LoadCredentials(fileName string) {
	var c Credentials
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Printf("File error: v\n", err)
		os.Exit(1)
	}

	json.Unmarshal(file, &c)
}
