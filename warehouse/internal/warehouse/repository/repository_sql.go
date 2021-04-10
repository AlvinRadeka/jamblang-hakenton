package repository

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
)

func (wr *warehouseRepository) Get(warehouseID int64) (domain.Warehouse, error) {
	var (
		warehouseData domain.Warehouse
	)

	query, args, err := squirrel.Select(
		"id",
		"name",
		"latitude",
		"longitude",
		"created_at",
		"updated_at",
	).From("warehouses").Where(
		squirrel.Eq{"id": warehouseID},
	).ToSql()

	if err != nil {
		return warehouseData, err
	}

	query = wr.sql.Rebind(query)
	row := wr.sql.QueryRow(query, args...)
	if err := row.Err(); err != nil {
		return warehouseData, err
	}

	err = row.Scan(
		&warehouseData.ID,
		&warehouseData.Name,
		&warehouseData.Latitude,
		&warehouseData.Longitude,
		&warehouseData.CreatedAt,
		&warehouseData.UpdatedAt,
	)
	if err != nil {
		return warehouseData, err
	}

	return warehouseData, nil
}

func (wr *warehouseRepository) Select(params domain.WarehouseQueryParameter) ([]domain.Warehouse, error) {
	var (
		warehousesData []domain.Warehouse
	)

	selector := squirrel.Select(
		"id",
		"name",
		"latitude",
		"longitude",
		"created_at",
		"updated_at",
	).From("warehouses")
	selector = params.BuildSQLQuery(selector)
	query, args, err := selector.ToSql()

	if err != nil {
		return warehousesData, err
	}

	query = wr.sql.Rebind(query)
	rows, err := wr.sql.Query(query, args...)
	if err != nil {
		return warehousesData, err
	}

	for rows.Next() {
		var warehouseData domain.Warehouse
		if err := rows.Scan(
			&warehouseData.ID,
			&warehouseData.Name,
			&warehouseData.Latitude,
			&warehouseData.Longitude,
			&warehouseData.CreatedAt,
			&warehouseData.UpdatedAt,
		); err != nil {
			return warehousesData, err
		}

		warehousesData = append(warehousesData, warehouseData)
	}

	return warehousesData, nil
}

func (wr *warehouseRepository) Create(data domain.WarehouseDataParameter) (domain.Warehouse, error) {
	var (
		warehouseData domain.Warehouse
		t             = time.Now()
	)

	query, args, err := squirrel.Insert("warehouses").Columns(
		"name",
		"latitude",
		"longitude",
		"created_at",
		"updated_at",
	).Values(
		data.Name,
		data.Latitude,
		data.Longitude,
		t, t,
	).ToSql()

	if err != nil {
		return warehouseData, err
	}

	query = wr.sql.Rebind(query)
	result, err := wr.sql.Exec(query, args...)
	if err != nil {
		return warehouseData, err
	}

	lastInserted, err := result.LastInsertId()
	if err != nil {
		return warehouseData, err
	}

	warehouseData, err = wr.Get(lastInserted)
	if err != nil {
		return warehouseData, err
	}

	return warehouseData, nil
}

func (wr *warehouseRepository) Update(warehouseID int64, data domain.WarehouseDataParameter) (domain.Warehouse, error) {
	var (
		warehouseData domain.Warehouse
	)

	query, args, err := squirrel.Update("warehouses").
		Set("name", data.Name).
		Set("latitude", data.Latitude).
		Set("longitude", data.Longitude).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": warehouseID}).
		ToSql()
	if err != nil {
		return warehouseData, err
	}

	query = wr.sql.Rebind(query)
	_, err = wr.sql.Exec(query, args...)
	if err != nil {
		return warehouseData, err
	}

	warehouseData, err = wr.Get(warehouseID)
	if err != nil {
		return warehouseData, err
	}

	return warehouseData, nil
}

func (wr *warehouseRepository) Delete(warehouseID int64) error {
	query, args, err := squirrel.Delete("warehouses").Where(squirrel.Eq{"id": warehouseID}).ToSql()
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
