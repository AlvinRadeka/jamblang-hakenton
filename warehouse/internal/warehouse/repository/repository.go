package repository

import (
	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type warehouseRepository struct {
	logger *logrus.Logger
	sql    *sqlx.DB
}

func NewSQL(logger *logrus.Logger, sql *sqlx.DB) domain.WarehouseRepository {
	return &warehouseRepository{
		logger: logger,
		sql:    sql,
	}
}
