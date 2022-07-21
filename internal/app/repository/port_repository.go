package repository

import (
	"context"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/app/model"
	"gorm.io/gorm"
)

type PortRepository interface {
	UpsertPort(ctx context.Context, port *model.Port) error
}

type portRepository struct {
	portDb *gorm.DB
}

func ProvideTierConfigRepository(db *gorm.DB) PortRepository {
	return &portRepository{portDb: db}
}
