package repository

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
)

func (wr *skuRepository) Get(skuID int64) (domain.SKU, error) {
	var (
		skuData domain.SKU
	)

	query, args, err := squirrel.Select(
		"id",
		"sku",
		"wh_code",
		"bin_code",
		"zone_id",
		"name",
		"created_at",
		"updated_at",
	).From("skus").Where(
		squirrel.Eq{"id": skuID},
	).ToSql()

	if err != nil {
		return skuData, err
	}

	query = wr.sql.Rebind(query)
	row := wr.sql.QueryRow(query, args...)
	if err := row.Err(); err != nil {
		return skuData, err
	}

	err = row.Scan(
		&skuData.ID,
		&skuData.SKU,
		&skuData.WHCode,
		&skuData.BinCode,
		&skuData.ZoneID,
		&skuData.Name,
		&skuData.CreatedAt,
		&skuData.UpdatedAt,
	)
	if err != nil {
		return skuData, err
	}

	return skuData, nil
}

func (wr *skuRepository) Select(params domain.SKUQueryParameter) ([]domain.SKU, error) {
	var (
		skusData []domain.SKU
	)

	selector := squirrel.Select(
		"id",
		"sku",
		"wh_code",
		"bin_code",
		"zone_id",
		"name",
		"created_at",
		"updated_at",
	).From("skus")
	selector = params.BuildSQLQuery(selector)
	query, args, err := selector.ToSql()

	if err != nil {
		return skusData, err
	}

	query = wr.sql.Rebind(query)
	rows, err := wr.sql.Query(query, args...)
	if err != nil {
		return skusData, err
	}

	for rows.Next() {
		var skuData domain.SKU
		if err := rows.Scan(
			&skuData.ID,
			&skuData.SKU,
			&skuData.WHCode,
			&skuData.BinCode,
			&skuData.ZoneID,
			&skuData.Name,
			&skuData.CreatedAt,
			&skuData.UpdatedAt,
		); err != nil {
			return skusData, err
		}

		skusData = append(skusData, skuData)
	}

	return skusData, nil
}

func (wr *skuRepository) Create(data domain.SKUDataParameter) (domain.SKU, error) {
	var (
		skuData domain.SKU
		t       = time.Now()
	)

	query, args, err := squirrel.Insert("skus").Columns(
		"sku",
		"wh_code",
		"bin_code",
		"zone_id",
		"name",
		"created_at",
		"updated_at",
	).Values(
		data.SKU,
		data.WHCode,
		data.BinCode,
		data.ZoneID,
		data.Name,
		t, t,
	).ToSql()

	if err != nil {
		wr.logger.Errorln(err)
		return skuData, err
	}

	query = wr.sql.Rebind(query)
	result, err := wr.sql.Exec(query, args...)
	if err != nil {
		wr.logger.Errorln(err)
		return skuData, err
	}

	lastInserted, err := result.LastInsertId()
	if err != nil {
		wr.logger.Errorln(err)
		return skuData, err
	}

	skuData, err = wr.Get(lastInserted)
	if err != nil {
		wr.logger.Errorln(err)
		return skuData, err
	}

	return skuData, nil
}

func (wr *skuRepository) Update(skuID int64, data domain.SKUDataParameter) (domain.SKU, error) {
	var (
		skuData domain.SKU
	)

	query, args, err := squirrel.Update("skus").
		Set("name", data.Name).
		Set("sku", data.SKU).
		Set("wh_code", data.WHCode).
		Set("bin_code", data.BinCode).
		Set("zone_id", data.ZoneID).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": skuID}).
		ToSql()
	if err != nil {
		return skuData, err
	}

	query = wr.sql.Rebind(query)
	_, err = wr.sql.Exec(query, args...)
	if err != nil {
		return skuData, err
	}

	skuData, err = wr.Get(skuID)
	if err != nil {
		return skuData, err
	}

	return skuData, nil
}

func (wr *skuRepository) Delete(skuID int64) error {
	query, args, err := squirrel.Delete("skus").Where(squirrel.Eq{"id": skuID}).ToSql()
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
