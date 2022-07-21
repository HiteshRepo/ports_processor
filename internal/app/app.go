package app

import (
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/repository"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/json_processor"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"
)

type App struct {
	JsonStream json_processor.Stream
	PortRepo   repository.PortRepository
	Logger     *logger.ZapLogger
}

func (a *App) Start() {}

func (a *App) Shutdown() {}
