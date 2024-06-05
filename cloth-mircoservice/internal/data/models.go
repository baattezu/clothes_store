package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}
type Models struct {
	ClotheInfo DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		ClotheInfo: DBModel{DB: db},
	}
}

func (m DBModel) Insert(clothe *ClotheInfo) error {
	query := `INSERT INTO clothes_info (cloth_name, cloth_cost, cloth_size)VALUES ($1, $2, $3 )RETURNING id, created_at, updated_at, version`
	args := []any{clothe.ClothName, clothe.ClothCost, clothe.ClothSize}
	return m.DB.QueryRow(query, args...).Scan(&clothe.ID, &clothe.CreatedAt, &clothe.UpdatedAt, &clothe.Version)
}

func (m DBModel) Get(id int64) (*ClotheInfo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, updated_at, cloth_name, cloth_cost, cloth_size, version FROM clothes_info WHERE id = $1`
	var clothe ClotheInfo

	err := m.DB.QueryRow(query, id).Scan(
		&clothe.ID,
		&clothe.CreatedAt,
		&clothe.UpdatedAt,
		&clothe.ClothName,
		&clothe.ClothCost,
		&clothe.ClothSize,
		&clothe.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &clothe, nil
}

func (m DBModel) Update(clothe *ClotheInfo) error {
	query := `UPDATE clothes_info SET cloth_name= $1, cloth_cost = $2, cloth_size = $3, version = version + 1 WHERE id = $4 RETURNING version`
	args := []any{
		clothe.ClothName,
		clothe.ClothCost,
		clothe.ClothSize,
		clothe.ID,
	}
	return m.DB.QueryRow(query, args...).Scan(&clothe.Version)
}

func (m DBModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM clothes_info WHERE id = $1`
	result, err := m.DB.Exec(query, id)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m DBModel) GetAll(cloth_name string, cloth_size string, filters Filters) ([]*ClotheInfo, error) {
	query := fmt.Sprintf(`SELECT id, created_at, updated_at, cloth_name, cloth_cost, cloth_size, version
	FROM clothes_info
	WHERE (to_tsvector('simple', cloth_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND  (LOWER(cloth_size) = LOWER($2) OR $2 = '')
	ORDER BY  %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{cloth_name, cloth_size, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	clothes_info := []*ClotheInfo{}

	for rows.Next() {
		var clothe ClotheInfo

		err := rows.Scan(
			&clothe.ID,
			&clothe.CreatedAt,
			&clothe.UpdatedAt,
			&clothe.ClothName,
			&clothe.ClothCost,
			&clothe.ClothSize,
			&clothe.Version,
		)
		if err != nil {
			return nil, err
		}
		clothes_info = append(clothes_info, &clothe)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return clothes_info, nil
}
