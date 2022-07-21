package db_test

import (
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db/mocks"
	"testing"
	"time"
)

func Test_ConfigureDatabasePool(t *testing.T) {
	connectionProperties := db.ConnectionPool{
		MaxOpenConnections: 1,
		MaxIdleConnections: 1,
		MaxIdleTime:        1,
		MaxLifeTime:        1,
		TimeOut:            1,
	}

	mockDb := mocks.DatabasePool{}
	mockDb.On("SetMaxOpenConns", connectionProperties.MaxOpenConnections).Return(nil)
	mockDb.On("SetConnMaxLifetime", time.Duration(connectionProperties.MaxLifeTime)*time.Second).Return(nil)
	mockDb.On("SetMaxIdleConns", connectionProperties.MaxIdleConnections).Return(nil)
	mockDb.On("SetConnMaxIdleTime", time.Duration(connectionProperties.MaxIdleTime)*time.Second).Return(nil)

	db.ConfigureDatabasePool(&mockDb, connectionProperties)

	mockDb.AssertCalled(t, "SetMaxOpenConns", connectionProperties.MaxOpenConnections)
	mockDb.AssertCalled(t, "SetConnMaxLifetime", time.Duration(connectionProperties.MaxLifeTime)*time.Second)
	mockDb.AssertCalled(t, "SetMaxIdleConns", connectionProperties.MaxIdleConnections)
	mockDb.AssertCalled(t, "SetConnMaxIdleTime", time.Duration(connectionProperties.MaxIdleTime)*time.Second)
}
