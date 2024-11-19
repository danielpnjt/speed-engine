package cmd

import (
	"github.com/danielpnjt/speed-engine/internal/infrastructure/container"
	"github.com/danielpnjt/speed-engine/internal/server"
)

func Run() {
	server.StartService(container.New())
}
