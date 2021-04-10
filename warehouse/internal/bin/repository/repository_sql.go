package repository

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
)

func (wr *binRepository) Get(binID int64) (domain.Bin, error) {
	var (
		binData domain.Bin
	)

	query, args, err := squirrel.Select(
		"id",
		"warehouse_id",
		"name",
		"latitude",
		"longitude",
		"created_at",
		"updated_at",
	).From("bins").Where(
		squirrel.Eq{"id": binID},
	).ToSql()

	if err != nil {
		return binData, err
	}

	query = wr.sql.Rebind(query)
	row := wr.sql.QueryRow(query, args...)
	if err := row.Err(); err != nil {
		return binData, err
	}

	err = row.Scan(
		&binData.ID,
		&binData.WarehouseID,
		&binData.Name,
		&binData.Latitude,
		&binData.Longitude,
		&binData.CreatedAt,
		&binData.UpdatedAt,
	)
	if err != nil {
		return binData, err
	}

	return binData, nil
}

func (wr *binRepository) GetByWarehouseID(warehouseID int64) ([]domain.Bin, error) {
	var (
		binsData []domain.Bin
	)

	query, args, err := squirrel.Select(
		"id",
		"warehouse_id",
		"name",
		"latitude",
		"longitude",
		"created_at",
		"updated_at",
	).From("bins").Where(
		squirrel.Eq{"warehouse_id": warehouseID},
	).ToSql()

	if err != nil {
		return binsData, err
	}

	query = wr.sql.Rebind(query)
	row, err := wr.sql.Query(query, args...)
	if err != nil {
		return binsData, err
	}

	for row.Next() {
		var binData domain.Bin

		if err := row.Scan(
			&binData.ID,
			&binData.WarehouseID,
			&binData.Name,
			&binData.Latitude,
			&binData.Longitude,
			&binData.CreatedAt,
			&binData.UpdatedAt,
		); err != nil {
			return binsData, err
		}

		binsData = append(binsData, binData)
	}

	return binsData, nil
}

func (wr *binRepository) Select(params domain.BinQueryParameter) ([]domain.Bin, error) {
	var (
		binsData []domain.Bin
	)

	selector := squirrel.Select(
		"id",
		"warehouse_id",
		"name",
		"latitude",
		"longitude",
		"created_at",
		"updated_at",
	).From("bins")
	selector = params.BuildSQLQuery(selector)
	query, args, err := selector.ToSql()

	if err != nil {
		return binsData, err
	}

	query = wr.sql.Rebind(query)
	rows, err := wr.sql.Query(query, args...)
	if err != nil {
		return binsData, err
	}

	for rows.Next() {
		var binData domain.Bin
		if err := rows.Scan(
			&binData.ID,
			&binData.WarehouseID,
			&binData.Name,
			&binData.Latitude,
			&binData.Longitude,
			&binData.CreatedAt,
			&binData.UpdatedAt,
		); err != nil {
			return binsData, err
		}

		binsData = append(binsData, binData)
	}

	return binsData, nil
}

func (wr *binRepository) Create(data domain.BinDataParameter) (domain.Bin, error) {
	var (
		binData domain.Bin
		t       = time.Now()
	)

	query, args, err := squirrel.Insert("bins").Columns(
		"warehouse_id",
		"name",
		"latitude",
		"longitude",
		"created_at",
		"updated_at",
	).Values(
		data.WarehouseID,
		data.Name,
		data.Latitude,
		data.Longitude,
		t, t,
	).ToSql()

	if err != nil {
		wr.logger.Errorln(err)
		return binData, err
	}

	query = wr.sql.Rebind(query)
	result, err := wr.sql.Exec(query, args...)
	if err != nil {
		wr.logger.Errorln(err)
		return binData, err
	}

	lastInserted, err := result.LastInsertId()
	if err != nil {
		wr.logger.Errorln(err)
		return binData, err
	}

	binData, err = wr.Get(lastInserted)
	if err != nil {
		wr.logger.Errorln(err)
		return binData, err
	}

	return binData, nil
}

func (wr *binRepository) Update(binID int64, data domain.BinDataParameter) (domain.Bin, error) {
	var (
		binData domain.Bin
	)

	query, args, err := squirrel.Update("bins").
		Set("warehouse_id", data.WarehouseID).
		Set("name", data.Name).
		Set("latitude", data.Latitude).
		Set("longitude", data.Longitude).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": binID}).
		ToSql()
	if err != nil {
		return binData, err
	}

	query = wr.sql.Rebind(query)
	_, err = wr.sql.Exec(query, args...)
	if err != nil {
		return binData, err
	}

	binData, err = wr.Get(binID)
	if err != nil {
		return binData, err
	}

	return binData, nil
}

func (wr *binRepository) Delete(binID int64) error {
	query, args, err := squirrel.Delete("bins").Where(squirrel.Eq{"id": binID}).ToSql()
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
