package sqlite

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mattn/go-sqlite3"
	"proxyChecker/internal/entity"
)

func (s *Storage) GetProxy(filter entity.Filters) ([]entity.ProxyItem, error) {
	const fn = "Storage.GetProxy"

	queryBuilder := squirrel.Select("id, proxy, port, out_ip, country, city, ISP, timezone, alive").From("proxy")

	queryBuilder = setWhereBlock(queryBuilder, filter)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s error build sql query: %s", fn, err.Error())
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s query: %s", fn, err)
	}

	var proxyList []entity.ProxyItem

	for rows.Next() {
		var proxy entity.ProxyItem
		if err := rows.Scan(
			&proxy.ID,
			&proxy.IP,
			&proxy.Port,
			&proxy.OutIP,
			&proxy.Country,
			&proxy.City,
			&proxy.ISP,
			&proxy.Timezone,
			&proxy.Alive,
		); err != nil {
			return nil, fmt.Errorf("%s scan: %s", fn, err)
		}
		proxyList = append(proxyList, proxy)
	}

	return proxyList, nil
}

func (s *Storage) UpdateProxyItemByID(item entity.ProxyItem) error {
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

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s exec update: %w", fn, err)
	}

	return nil
}

func setSetBlock(updateBuilder squirrel.UpdateBuilder, item entity.ProxyItem) squirrel.UpdateBuilder {
	if item.Alive != 0 {
		updateBuilder = updateBuilder.Set("alive", item.Alive)
	}

	if item.Country != "" {
		updateBuilder = updateBuilder.Set("Country", item.Country)
	}

	if item.Country != "" {
		updateBuilder = updateBuilder.Set("City", item.City)
	}

	if item.ISP != "" {
		updateBuilder = updateBuilder.Set("ISP", item.ISP)
	}

	if item.Timezone != -1 {
		updateBuilder = updateBuilder.Set("Timezone", item.Timezone)
	}

	if item.OutIP != "" {
		updateBuilder = updateBuilder.Set("out_ip", item.OutIP)
	}

	return updateBuilder
}

func setWhereBlock(queryBuilder squirrel.SelectBuilder, filter entity.Filters) squirrel.SelectBuilder {
	if filter.AliveOnly != nil {
		alive := 1
		if *filter.AliveOnly {
			alive = 2
		}
		queryBuilder = queryBuilder.Where(squirrel.Eq{"alive": alive})
	}

	if filter.Country != nil {
		switch filter.Country.Op {
		case entity.Eq:
			queryBuilder = queryBuilder.Where(squirrel.Eq{"country": filter.Country.Val})
		case entity.Ne:
			queryBuilder = queryBuilder.Where(squirrel.NotEq{"country": filter.Country.Val})
		}
	}

	if filter.City != nil {
		switch filter.City.Op {
		case entity.Eq:
			queryBuilder = queryBuilder.Where(squirrel.Eq{"city": filter.City.Val})
		case entity.Ne:
			queryBuilder = queryBuilder.Where(squirrel.NotEq{"City": filter.City.Val})
		}
	}

	if filter.ISP != nil {
		switch filter.ISP.Op {
		case entity.Eq:
			queryBuilder = queryBuilder.Where(squirrel.Eq{"ISP": filter.ISP.Val})
		case entity.Ne:
			queryBuilder = queryBuilder.Where(squirrel.NotEq{"ISP": filter.ISP.Val})
		}
	}

	if filter.OutIP != nil {
		switch filter.OutIP.Op {
		case entity.Eq:
			queryBuilder = queryBuilder.Where(squirrel.Eq{"out_ip": filter.OutIP.Val})
		case entity.Ne:
			queryBuilder = queryBuilder.Where(squirrel.NotEq{"out_ip": filter.OutIP.Val})
		}
	}

	return queryBuilder
}

func (s *Storage) SaveProxy(proxyList []entity.ProxyItem) error {
	const fn = "sqlite.SaveProxy"

	query, err := s.db.Prepare("INSERT INTO PROXY(proxy, port, country, city, ISP, timezone, alive) VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("%s prepare insert: %s", fn, err)
	}

	for _, proxy := range proxyList {
		_, err := query.Exec(proxy.IP, proxy.Port, proxy.Country, proxy.City, proxy.ISP, proxy.Timezone, proxy.Alive)
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
