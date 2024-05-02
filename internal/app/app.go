package app

import (
	"database/sql"
	"sync"
	"tg_reader_bot/internal/config"
)

type Container struct {
	Config   *config.ConfigStructure
	Database *sql.DB
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
