package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mattn/go-sqlite3"
	"log/slog"
	"proxyChecker/internal/entity"
)

func New(storagePath string, log *slog.Logger) (*Storage, error) {
	const fn = "sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s opening sqlite db: %w", fn, err)
	}

	storage := &Storage{db: db, log: log}

	err = storage.MigrationsUP()
	if err != nil {
		return nil, fmt.Errorf("%s MigrationsUP failed: %e", fn, err)
	}

	return storage, nil
}

func (s *Storage) Get(filter entity.Filters) ([]entity.ProxyItem, error) {
	const fn = "Storage.Get"

	queryBuilder := squirrel.Select("proxy, port, country, ISP, timezone, alive, status").From("proxy")

	queryBuilder = setWhereByFilter(queryBuilder, filter)

	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s error build sql query: %s", fn, err.Error())
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s query: %s", fn, err)
	}

	var proxyList []entity.ProxyItem

	for rows.Next() {
		var proxy entity.ProxyItem
		if err := rows.Scan(&proxy.Ip, &proxy.Port, &proxy.Country, &proxy.ISP, &proxy.Timezone, &proxy.Alive, &proxy.Status); err != nil {
			return nil, fmt.Errorf("%s scan: %s", fn, err)
		}
		proxyList = append(proxyList, proxy)
	}

	return proxyList, nil
}

func (s *Storage) ProxyUpdate(items []entity.ProxyItem) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetCountByFilter(filter entity.Filters) (int, error) {
	const fn = "GetCountByFilter"

	queryBuilder := squirrel.Select("count(id)").From("proxy")

	queryBuilder = setWhereByFilter(queryBuilder, filter)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return -1, fmt.Errorf("%s error build sql query: %s", fn, err.Error())
	}

	finalArgs := args
	if finalArgs == nil {
		finalArgs = []interface{}{}
	}

	var count int
	err = s.db.QueryRow(query, finalArgs...).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("%s query: %s", fn, err.Error())
	}

	return count, nil
}

func (s *Storage) GetDistinctField(fieldName string, filter entity.Filters) ([]string, error) {
	const fn = "Storage.GetDistinctCountry"

	queryBuilder := squirrel.Select("distinct " + fieldName).From("proxy").Where(squirrel.NotEq{fieldName: ""})

	queryBuilder = setWhereByFilter(queryBuilder, filter)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s error build sql query: %s", fn, err.Error())
	}

	finalArgs := args
	if finalArgs == nil {
		finalArgs = []interface{}{}
	}

	rows, err := s.db.Query(query, finalArgs...)
	if err != nil {
		return nil, fmt.Errorf("%s query: %s", fn, err)
	}

	var fieldList []string

	for rows.Next() {
		var field string
		if err := rows.Scan(&field); err != nil {
			return nil, fmt.Errorf("%s scan: %s", fn, err)
		}
		fieldList = append(fieldList, field)
	}

	return fieldList, nil
}

func setWhereByFilter(queryBuilder squirrel.SelectBuilder, filter entity.Filters) squirrel.SelectBuilder {
	if filter.AliveOnly != nil {
		alive := 0
		if *filter.AliveOnly {
			alive = 1
		}
		queryBuilder = queryBuilder.Where(squirrel.Eq{"alive": alive})
	}

	if filter.Country != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"country": filter.Country})
	}

	if filter.ISP != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"ISP": filter.Country})
	}

	return queryBuilder
}

func (s *Storage) SaveProxy(proxyList []entity.ProxyItem) error {
	const fn = "sqlite.SaveProxy"

	query, err := s.db.Prepare("INSERT INTO PROXY(proxy, port, country, ISP, timezone, alive, status) VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("%s prepare insert: %s", fn, err)
	}

	for _, proxy := range proxyList {
		_, err := query.Exec(proxy.Ip, proxy.Port, proxy.Country, proxy.ISP, proxy.Timezone, proxy.Alive, proxy.Status)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode != sqlite3.ErrConstraintUnique {
				return fmt.Errorf("%s exec insert: %s", fn, err)
			}
		}
	}

	return nil
}

