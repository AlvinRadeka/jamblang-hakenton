package usecase

import (
	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
	"github.com/sirupsen/logrus"
)

type warehouseUsecase struct {
	logger    *logrus.Logger
	warehouse domain.WarehouseRepository
}

func NewUsecase(logger *logrus.Logger, warehouse domain.WarehouseRepository) domain.WarehouseUsecase {
	return &warehouseUsecase{
		logger:    logger,
		warehouse: warehouse,
	}
}

func (uc *warehouseUsecase) Get(warehouseID int64) (domain.WarehouseResponse, error) {
	var (
		warehouseResponse domain.WarehouseResponse
	)

	warehouseData, err := uc.warehouse.Get(warehouseID)
	if err != nil {
		return warehouseResponse, err
	}

	warehouseResponse = warehouseData.WarehouseResponse()
	return warehouseResponse, nil
}

func (uc *warehouseUsecase) Select(params domain.WarehouseQueryParameter) ([]domain.WarehouseResponse, error) {
	var (
		warehouseResponses = []domain.WarehouseResponse{}
	)

	warehousesData, err := uc.warehouse.Select(params)
	if err != nil {
		return warehouseResponses, err
	}

	for _, warehouse := range warehousesData {
		warehouseResponses = append(warehouseResponses, warehouse.WarehouseResponse())
	}

	return warehouseResponses, nil
}

func (uc *warehouseUsecase) Create(data domain.WarehouseDataParameter) (domain.WarehouseResponse, error) {
	var (
		warehouseResponse domain.WarehouseResponse
	)

	warehouseData, err := uc.warehouse.Create(data)
	if err != nil {
		return warehouseResponse, err
	}

	warehouseResponse = warehouseData.WarehouseResponse()
	return warehouseResponse, nil
}

func (uc *warehouseUsecase) Update(warehouseID int64, data domain.WarehouseDataParameter) (domain.WarehouseResponse, error) {
	var (
		warehouseResponse domain.WarehouseResponse
	)

	warehouseData, err := uc.warehouse.Update(warehouseID, data)
	if err != nil {
		return warehouseResponse, err
	}

	warehouseResponse = warehouseData.WarehouseResponse()
	return warehouseResponse, nil
}

func (uc *warehouseUsecase) Delete(warehouseID int64) (domain.GenericResponse, error) {
	err := uc.warehouse.Delete(warehouseID)
	if err != nil {
		return domain.GenericResponse{}, err
	}

	return domain.GenericResponse{
		Success: true,
	}, nil
}
