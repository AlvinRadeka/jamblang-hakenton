package domain

import "github.com/Masterminds/squirrel"

type PaginationQuery struct {
	Page  int64
	Limit int64
}

type GenericResponse struct {
	Success bool `json:"success"`
}

func (pg PaginationQuery) generatePaginationQuery(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
	if pg.Page < 1 {
		pg.Page = 1
	}

	if pg.Limit < 1 {
		pg.Limit = 10
	}

	offset := pg.Limit * (pg.Page - 1)

	sb = sb.Limit(uint64(pg.Limit)).Offset(uint64(offset))
	return sb
}
