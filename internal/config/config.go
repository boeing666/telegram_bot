package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

type ConfigStructure struct {
	AppID      int    `json:"app_id"`
	AppHash    string `json:"app_hash"`
	APIToken   string `json:"api_token"`
	DbSettings struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"db_settings"`
}

func Init() (*ConfigStructure, error) {
	jsonFile, err := os.ReadFile("configs/config.json")
	if err != nil {
		return nil, fmt.Errorf("error on reading config: %w", err)
	}
	var config ConfigStructure
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error on parsing config file: %w", err)
	}
	return &config, nil
}

func (config *ConfigStructure) GetDatabaseQuery() string {
	query := url.URL{
		Scheme: "mysql",
		User:   url.UserPassword(config.DbSettings.Username, config.DbSettings.Password),
		Host:   fmt.Sprintf("tcp(%s)", config.DbSettings.Host),
		Path:   config.DbSettings.Database,
	}
	return query.String()
}
