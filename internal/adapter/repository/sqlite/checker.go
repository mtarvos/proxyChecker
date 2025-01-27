package sqlite

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"proxyChecker/internal/lib/helpers"
)

func (s *Storage) SetAlive(proxy string, port int, ip string, alive bool) error {
	const fn = "sqlite.SetAlive"

	queryBuilder := squirrel.Update("proxy").
		Set("alive", helpers.BoolToInt(alive)).
		Set("out_ip", ip).
		Where(squirrel.Eq{"proxy": proxy}).
		Where(squirrel.Eq{"port": port})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%s error build sql query: %s", fn, err.Error())
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s exec delete: %s", fn, err.Error())
	}

	return nil
}
