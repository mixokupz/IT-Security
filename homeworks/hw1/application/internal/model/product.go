package model

import (
	"database/sql/driver"
	"errors"
)

type ProductName string

// чтение из БД
func (n *ProductName) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("invalid name")
	}
	*n = ProductName(s)
	return nil
}

// запись в БД
func (n ProductName) Value() (driver.Value, error) {
	return string(n), nil
}

type Product struct {
	ID   int64       `json:"id"`
	Name ProductName `json:"name"`
}
