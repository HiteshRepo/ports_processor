package db

import "time"

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
