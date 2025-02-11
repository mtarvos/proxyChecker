package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mattn/go-sqlite3"
	"proxyChecker/internal/entity"
)

func (s *Storage) GetProxy(ctx context.Context, filter entity.Filters) ([]entity.ProxyItem, error) {
	const fn = "Storage.GetProxy"

	queryBuilder := squirrel.Select("id, proxy, port, out_ip, country, city, ISP, timezone, alive").From("proxy")

	if filter.Page != 0 && filter.Limit != 0 {
		queryBuilder = queryBuilder.Offset(uint64((filter.Page - 1) * filter.Limit))
	}

	if filter.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(filter.Limit))
	}

	queryBuilder = setWhereBlock(queryBuilder, filter)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s error build sql query: %w", fn, err)
	}

	var proxyList []entity.ProxyItem
	err = s.db.SelectContext(ctx, &proxyList, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s query: %w", fn, err)
	}

	return proxyList, nil
}

func (s *Storage) UpdateProxyItemByID(ctx context.Context, item entity.ProxyItem) error {
	const fn = "sqlite.UpdateProxyItemByID"

	if item.ID == 0 {
		return fmt.Errorf("%s bad item: ID is not defined ", fn)
	}

	updateBuilder := squirrel.Update("proxy")
	updateBuilder = setSetBlock(updateBuilder, item)
	updateBuilder = updateBuilder.
		Where(squirrel.Eq{"id": item.ID})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%s error build sql update query: %w", fn, err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s exec update: %w", fn, err)
	}

	return nil
}

func setSetBlock(updateBuilder squirrel.UpdateBuilder, item entity.ProxyItem) squirrel.UpdateBuilder {
	updateBuilder = updateBuilder.Set("alive", item.Alive)
	updateBuilder = updateBuilder.Set("Country", item.Country)
	updateBuilder = updateBuilder.Set("City", item.City)
	updateBuilder = updateBuilder.Set("ISP", item.ISP)
	updateBuilder = updateBuilder.Set("Timezone", item.Timezone)
	updateBuilder = updateBuilder.Set("out_ip", item.OutIP)
	return updateBuilder
}

func setWhereBlock(queryBuilder squirrel.SelectBuilder, filter entity.Filters) squirrel.SelectBuilder {
	if filter.Alive != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"alive": *filter.Alive})
	}

	if filter.Country != nil {
		queryBuilder = buildWhereBlock(queryBuilder, *filter.Country, "country")
	}

	if filter.City != nil {
		queryBuilder = buildWhereBlock(queryBuilder, *filter.City, "city")
	}

	if filter.ISP != nil {
		queryBuilder = buildWhereBlock(queryBuilder, *filter.ISP, "ISP")
	}

	if filter.OutIP != nil {
		queryBuilder = buildWhereBlock(queryBuilder, *filter.OutIP, "out_ip")
	}

	return queryBuilder
}

func buildWhereBlock(queryBuilder squirrel.SelectBuilder, filterField entity.StringFilter, fieldName string) squirrel.SelectBuilder {
	if val, ok := filterField.Val.(sql.NullString); ok && !val.Valid {
		switch filterField.Op {
		case entity.Eq:
			queryBuilder = queryBuilder.Where(squirrel.Eq{fieldName: filterField.Val})
		case entity.Ne:
			queryBuilder = queryBuilder.Where(squirrel.NotEq{fieldName: filterField.Val})
		}
	} else {
		switch filterField.Op {
		case entity.Eq:
			queryBuilder = queryBuilder.Where(squirrel.Expr("LOWER("+fieldName+") = LOWER(?)", filterField.Val))
		case entity.Ne:
			queryBuilder = queryBuilder.Where(squirrel.Expr("LOWER("+fieldName+") <> LOWER(?)", filterField.Val))
		}
	}

	return queryBuilder
}

func (s *Storage) SaveProxy(ctx context.Context, proxyList []entity.ProxyItem) error {
	const fn = "sqlite.SaveProxy"

	query, err := s.db.Prepare("INSERT INTO PROXY(proxy, port, country, city, ISP, timezone, alive) VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("%s prepare insert: %s", fn, err)
	}

	for _, proxy := range proxyList {
		_, err := query.ExecContext(ctx, proxy.IP, proxy.Port, proxy.Country, proxy.City, proxy.ISP, proxy.Timezone, proxy.Alive)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode != sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s exec insert: %s", fn, err)
			}
		}
	}

	return nil
}

func (s *Storage) DeleteProxy(proxyList []entity.ProxyItem) error {
	const fn = "sqlite.DeleteProxy"

	query, err := s.db.Prepare("DELETE FROM PROXY WHERE proxy = '?';")
	if err != nil {
		return fmt.Errorf("%s prepare delete: %s", fn, err)
	}

	for _, proxy := range proxyList {
		_, err := query.Exec(proxy.IP)
		if err != nil {
			return fmt.Errorf("%s exec delete: %s", fn, err)
		}
	}

	return nil
}

func (s *Storage) ClearProxy() error {
	const fn = "sqlite.GetProxyByISP"

	_, err := s.db.Exec(`DELETE FROM PROXY`)
	if err != nil {
		return fmt.Errorf("%s query: %s", fn, err)
	}

	return nil
}
