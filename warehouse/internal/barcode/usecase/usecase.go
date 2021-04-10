package usecase

import (
	"encoding/base64"
	"io/ioutil"
	"os"

	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
	"github.com/sirupsen/logrus"
)

type barcodeUsecase struct {
	logger    *logrus.Logger
	barcode   domain.BarcodeRepository
	warehouse domain.WarehouseRepository
	sku       domain.SKURepository
}

var (
	zoneMap = make(map[string]string)
)

func NewUsecase(logger *logrus.Logger, barcode domain.BarcodeRepository, warehouse domain.WarehouseRepository, sku domain.SKURepository) domain.BarcodeUsecase {
	zoneMap = map[string]string{
		"1":  "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+1.jpg",
		"2":  "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+2.jpg",
		"3":  "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+3.jpg",
		"4":  "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+4.jpg",
		"5":  "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+5.jpg",
		"6":  "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+6.jpg",
		"7":  "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+7.jpg",
		"8":  "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+8.jpg",
		"9":  "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+9.jpg",
		"10": "https://ocr-zone.s3.amazonaws.com/warehouse/Floor_Plan+10.jpg",
	}

	return &barcodeUsecase{
		logger:    logger,
		barcode:   barcode,
		warehouse: warehouse,
		sku:       sku,
	}
}

func (b *barcodeUsecase) ParseBarcodeFromFileToLambda(file *os.File) ([]domain.WarehouseBarcode, error) {
	var (
		whBarcode = []domain.WarehouseBarcode{}
	)

	// Encode to base64
	readerFile, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return whBarcode, err
	}
	encodedFile := base64.StdEncoding.EncodeToString(readerFile)

	// Fetch data from Lambda
	barcodes, err := b.barcode.ParseToLambda(encodedFile)
	if err != nil {
		return whBarcode, err
	}

	// Iterate each Barcodes, and find the correct barcode
	for _, barcode := range barcodes.Data {
		skuFound, err := b.sku.Select(domain.SKUQueryParameter{
			SKU: []string{barcode.DetectedText},
			PaginationQuery: domain.PaginationQuery{
				Limit: 1,
				Page:  1,
			},
		})
		if err != nil {
			whBarcode = append(whBarcode, domain.WarehouseBarcode{
				SKU:      barcode.DetectedText,
				Geometry: barcode.Geometry,
				Error:    "An error occured when trying to find SKU " + barcode.DetectedText,
			})
			continue
		}

		if len(skuFound) < 1 {
			whBarcode = append(whBarcode, domain.WarehouseBarcode{
				SKU:      barcode.DetectedText,
				Geometry: barcode.Geometry,
				Error:    barcode.DetectedText + " not found",
			})
			continue
		}

		tempZoneMap := ""
		if v, ok := zoneMap[skuFound[0].ZoneID]; ok {
			tempZoneMap = v
		}

		whBarcode = append(whBarcode, domain.WarehouseBarcode{
			SKU:      skuFound[0].SKU,
			Geometry: barcode.Geometry,
			BinCode:  skuFound[0].BinCode,
			Zone:     tempZoneMap,
		})
	}

	return whBarcode, nil
}
