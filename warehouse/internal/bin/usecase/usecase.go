package usecase

import (
	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
	"github.com/sirupsen/logrus"
)

type binUsecase struct {
	logger    *logrus.Logger
	bin       domain.BinRepository
	warehouse domain.WarehouseRepository
}

func NewUsecase(logger *logrus.Logger, bin domain.BinRepository, warehouse domain.WarehouseRepository) domain.BinUsecase {
	return &binUsecase{
		logger:    logger,
		bin:       bin,
		warehouse: warehouse,
	}
}

func (uc *binUsecase) Get(binID int64) (domain.BinResponse, error) {
	var (
		binResponse domain.BinResponse
	)

	binData, err := uc.bin.Get(binID)
	if err != nil {
		return binResponse, err
	}

	binResponse = binData.BinResponse()
	return binResponse, nil
}

func (uc *binUsecase) Select(params domain.BinQueryParameter) ([]domain.BinResponse, error) {
	var (
		binResponses = []domain.BinResponse{}
	)

	binsData, err := uc.bin.Select(params)
	if err != nil {
		return binResponses, err
	}

	for _, bin := range binsData {
		binResponses = append(binResponses, bin.BinResponse())
	}

	return binResponses, nil
}

func (uc *binUsecase) Create(data domain.BinDataParameter) (domain.BinResponse, error) {
	var (
		binResponse domain.BinResponse
	)

	// Check if warehouse exists
	_, err := uc.warehouse.Get(data.WarehouseID)
	if err != nil {
		return binResponse, err
	}

	binData, err := uc.bin.Create(data)
	if err != nil {
		return binResponse, err
	}

	binResponse = binData.BinResponse()
	return binResponse, nil
}

func (uc *binUsecase) Update(binID int64, data domain.BinDataParameter) (domain.BinResponse, error) {
	var (
		binResponse domain.BinResponse
	)

	// Check if warehouse exists
	_, err := uc.warehouse.Get(data.WarehouseID)
	if err != nil {
		return binResponse, err
	}

	binData, err := uc.bin.Update(binID, data)
	if err != nil {
		return binResponse, err
	}

	binResponse = binData.BinResponse()
	return binResponse, nil
}

func (uc *binUsecase) Delete(binID int64) (domain.GenericResponse, error) {
	err := uc.bin.Delete(binID)
	if err != nil {
		return domain.GenericResponse{}, err
	}

	return domain.GenericResponse{
		Success: true,
	}, nil
}
