package test

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	DriverName       string
	DriverDataSource string
}

type allConfig struct {
	config
}

var conf *allConfig

const confFile = "./config.json"

func NewConfig() {
	var cjson config
	cfile, err := ioutil.ReadFile(confFile)
	if err != nil {
		panic("fail to load configuration file")
	}

	if err := json.Unmarshal(cfile, &cjson); err != nil {
		panic(err)
	}

	conf = &allConfig{config: cjson}
}
