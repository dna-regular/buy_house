package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type ProxyConf struct {
	RawProxy string `json:"raw_proxy"`
}

type Config struct {
	Proxies ProxyConf
}

func ParseConf(confPath string) (*Config, error) {
	jsonFile, err := os.Open(confPath)
	if err != nil {
		log.Println("Error opening json file:", err)
		return nil, err
	}

	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	var config Config
	for {
		err := decoder.Decode(&config)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error decoding json:", err)
			return nil, err
		}

	}
	return &config, err
}
