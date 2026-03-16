package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"challenge2/models"
)

type ItemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetAll() ([]models.Item, error) {

	var items []models.Item

	query := `SELECT id, name, description FROM items`

	err := r.db.Select(&items, query)
	if err != nil {
		return nil, fmt.Errorf("Gagal dapat items: %w", err)
	}

	return items, nil
}

func (r *ItemRepository) Create(task *models.Item) error {

	query := `
	INSERT INTO items (name, description)
	VALUES ($1,$2)
	`

	_, err := r.db.Exec(query,
		task.Name,
		task.Description,
	)

	if err != nil {
		return fmt.Errorf("Gagal insert item: %w", err)
	}

	return nil
}

func (r *ItemRepository) DeleteByID(id int) error {

	query := `DELETE FROM items WHERE id=$1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("gagal delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal cek rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item dengan id %d tidak ditemukan", id)
	}

	return nil
}