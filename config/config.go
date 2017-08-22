package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var ConfigFile string

func init() {
	binaryAbsolute, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal("not able to find absolute path to config")
	}

	defaultConfigPath := filepath.Join(filepath.Dir(binaryAbsolute), "config.json")

	flag.StringVar(&ConfigFile, "config", defaultConfigPath, "Configuration file")
}

type Configuration struct {
	Database struct {
		Host     string `json:"psqlhost"`
		Name     string `json:"psqlname"`
		User     string `json:"psqluser"`
		Password string `json:"psqlpassword"`
		Port     string `json:"psqlport"`
	} `json:"psqlinfo"`
	// TODO : add (path to signatures/releases static folder)
}

func (config *Configuration) Parse(filename string) error {
	ConfigJSON, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}
	return json.Unmarshal(ConfigJSON, config)

}

func Get() (*Configuration, error) {
	config := &Configuration{}
	err := config.Parse(ConfigFile)

	return config, err
}
