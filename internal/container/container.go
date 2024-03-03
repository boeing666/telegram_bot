package container

import (
	"database/sql"
	"sync"
	"tg_reader_bot/internal/commands"
	"tg_reader_bot/internal/config"
)

type Container struct {
	Config   *config.ConfigStructure
	Database *sql.DB
	Handler  *commands.Handler
}

var (
	container *Container
	once      sync.Once
)

func (c *Container) Init(config *config.ConfigStructure, database *sql.DB, handler *commands.Handler) {
	container.Config = config
	container.Database = database
	container.Handler = handler
}

func GetContainer() *Container {
	once.Do(func() {
		container = &Container{}
	})
	return container
}
