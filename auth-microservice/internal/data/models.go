package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrEditConflict = errors.New("edit conflict")
)

type DBModel struct {
	DB *sql.DB
}
type Models struct {
	ClotheInfo  DBModel
	Tokens      TokenModel
	Users       UserModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		ClotheInfo:  DBModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}

func (m DBModel) Insert(clothe *ClotheInfo) error {

	query := `INSERT INTO module_info (module_name, module_duration, exam_type)VALUES ($1, $2, $3 )RETURNING id, created_at, updated_at, version`
	args := []any{clothe.ModuleName, clothe.ModuleDuration, clothe.ExamType}
	return m.DB.QueryRow(query, args...).Scan(&clothe.ID, &clothe.CreatedAt, &clothe.UpdatedAt, &clothe.Version)
}

func (m DBModel) Get(id int64) (*ClotheInfo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, updated_at, module_name, module_duration, exam_type, version FROM module_info WHERE id = $1`
	var clothe ClotheInfo

	err := m.DB.QueryRow(query, id).Scan(
		&clothe.ID,
		&clothe.CreatedAt,
		&clothe.UpdatedAt,
		&clothe.ModuleName,
		&clothe.ModuleDuration,
		&clothe.ExamType,
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

	query := `UPDATE module_info SET  module_name= $1, module_duration = $2, exam_type = $3, version = version + 1 WHERE id = $4 RETURNING version`
	args := []any{
		clothe.ModuleName,
		clothe.ModuleDuration,
		clothe.ExamType,
		clothe.ID,
	}

	return m.DB.QueryRow(query, args...).Scan(&clothe.Version)
}

func (m DBModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM module_info WHERE id = $1`
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

func (m DBModel) GetAll(module_name string, exam_type string, filters Filters) ([]*ClotheInfo, error) {

	query := fmt.Sprintf(`SELECT id, created_at, updated_at, module_name, module_duration, exam_type, version
	FROM module_info
	WHERE (to_tsvector('simple', module_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND  (LOWER(exam_type) = LOWER($2) OR $2 = '')
	ORDER BY  %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{module_name, exam_type, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	module_info := []*ClotheInfo{}

	for rows.Next() {
		var module_infos ClotheInfo

		err := rows.Scan(
			&module_infos.ID,
			&module_infos.CreatedAt,
			&module_infos.UpdatedAt,
			&module_infos.ModuleName,
			&module_infos.ModuleDuration,
			&module_infos.ExamType,
			&module_infos.Version,
		)
		if err != nil {
			return nil, err
		}
		module_info = append(module_info, &module_infos)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return module_info, nil
}
