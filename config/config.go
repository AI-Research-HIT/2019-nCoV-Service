package config

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
)

type ConfigT struct {
	PidFile    string `json:"pidfile"`
	LogFile    string `json:"logfile"`
	LogLevel   int    `json:"loglevel"`
	ServerAddr string `json:"serveraddr"`
}

func (c *ConfigT) Check() bool {
	if c.LogLevel <= 0 && c.LogLevel > 5 {
		c.LogLevel = 2
	}

	return len(c.PidFile) > 0 && len(c.LogFile) > 0 &&
		len(c.ServerAddr) > 0
}

var conf string

func init() {
	const (
		defaultConf = "2019-nCoV-Service.conf"
		usage       = "the config file"
	)
	flag.StringVar(&conf, "conf", defaultConf, usage)
	flag.StringVar(&conf, "c", defaultConf, usage+" (shorthand)")
}

var Config ConfigT

func ParseConfig() error {
	flag.Parse()
	configData, err := ioutil.ReadFile(conf)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configData, &Config)
	if err != nil {
		return err
	}

	if !Config.Check() {
		return errors.New("config file check failed")
	}

	return nil
}
