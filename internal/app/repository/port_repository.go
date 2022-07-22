package repository

//go:generate mockery --name=PortRepository

import (
	"context"
	dbModel "github.com/hiteshpattanayak-tw/ports_processor/internal/app/db/migrations/model"
	"gorm.io/gorm"
)

type PortRepository interface {
	UpsertPort(ctx context.Context, tblName string, port *dbModel.Port) error
}

type portRepository struct {
	portDb *gorm.DB
}

func ProvidePortRepository(db *gorm.DB) PortRepository {
	return &portRepository{portDb: db}
}

func (pr *portRepository) UpsertPort(ctx context.Context, tblName string, port *dbModel.Port) error {
	err := pr.portDb.WithContext(ctx).Table(tblName).Where(dbModel.Port{Name: port.Name}).Save(port).Error
	if err != nil {
		return err
	}

	return nil
}