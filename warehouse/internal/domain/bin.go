package domain

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
)

type Bin struct {
	ID          int64
	WarehouseID int64
	Name        string
	Latitude    float64
	Longitude   float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (b Bin) BinResponse() BinResponse {
	return BinResponse{
		ID:          b.ID,
		WarehouseID: b.WarehouseID,
		Name:        b.Name,
		Latitude:    b.Latitude,
		Longitude:   b.Longitude,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

type BinResponse struct {
	ID          int64     `json:"id"`
	WarehouseID int64     `json:"warehouse_id"`
	Name        string    `json:"name"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BinDataParameter struct {
	WarehouseID int64   `json:"warehouse_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Latitude    float64 `json:"latitude" validate:"required,latitude"`
	Longitude   float64 `json:"longitude" validate:"required,longitude"`
}

type BinQueryParameter struct {
	PaginationQuery
	ID          []int64
	WarehouseID []int64
}

func (wh *BinQueryParameter) Parse(uv url.Values) error {
	if page := uv.Get("page"); len(page) > 0 {
		i, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			return errors.New("Invalid Page Parameter")
		}
		wh.Page = i
	}

	if limit := uv.Get("limit"); len(limit) > 0 {
		i, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return errors.New("Invalid Limit Parameter")
		}
		wh.Limit = i
	}

	if uid := uv["id"]; len(uid) > 0 {
		for _, _uid := range uid {
			i, err := strconv.ParseInt(_uid, 10, 64)
			if err != nil {
				return errors.New("Invalid ID Parameter")
			}

			wh.ID = append(wh.ID, i)
		}
	}

	if whID := uv["warehouse_id"]; len(whID) > 0 {
		for _, whID := range whID {
			i, err := strconv.ParseInt(whID, 10, 64)
			if err != nil {
				return errors.New("Invalid Warehouse ID Parameter")
			}

			wh.WarehouseID = append(wh.WarehouseID, i)
		}
	}

	return nil
}

func (wh BinQueryParameter) BuildSQLQuery(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
	sb = wh.generatePaginationQuery(sb)

	if len(wh.ID) > 0 {
		sb = sb.Where(squirrel.Eq{"id": wh.ID})
	}

	if len(wh.WarehouseID) > 0 {
		sb = sb.Where(squirrel.Eq{"warehouse_id": wh.WarehouseID})
	}

	return sb
}

type BinRepository interface {
	Get(binID int64) (Bin, error)
	Select(params BinQueryParameter) ([]Bin, error)
	Create(data BinDataParameter) (Bin, error)
	Update(binID int64, data BinDataParameter) (Bin, error)
	Delete(binID int64) error
}

type BinUsecase interface {
	Get(binID int64) (BinResponse, error)
	Select(params BinQueryParameter) ([]BinResponse, error)
	Create(data BinDataParameter) (BinResponse, error)
	Update(binID int64, data BinDataParameter) (BinResponse, error)
	Delete(binID int64) (GenericResponse, error)
}
