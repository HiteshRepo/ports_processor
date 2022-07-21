package app

import (
	"context"
	"fmt"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/config"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/model"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/repository"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/json_processor"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"
)

type App struct {
	Ctx        context.Context
	Cancel     context.CancelFunc
	JsonStream json_processor.Stream
	PortRepo   repository.PortRepository
	Logger     *logger.ZapLogger
	AppConfig  config.AppConfig
}

func (a *App) Start(_ func(err error)) {
	go a.watchJsonStream()
	a.JsonStream.Start(a.AppConfig.GetServerConfig().PortsFilePath)
}

func (a *App) Shutdown(_ func(err error)) {}

func (a *App) watchJsonStream() {
	for data := range a.JsonStream.Watch() {
		if data.Error != nil {
			a.Logger.With(a.Ctx).Error(data.Error.Error())
		}
		a.handlePortData(data.Data)
	}
}

func (a *App) handlePortData(data interface{}) {
	port, ok := data.(*model.Port)
	if !ok {
		a.Logger.With(a.Ctx).Warn(fmt.Sprintf("invalid port details: %v", data))
		return
	}

	err := a.PortRepo.UpsertPort(a.Ctx, port)
	if err != nil {
		a.Logger.With(a.Ctx).Error(fmt.Sprintf("error while upserting port details: %v", err.Error()))
		return
	}
}
