package sqlite

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"proxyChecker/internal/entity"
)

func (s *Storage) GetCountByFilter(ctx context.Context, filter entity.Filters) (int, error) {
	const fn = "GetCountByFilter"

	queryBuilder := squirrel.Select("count(id)").From("proxy")

	queryBuilder = setWhereBlock(queryBuilder, filter)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return -1, fmt.Errorf("%s error build sql query: %w", fn, err)
	}

	finalArgs := args
	if finalArgs == nil {
		finalArgs = []interface{}{}
	}

	var count int
	err = s.db.QueryRowContext(ctx, query, finalArgs...).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("%s query: %w", fn, err)
	}

	return count, nil
}

func (s *Storage) GetDistinctField(ctx context.Context, fieldName string, filter entity.Filters) ([]string, error) {
	const fn = "Storage.GetDistinctCountry"

	queryBuilder := squirrel.Select("distinct " + fieldName).From("proxy").Where(squirrel.NotEq{fieldName: ""})

	queryBuilder = setWhereBlock(queryBuilder, filter)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s error build sql query: %w", fn, err)
	}

	finalArgs := args
	if finalArgs == nil {
		finalArgs = []interface{}{}
	}

	rows, err := s.db.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, fmt.Errorf("%s query: %w", fn, err)
	}

	var fieldList []string

	for rows.Next() {
		var field string
		if err = rows.Scan(&field); err != nil {
			return nil, fmt.Errorf("%s scan: %w", fn, err)
		}
		fieldList = append(fieldList, field)
	}

	return fieldList, nil
}
