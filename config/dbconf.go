package config

import (
	"encoding/json"
	"os"
)

type DbConf struct {
	User   string `json:"user"`
	DbName string `json:"db_name"`
	Pass   string `json:"pass"`
	Host   string `json:"host"`
}

func GetConfig() *DbConf {
	configFile, _ := os.Open("config.json")

	decoder := json.NewDecoder(configFile)
	var config DbConf
	decoder.Decode(&config)

	defer configFile.Close()

	return &config
}

func (d *DbConf) GetUser() string {
	return d.User
}

func (d *DbConf) GetDbName() string {
	return d.DbName
}

func (d *DbConf) GetPass() string {
	return d.Pass
}

func (d *DbConf) GetHost() string {
	return d.Host
}
