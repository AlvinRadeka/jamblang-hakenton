package repository

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
)

func (wr *commodityRepository) Get(commodityID int64) (domain.Commodity, error) {
	var (
		commodityData domain.Commodity
	)

	query, args, err := squirrel.Select(
		"id",
		"name",
		"description",
		"created_at",
		"updated_at",
	).From("commodities").Where(
		squirrel.Eq{"id": commodityID},
	).ToSql()

	if err != nil {
		return commodityData, err
	}

	query = wr.sql.Rebind(query)
	row := wr.sql.QueryRow(query, args...)
	if err := row.Err(); err != nil {
		return commodityData, err
	}

	err = row.Scan(
		&commodityData.ID,
		&commodityData.Name,
		&commodityData.Description,
		&commodityData.CreatedAt,
		&commodityData.UpdatedAt,
	)
	if err != nil {
		return commodityData, err
	}

	return commodityData, nil
}

func (wr *commodityRepository) Select(params domain.CommodityQueryParameter) ([]domain.Commodity, error) {
	var (
		commoditiesData []domain.Commodity
	)

	selector := squirrel.Select(
		"id",
		"name",
		"description",
		"created_at",
		"updated_at",
	).From("commodities")
	selector = params.BuildSQLQuery(selector)
	query, args, err := selector.ToSql()

	if err != nil {
		return commoditiesData, err
	}

	query = wr.sql.Rebind(query)
	rows, err := wr.sql.Query(query, args...)
	if err != nil {
		return commoditiesData, err
	}

	for rows.Next() {
		var commodityData domain.Commodity
		if err := rows.Scan(
			&commodityData.ID,
			&commodityData.Name,
			&commodityData.Description,
			&commodityData.CreatedAt,
			&commodityData.UpdatedAt,
		); err != nil {
			return commoditiesData, err
		}

		commoditiesData = append(commoditiesData, commodityData)
	}

	return commoditiesData, nil
}

func (wr *commodityRepository) Create(data domain.CommodityDataParameter) (domain.Commodity, error) {
	var (
		commodityData domain.Commodity
		t             = time.Now()
	)

	query, args, err := squirrel.Insert("commodities").Columns(
		"name",
		"description",
		"created_at",
		"updated_at",
	).Values(
		data.Name,
		data.Description,
		t, t,
	).ToSql()

	if err != nil {
		wr.logger.Errorln(err)
		return commodityData, err
	}

	query = wr.sql.Rebind(query)
	result, err := wr.sql.Exec(query, args...)
	if err != nil {
		wr.logger.Errorln(err)
		return commodityData, err
	}

	lastInserted, err := result.LastInsertId()
	if err != nil {
		wr.logger.Errorln(err)
		return commodityData, err
	}

	commodityData, err = wr.Get(lastInserted)
	if err != nil {
		wr.logger.Errorln(err)
		return commodityData, err
	}

	return commodityData, nil
}

func (wr *commodityRepository) Update(commodityID int64, data domain.CommodityDataParameter) (domain.Commodity, error) {
	var (
		commodityData domain.Commodity
	)

	query, args, err := squirrel.Update("commodities").
		Set("name", data.Name).
		Set("description", data.Description).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": commodityID}).
		ToSql()
	if err != nil {
		return commodityData, err
	}

	query = wr.sql.Rebind(query)
	_, err = wr.sql.Exec(query, args...)
	if err != nil {
		return commodityData, err
	}

	commodityData, err = wr.Get(commodityID)
	if err != nil {
		return commodityData, err
	}

	return commodityData, nil
}

func (wr *commodityRepository) Delete(commodityID int64) error {
	query, args, err := squirrel.Delete("commodities").Where(squirrel.Eq{"id": commodityID}).ToSql()
	if err != nil {
		return err
	}

	query = wr.sql.Rebind(query)
	_, err = wr.sql.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
