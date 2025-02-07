package entity

import (
	"database/sql"
	"encoding/json"
)

type CustomNullString struct {
	sql.NullString
}

func (cns *CustomNullString) MarshalJSON() ([]byte, error) {
	if !cns.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(cns.String)
}

type CustomNullInt32 struct {
	sql.NullInt32
}

func (cni *CustomNullInt32) MarshalJSON() ([]byte, error) {
	if !cni.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(cni.Int32)
}
