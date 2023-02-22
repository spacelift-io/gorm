package gorm

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

/*
Copy-pasta of: https://github.com/go-gorm/gorm/blob/04cbd956ebed5fec1b61a819a3f7494c00d276b3/soft_delete.go
*/

var deletedAtType reflect.Type = reflect.TypeOf(DeletedAt{})

type DeletedAt sql.NullTime

// Scan implements the Scanner interface.
func (n *DeletedAt) Scan(value interface{}) error {
	return (*sql.NullTime)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n DeletedAt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n DeletedAt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *DeletedAt) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Time)
	if err == nil {
		n.Valid = true
	}
	return err
}
