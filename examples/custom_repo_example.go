package examples

import (
	d "acto/domain/points"
	uc "acto/points"
	"context"
)

// CustomPointTypeRepo is an example override implementing PointTypeRepository.
type CustomPointTypeRepo struct{}

var _ uc.PointTypeRepository = (*CustomPointTypeRepo)(nil)

func (r *CustomPointTypeRepo) CreatePointType(ctx context.Context, pt d.PointType) (string, error) {
	return "custom-id", nil
}
func (r *CustomPointTypeRepo) UpdatePointType(ctx context.Context, pt d.PointType) error { return nil }
func (r *CustomPointTypeRepo) DeletePointType(ctx context.Context, pointTypeID string) error {
	return nil
}
func (r *CustomPointTypeRepo) SoftDeletePointType(ctx context.Context, name string) error { return nil }
func (r *CustomPointTypeRepo) GetPointTypeByID(ctx context.Context, pointTypeID string) (*d.PointType, error) {
	return &d.PointType{ID: pointTypeID, Name: "gold-points"}, nil
}
func (r *CustomPointTypeRepo) GetPointTypeByName(ctx context.Context, name string) (*d.PointType, error) {
	return &d.PointType{ID: "custom-id", Name: name}, nil
}
func (r *CustomPointTypeRepo) ListPointTypes(ctx context.Context, limit, offset int) ([]d.PointType, error) {
	return []d.PointType{}, nil
}
func (r *CustomPointTypeRepo) HasBalances(ctx context.Context, pointTypeID string) (bool, error) {
	return false, nil
}
