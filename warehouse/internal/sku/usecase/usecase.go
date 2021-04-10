package usecase

import (
	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
	"github.com/sirupsen/logrus"
)

type skuUsecase struct {
	logger *logrus.Logger
	sku    domain.SKURepository
}

func NewUsecase(logger *logrus.Logger, sku domain.SKURepository) domain.SKUUsecase {
	return &skuUsecase{
		logger: logger,
		sku:    sku,
	}
}

func (uc *skuUsecase) Get(skuID int64) (domain.SKUResponse, error) {
	var (
		skuResponse domain.SKUResponse
	)

	warehouseData, err := uc.sku.Get(skuID)
	if err != nil {
		return skuResponse, err
	}

	skuResponse = warehouseData.SKUResponse()
	return skuResponse, nil
}

func (uc *skuUsecase) Select(params domain.SKUQueryParameter) ([]domain.SKUResponse, error) {
	var (
		skuResponses = []domain.SKUResponse{}
	)

	skusData, err := uc.sku.Select(params)
	if err != nil {
		return skuResponses, err
	}

	for _, sku := range skusData {
		skuResponses = append(skuResponses, sku.SKUResponse())
	}

	return skuResponses, nil
}

func (uc *skuUsecase) Create(data domain.SKUDataParameter) (domain.SKUResponse, error) {
	var (
		skuResponse domain.SKUResponse
	)

	skuData, err := uc.sku.Create(data)
	if err != nil {
		return skuResponse, err
	}

	skuResponse = skuData.SKUResponse()
	return skuResponse, nil
}

func (uc *skuUsecase) Update(skuID int64, data domain.SKUDataParameter) (domain.SKUResponse, error) {
	var (
		skuResponse domain.SKUResponse
	)

	skuData, err := uc.sku.Update(skuID, data)
	if err != nil {
		return skuResponse, err
	}

	skuResponse = skuData.SKUResponse()
	return skuResponse, nil
}

func (uc *skuUsecase) Delete(skuID int64) (domain.GenericResponse, error) {
	err := uc.sku.Delete(skuID)
	if err != nil {
		return domain.GenericResponse{}, err
	}

	return domain.GenericResponse{
		Success: true,
	}, nil
}
