package domain

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
)

type Warehouse struct {
	ID        int64
	Name      string
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (wh Warehouse) WarehouseResponse() WarehouseResponse {
	return WarehouseResponse{
		ID:        wh.ID,
		Name:      wh.Name,
		Latitude:  wh.Latitude,
		Longitude: wh.Longitude,
		CreatedAt: wh.CreatedAt,
		UpdatedAt: wh.UpdatedAt,
	}
}

type WarehouseResponse struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	Bins      []BinResponse `json:"bins"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type WarehouseDataParameter struct {
	Name      string  `json:"name" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type WarehouseQueryParameter struct {
	PaginationQuery
	ID []int64
}

func (wh *WarehouseQueryParameter) Parse(uv url.Values) error {
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

	return nil
}

func (wh WarehouseQueryParameter) BuildSQLQuery(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
	sb = wh.generatePaginationQuery(sb)

	if len(wh.ID) > 0 {
		sb = sb.Where(squirrel.Eq{"id": wh.ID})
	}

	return sb
}

type WarehouseRepository interface {
	Get(warehouseID int64) (Warehouse, error)
	Select(params WarehouseQueryParameter) ([]Warehouse, error)
	Create(data WarehouseDataParameter) (Warehouse, error)
	Update(warehouseID int64, data WarehouseDataParameter) (Warehouse, error)
	Delete(warehouseID int64) error
}

type WarehouseUsecase interface {
	Get(warehouseID int64) (WarehouseResponse, error)
	Select(params WarehouseQueryParameter) ([]WarehouseResponse, error)
	Create(data WarehouseDataParameter) (WarehouseResponse, error)
	Update(warehouseID int64, data WarehouseDataParameter) (WarehouseResponse, error)
	Delete(warehouseID int64) (GenericResponse, error)
}
