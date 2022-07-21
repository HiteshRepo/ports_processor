package db

import (
	"fmt"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

func ProvideDatabase(config DatabaseConfig, serviceName string) (*gorm.DB, error) {
	dbName := config.Dbname
	username := config.Username
	password := config.Password
	host := config.Host
	port := config.Port
	sslMode := config.SslMode
	connectionTimeOut := config.Connection.TimeOut

	args := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s connect_timeout=%d", host, port, username, dbName, password, sslMode, connectionTimeOut)

	sqltrace.Register("pgx", &stdlib.Driver{}, sqltrace.WithServiceName(serviceName))
	sqlDB, err := sqltrace.Open("pgx", args)
	if err != nil {
		return nil, errors.WithMessage(err, "Setup.db.DB")
	}

	conn, err := gormtrace.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: fmt.Sprintf("%s.", config.Schema),
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, errors.WithMessagef(err, "GetDB.gorm.Open: failed to connect to db")
	}

	ConfigureDatabasePool(sqlDB, config.Connection)

	if err := sqlDB.Ping(); err != nil {
		return nil, errors.WithMessage(err, "Setup.sqlDB.Ping")
	}

	log.Println("Created DB connection pool",
		zap.Int("openConnections", sqlDB.Stats().OpenConnections),
		zap.Int("maxOpenConnections", sqlDB.Stats().MaxOpenConnections),
		zap.Int("connectionsInUse", sqlDB.Stats().InUse),
		zap.Int("idleConnections", sqlDB.Stats().Idle),
	)

	return conn, nil
}

type DatabasePool interface {
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetConnMaxIdleTime(d time.Duration)
}

func ConfigureDatabasePool(pool DatabasePool, connection ConnectionPool) {
	pool.SetMaxOpenConns(connection.MaxOpenConnections)
	pool.SetConnMaxLifetime(time.Second * time.Duration(connection.MaxLifeTime))
	pool.SetMaxIdleConns(connection.MaxIdleConnections)
	pool.SetConnMaxIdleTime(time.Second * time.Duration(connection.MaxIdleTime))
}
