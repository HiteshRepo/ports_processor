//go:build wireinject
// +build wireinject

//go:generate wire

package di

import (
	"context"
	"github.com/google/wire"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/config"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/repository"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db/migrations"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/json_processor"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"
)

var configSet = wire.NewSet(
	config.ProvideAppConfig,
	wire.FieldsOf(new(config.AppConfig), "DatabaseConfig"),
	wire.FieldsOf(new(config.AppConfig), "ServerConfig"),
	wire.FieldsOf(new(config.ServerConfig), "ServiceName"),
	wire.FieldsOf(new(config.ServerConfig), "LogLevel"),
)

var logSet = wire.NewSet(
	logger.ProvideLogger,
	wire.Bind(new(logger.Logger), new(*logger.ZapLogger)),
)

var dbSet = wire.NewSet(
	db.ProvideDatabase,
	migrations.ProvideMigrator,
)

var repoSet = wire.NewSet(
	repository.ProvidePortRepository,
)

var pkgSet = wire.NewSet(
	json_processor.ProvideJSONStream,
)

func InitializeApp(ctx context.Context, cancel context.CancelFunc) (*app.App, error) {
	wire.Build(
		configSet,
		pkgSet,
		logSet,
		dbSet,
		repoSet,
		wire.Struct(new(app.App), "*"),
	)

	return &app.App{}, nil
}