package domain

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
)

type Commodity struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (c Commodity) CommodityResponse() CommodityResponse {
	return CommodityResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type CommodityResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CommodityDataParameter struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type CommodityQueryParameter struct {
	PaginationQuery
	ID []int64
}

func (wh *CommodityQueryParameter) Parse(uv url.Values) error {
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

func (wh CommodityQueryParameter) BuildSQLQuery(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
	sb = wh.generatePaginationQuery(sb)

	if len(wh.ID) > 0 {
		sb = sb.Where(squirrel.Eq{"id": wh.ID})
	}

	return sb
}

type CommodityRepository interface {
	Get(commodityID int64) (Commodity, error)
	Select(params CommodityQueryParameter) ([]Commodity, error)
	Create(data CommodityDataParameter) (Commodity, error)
	Update(commodityID int64, data CommodityDataParameter) (Commodity, error)
	Delete(commodityID int64) error
}

type CommodityUsecase interface {
	Get(commodityID int64) (CommodityResponse, error)
	Select(params CommodityQueryParameter) ([]CommodityResponse, error)
	Create(data CommodityDataParameter) (CommodityResponse, error)
	Update(commodityID int64, data CommodityDataParameter) (CommodityResponse, error)
	Delete(commodityID int64) (GenericResponse, error)
}
