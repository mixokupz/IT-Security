package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nsu.secure/internal/model"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{DB: db}
}

var ErrNotFound = errors.New("order not found")

// FindByID — найти заказ по id.
// Возвращает ErrNotFound, если заказа нет.
func (r *Database) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	const query = `SELECT id, Name FROM Products WHERE id = $1`

	var o model.Product
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&o.ID, &o.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find Product by id: %w", err)
	}
	return &o, nil
}

// Save — создать новый заказ, возвращает сгенерированный id.
func (r *Database) Save(ctx context.Context, o *model.Product) (int64, error) {
	const query = `INSERT INTO Products (Name) VALUES ($1) RETURNING id`

	var id int64
	err := r.DB.QueryRowContext(ctx, query, o.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert Product: %w", err)
	}

	o.ID = id // заполняем id в самой структуре
	return id, nil
}

// Update — обновить существующий заказ.
// Возвращает ErrNotFound, если записи с таким id нет.
func (r *Database) Update(ctx context.Context, o *model.Product) error {
	const query = `UPDATE Products SET Name = $1 WHERE id = $2`

	res, err := r.DB.ExecContext(ctx, query, o.Name, o.ID)
	if err != nil {
		return fmt.Errorf("update Product: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete — удалить заказ по id.
// Возвращает ErrNotFound, если записи не было.
func (r *Database) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM Products WHERE id = $1`

	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete Product: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// FindAll — вернуть все заказы (с пагинацией).
func (r *Database) FindAll(ctx context.Context, limit, offset int) ([]model.Product, error) {
	const query = `SELECT id, Name FROM Products Product BY id LIMIT $1 OFFSET $2`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list Products: %w", err)
	}
	defer rows.Close()

	Products := make([]model.Product, 0, limit)
	for rows.Next() {
		var o model.Product
		if err := rows.Scan(&o.ID, &o.Name); err != nil {
			return nil, fmt.Errorf("scan Product: %w", err)
		}
		Products = append(Products, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return Products, nil
}
