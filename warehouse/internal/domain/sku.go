package domain

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
)

type SKU struct {
	ID        int64
	SKU       string
	Name      string
	WHCode    string
	BinCode   string
	ZoneID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (sk SKU) SKUResponse() SKUResponse {
	return SKUResponse{
		ID:        sk.ID,
		SKU:       sk.SKU,
		Name:      sk.Name,
		WHCode:    sk.WHCode,
		BinCode:   sk.BinCode,
		ZoneID:    sk.ZoneID,
		CreatedAt: sk.CreatedAt,
		UpdatedAt: sk.UpdatedAt,
	}
}

type SKUResponse struct {
	ID        int64     `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	WHCode    string    `json:"wh_code"`
	BinCode   string    `json:"bin_code"`
	ZoneID    string    `json:"zone_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SKUDataParameter struct {
	SKU     string `json:"sku" validate:"required"`
	WHCode  string `json:"wh_code" validate:"required"`
	BinCode string `json:"bin_code" validate:"required"`
	ZoneID  string `json:"zone_id" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

type SKUQueryParameter struct {
	PaginationQuery
	SKU []string
}

func (wh *SKUQueryParameter) Parse(uv url.Values) error {
	if page := uv.Get("page"); len(page) > 0 {
		i, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			return errors.New("invalid Page Parameter")
		}
		wh.Page = i
	}

	if limit := uv.Get("limit"); len(limit) > 0 {
		i, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return errors.New("invalid Limit Parameter")
		}
		wh.Limit = i
	}

	if skus := uv["sku"]; len(skus) > 0 {
		wh.SKU = append(wh.SKU, skus...)
	}

	return nil
}

func (wh SKUQueryParameter) BuildSQLQuery(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
	sb = wh.generatePaginationQuery(sb)

	if len(wh.SKU) > 0 {
		sb = sb.Where(squirrel.Eq{"sku": wh.SKU})
	}

	return sb
}

type SKURepository interface {
	Get(skuID int64) (SKU, error)
	Select(params SKUQueryParameter) ([]SKU, error)
	Create(data SKUDataParameter) (SKU, error)
	Update(skuID int64, data SKUDataParameter) (SKU, error)
	Delete(skuID int64) error
}

type SKUUsecase interface {
	Get(skuID int64) (SKUResponse, error)
	Select(params SKUQueryParameter) ([]SKUResponse, error)
	Create(data SKUDataParameter) (SKUResponse, error)
	Update(skuID int64, data SKUDataParameter) (SKUResponse, error)
	Delete(skuID int64) (GenericResponse, error)
}
