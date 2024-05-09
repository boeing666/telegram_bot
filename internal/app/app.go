package app

import (
	"database/sql"
	"sync"
	"tg_reader_bot/internal/config"
	"tg_reader_bot/internal/telegram"
)

type Container struct {
	Config   *config.ConfigStructure
	Database *sql.DB
	Client   *telegram.TGClient
}

var (
	container *Container
	once      sync.Once
)

func (c *Container) Init(config *config.ConfigStructure, database *sql.DB) {
	container.Config = config
	container.Database = database
}

func GetContainer() *Container {
	once.Do(func() {
		container = &Container{}
	})
	return container
}

func GetDatabase() *sql.DB {
	return GetContainer().Database
}

func GetConfig() *config.ConfigStructure {
	return GetContainer().Config
}

func GetClient() *telegram.TGClient {
	return GetContainer().Client
}