func (s *Storage) RemoveProxy(proxyList []entity.ProxyItem) error {
	const fn = "sqlite.RemoveProxy"

	query, err := s.db.Prepare("DELETE FROM PROXY WHERE proxy = '?';")
	if err != nil {
		return fmt.Errorf("%s prepare delete: %s", fn, err)
	}

	for _, proxy := range proxyList {
		_, err := query.Exec(proxy.Ip)
		if err != nil {
			return fmt.Errorf("%s exec delete: %s", fn, err)
		}
	}

	return nil
}

func (s *Storage) GetAll() ([]entity.ProxyItem, error) {
	const fn = "sqlite.GetAll"

	rows, err := s.db.Query(`SELECT proxy, port, country, ISP, timezone, alive, status FROM PROXY;`)
	if err != nil {
		return nil, fmt.Errorf("%s query: %s", fn, err)
	}

	var proxyList []entity.ProxyItem

	for rows.Next() {
		var proxy entity.ProxyItem
		if err := rows.Scan(&proxy.Ip, &proxy.Port, &proxy.Country, &proxy.ISP, &proxy.Timezone, &proxy.Alive, &proxy.Status); err != nil {
			return nil, fmt.Errorf("%s scan: %s", fn, err)
		}
		proxyList = append(proxyList, proxy)
	}

	return proxyList, nil
}

func (s *Storage) GetProxyByAlive(alive bool) ([]entity.ProxyItem, error) {
	const fn = "sqlite.GetProxyByAlive"

	aliveInt := 0
	if alive {
		aliveInt = 1
	}

	rows, err := s.db.Query(`SELECT proxy, port, country, ISP, timezone, alive, status FROM PROXY where alive = ?`, aliveInt)
	if err != nil {
		return nil, fmt.Errorf("%s query: %s", fn, err)
	}

	var proxyList []entity.ProxyItem

	for rows.Next() {
		var proxy entity.ProxyItem
		if err := rows.Scan(&proxy.Ip, &proxy.Port, &proxy.Country, &proxy.ISP, &proxy.Timezone, &proxy.Alive, &proxy.Status); err != nil {
			return nil, fmt.Errorf("%s scan: %s", fn, err)
		}
		proxyList = append(proxyList, proxy)
	}

	return proxyList, nil
}

func (s *Storage) GetProxyByCountry(country string) ([]entity.ProxyItem, error) {
	const fn = "sqlite.GetProxyByCountry"

	rows, err := s.db.Query(`SELECT proxy, port, country, ISP, timezone, alive, status FROM PROXY where country = ?`, country)
	if err != nil {
		return nil, fmt.Errorf("%s query: %s", fn, err)
	}

	var proxyList []entity.ProxyItem

	for rows.Next() {
		var proxy entity.ProxyItem
		if err := rows.Scan(&proxy.Ip, &proxy.Port, &proxy.Country, &proxy.ISP, &proxy.Timezone, &proxy.Alive, &proxy.Status); err != nil {
			return nil, fmt.Errorf("%s scan: %s", fn, err)
		}
		proxyList = append(proxyList, proxy)
	}

	return proxyList, nil
}

func (s *Storage) GetProxyByISP(ISP string) ([]entity.ProxyItem, error) {
	const fn = "sqlite.GetProxyByISP"

	rows, err := s.db.Query(`SELECT proxy, port, country, ISP, timezone, alive, status FROM PROXY where ISP = ?`, ISP)
	if err != nil {
		return nil, fmt.Errorf("%s query: %s", fn, err)
	}

	var proxyList []entity.ProxyItem

	for rows.Next() {
		var proxy entity.ProxyItem
		if err := rows.Scan(&proxy.Ip, &proxy.Port, &proxy.Country, &proxy.ISP, &proxy.Timezone, &proxy.Alive, &proxy.Status); err != nil {
			return nil, fmt.Errorf("%s scan: %s", fn, err)
		}
		proxyList = append(proxyList, proxy)
	}

	return proxyList, nil
}

func (s *Storage) ClearProxy() error {
	const fn = "sqlite.GetProxyByISP"

	_, err := s.db.Exec(`DELETE FROM PROXY`)
	if err != nil {
		return fmt.Errorf("%s query: %s", fn, err)
	}

	return nil
}
