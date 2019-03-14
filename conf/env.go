package conf

import (
	"io/ioutil"
)

var defaultPath = "./config.yml"

type Env struct {
	//etcd的地址
	EtcdAddr string
}

func NewConfig(path string) *Env {
	if path == "" {
		path = defaultPath
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	yamlString := string(file)
	cfg, err := ParseYaml(yamlString)
	return &Env{
		EtcdAddr: cfg.getConfigString("etcd_addr"),
	}
}

func (cfg *Config) getConfigString(key string) string {
	s, _ := cfg.String(key)
	return s
}
