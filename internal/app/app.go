package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/config"
	dbModel "github.com/hiteshpattanayak-tw/ports_processor/internal/app/db/migrations/model"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/model"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/repository"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db/migrations"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/json_processor"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"
	"strings"
)

type App struct {
	Ctx        context.Context
	Cancel     context.CancelFunc
	JsonStream json_processor.Stream
	PortRepo   repository.PortRepository
	Logger     *logger.ZapLogger
	AppConfig  config.AppConfig
	Migrator   *migrations.Migrator
}

func (a *App) Start(_ func(err error)) {
	a.Migrator.RunMigrations()
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
	b, err := json.Marshal(data)
	if err != nil {
		a.Logger.With(a.Ctx).Warn(fmt.Sprintf("error while seriallizing port data: %v", err.Error()))
		return
	}

	var port *model.Port
	err = json.Unmarshal(b, &port)
	if err != nil {
		a.Logger.With(a.Ctx).Warn(fmt.Sprintf("error while deseriallizing port data: %v", err.Error()))
		return
	}

	dbPort := convertToDbModel(port)

	err = a.PortRepo.UpsertPort(a.Ctx, a.AppConfig.GetDatabaseConfig().PortsTableName, dbPort)
	if err != nil {
		a.Logger.With(a.Ctx).Error(fmt.Sprintf("error while upserting port details: %v", err.Error()))
		return
	}
}

func convertToDbModel(port *model.Port) *dbModel.Port {
	return &dbModel.Port{
		Name:        port.Name,
		City:        port.City,
		Country:     port.Country,
		Alias:       strings.Join(port.Alias, ","),
		Regions:     strings.Join(port.Regions, ","),
		Coordinates: concatCoordinates(port.Coordinates),
		Province:    port.Province,
		Timezone:    port.Timezone,
		Unlocs:      strings.Join(port.Unlocs, ","),
		Code:        port.Code,
	}
}

func concatCoordinates(coord []float64) string {
	coordinates := ""
	for _, c := range coord {
		coordinates = fmt.Sprintf("%s,%f", coordinates, c)
	}

	return coordinates
}
