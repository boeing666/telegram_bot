package container

import (
	"database/sql"
	"sync"
	"tg_reader_bot/internal/config"
)

type Container struct {
	config   *config.ConfigStructure
	database *sql.DB
}

var (
	container *Container
	once      sync.Once
)

func (c *Container) Init(config *config.ConfigStructure, database *sql.DB) {
	container.config = config
	container.database = database
}

func GetContainer() *Container {
	once.Do(func() {
		container = &Container{}
	})
	return container
}
