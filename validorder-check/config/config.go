package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GConfig Config
var GPreConfig PreConfig

func LoadConfig(path string) *Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("Read yml config[%s] err[%v]", path, err).Error())
	}

	err = yaml.Unmarshal([]byte(data), &GConfig)
	if err != nil {
		panic(fmt.Errorf("yml.Unmarshal config[%s] err[%v]", path, err).Error())
	}
	err = PreProcess()
	if err != nil {
		panic(fmt.Errorf("PreProcess config[%s] err[%v]", path, err).Error())
	}
	return &GConfig
}

func PreProcess() error {

	return nil
}

type PreConfig struct{

}

type Config struct {
	Server  System      `yaml:"system"`
	ChuangLan ChuangLan     `yaml:"chuanglan"`
}

type System struct {
	OutputPath    string    `yaml:"output_path"`
	InputPath   string        `yaml:"input_path"`
	LogPath string      `yaml:"log_path"`
	Intvl    int      `yaml:"interval"`
	PhoneCheckFlag   int     `yaml:"phone_check"`
	IdcardCheckFlag   int     `yaml:"idcard_check"`
}

type ChuangLan struct {
	AppId   string      `yaml:"appId"`
	AppKey  string      `yaml:"appKey"`
	IdcardCheck_url string    `yaml:"idcardcheck_url"`
	UnnCheck_url   string     `yaml:"unncheck_url"`
	Host      string
}

