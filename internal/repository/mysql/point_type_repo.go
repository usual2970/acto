package mysql

import (
	"context"
	"database/sql"
	"time"

	d "github.com/usual2970/acto/domain/points"
	uc "github.com/usual2970/acto/points"
)

type PointTypeRepository struct {
	db *sql.DB
}

func NewPointTypeRepository(db *sql.DB) *PointTypeRepository { return &PointTypeRepository{db: db} }

var _ uc.PointTypeRepository = (*PointTypeRepository)(nil)

func (r *PointTypeRepository) CreatePointType(ctx context.Context, pt d.PointType) (string, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO point_types (id,name,display_name,description,enabled,created_at) VALUES (UUID(),?,?,?,?,?)`, pt.Name, pt.DisplayName, pt.Description, pt.Enabled, time.Now().Unix())
	if err != nil {
		return "", err
	}

	return pt.Name, nil
}

func (r *PointTypeRepository) UpdatePointType(ctx context.Context, pt d.PointType) error {
	_, err := r.db.ExecContext(ctx, `UPDATE point_types SET display_name=?, description=?, enabled=? WHERE id=?`, pt.DisplayName, pt.Description, pt.Enabled, pt.ID)
	return err
}

func (r *PointTypeRepository) DeletePointType(ctx context.Context, pointTypeID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM point_types WHERE id=?`, pointTypeID)
	return err
}

func (r *PointTypeRepository) SoftDeletePointType(ctx context.Context, name string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE point_types SET deleted_at=? WHERE name=?`, time.Now().Unix(), name)
	return err
}

func (r *PointTypeRepository) GetPointTypeByID(ctx context.Context, pointTypeID string) (*d.PointType, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id,name,display_name,description,enabled,deleted_at,created_at FROM point_types WHERE id=? AND deleted_at IS NULL`, pointTypeID)
	var pt d.PointType
	var deletedAt *int64
	if err := row.Scan(&pt.ID, &pt.Name, &pt.DisplayName, &pt.Description, &pt.Enabled, &deletedAt, &pt.CreatedAt); err != nil {
		return nil, err
	}
	pt.DeletedAt = deletedAt
	return &pt, nil
}

func (r *PointTypeRepository) GetPointTypeByName(ctx context.Context, name string) (*d.PointType, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id,name,display_name,description,enabled,deleted_at,created_at FROM point_types WHERE name=? AND deleted_at IS NULL`, name)
	var pt d.PointType
	var deletedAt *int64
	if err := row.Scan(&pt.ID, &pt.Name, &pt.DisplayName, &pt.Description, &pt.Enabled, &deletedAt, &pt.CreatedAt); err != nil {
		return nil, err
	}
	pt.DeletedAt = deletedAt
	return &pt, nil
}

func (r *PointTypeRepository) ListPointTypes(ctx context.Context, limit, offset int) ([]d.PointType, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,name,display_name,description,enabled,deleted_at,created_at FROM point_types WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []d.PointType
	for rows.Next() {
		var pt d.PointType
		var deletedAt *int64
		if err := rows.Scan(&pt.ID, &pt.Name, &pt.DisplayName, &pt.Description, &pt.Enabled, &deletedAt, &pt.CreatedAt); err != nil {
			return nil, err
		}
		pt.DeletedAt = deletedAt
		res = append(res, pt)
	}
	return res, rows.Err()
}

func (r *PointTypeRepository) HasBalances(ctx context.Context, pointTypeID string) (bool, error) {
	// Placeholder: real check when balances table exists
	return false, nil
}
