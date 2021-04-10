package domain

import "os"

type BarcodeLambdaResponse struct {
	Data []BarcodeLambda `json:"data"`
}

type BarcodeLambda struct {
	DetectedText string          `json:"DetectedText"`
	Type         string          `json:"Type"`
	ID           int64           `json:"ID"`
	Confidence   float64         `json:"confidence"`
	Geometry     BarcodeGeometry `json:"Geometry"`
}

type BarcodeGeometry struct {
	BoundingBox BarcodeGeometryBoundingBox `json:"BoundingBox"`
	Polygon     []BarcodeGeometryPolygon   `json:"Polygon"`
}

type BarcodeGeometryBoundingBox struct {
	Width  float64 `json:"Width"`
	Height float64 `json:"Height"`
	Left   float64 `json:"Left"`
	Top    float64 `json:"Top"`
}

type BarcodeGeometryPolygon struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}

type WarehouseBarcodeResponse struct {
	Barcodes []WarehouseBarcode `json:"Barcodes"`
}

type WarehouseBarcode struct {
	SKU      string          `json:"SKU"`
	Geometry BarcodeGeometry `json:"Geometry"`
	BinCode  string          `json:"BinCode,omitempty"`
	Zone     string          `json:"Zone,omitempty"`
	Error    string          `json:"Error,omitempty"`
}

type BarcodeRepository interface {
	ParseToLambda(file64 string) (BarcodeLambdaResponse, error)
}

type BarcodeRepositoryConfig struct {
	LambdaURL string
}

type BarcodeUsecase interface {
	ParseBarcodeFromFileToLambda(file *os.File) ([]WarehouseBarcode, error)
}
