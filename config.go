package main

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type GstpConfig struct {
	Query          string        `yaml:"query"`
	CredentialFile string        `yaml:"credential"`
	TokenFile      string        `yaml:"token"`
	UserID         string        `yaml:"userid"`
	WebPath        string        `yaml:"webpath"`
	Port           string        `yaml:"port"`
	Filters        []GstpFilter  `yaml:"filters"`
	CheckInterval  time.Duration `yaml:"check_interval"`
}

type GstpFilter struct {
	SubjectRegex string `yaml:"subject_regex"`
	Label        string `yaml:"label"`
}

func ReadFile(filename string) (data []byte, err error) {

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	return data, err
}

func UnMarshallConfig(configdata []byte) (config *GstpConfig, err error) {

	err = yaml.Unmarshal(configdata, &config)
	if err != nil {
		return config, err
	}

	return config, err
}
