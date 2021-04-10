package usecase

import (
	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
	"github.com/sirupsen/logrus"
)

type commodityUsecase struct {
	logger    *logrus.Logger
	commodity domain.CommodityRepository
}

func NewUsecase(logger *logrus.Logger, commodity domain.CommodityRepository) domain.CommodityUsecase {
	return &commodityUsecase{
		logger:    logger,
		commodity: commodity,
	}
}

func (uc *commodityUsecase) Get(commodityID int64) (domain.CommodityResponse, error) {
	var (
		commodityResponse domain.CommodityResponse
	)

	warehouseData, err := uc.commodity.Get(commodityID)
	if err != nil {
		return commodityResponse, err
	}

	commodityResponse = warehouseData.CommodityResponse()
	return commodityResponse, nil
}

func (uc *commodityUsecase) Select(params domain.CommodityQueryParameter) ([]domain.CommodityResponse, error) {
	var (
		commodityResponses = []domain.CommodityResponse{}
	)

	commoditysData, err := uc.commodity.Select(params)
	if err != nil {
		return commodityResponses, err
	}

	for _, commodity := range commoditysData {
		commodityResponses = append(commodityResponses, commodity.CommodityResponse())
	}

	return commodityResponses, nil
}

func (uc *commodityUsecase) Create(data domain.CommodityDataParameter) (domain.CommodityResponse, error) {
	var (
		commodityResponse domain.CommodityResponse
	)

	commodityData, err := uc.commodity.Create(data)
	if err != nil {
		return commodityResponse, err
	}

	commodityResponse = commodityData.CommodityResponse()
	return commodityResponse, nil
}

func (uc *commodityUsecase) Update(commodityID int64, data domain.CommodityDataParameter) (domain.CommodityResponse, error) {
	var (
		commodityResponse domain.CommodityResponse
	)

	commodityData, err := uc.commodity.Update(commodityID, data)
	if err != nil {
		return commodityResponse, err
	}

	commodityResponse = commodityData.CommodityResponse()
	return commodityResponse, nil
}

func (uc *commodityUsecase) Delete(commodityID int64) (domain.GenericResponse, error) {
	err := uc.commodity.Delete(commodityID)
	if err != nil {
		return domain.GenericResponse{}, err
	}

	return domain.GenericResponse{
		Success: true,
	}, nil
}
