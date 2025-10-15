package points

import (
	d "acto/domain/points"
	"context"
	"strings"
)

type PointTypeService struct {
	repo PointTypeRepository
}

func NewPointTypeService(repo PointTypeRepository) *PointTypeService {
	return &PointTypeService{repo: repo}
}

func (s *PointTypeService) Create(ctx context.Context, name, displayName, description string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", d.ErrDuplicatePointTypeName
	}
	return s.repo.CreatePointType(ctx, d.PointType{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Enabled:     true,
	})
}

func (s *PointTypeService) Update(ctx context.Context, pt d.PointType) error {
	return s.repo.UpdatePointType(ctx, pt)
}

func (s *PointTypeService) Delete(ctx context.Context, id string, hasBalances bool) error {
	if hasBalances {
		return d.ErrPointTypeInUse
	}
	return s.repo.DeletePointType(ctx, id)
}

func (s *PointTypeService) GetByID(ctx context.Context, id string) (*d.PointType, error) {
	return s.repo.GetPointTypeByID(ctx, id)
}

func (s *PointTypeService) List(ctx context.Context, limit, offset int) ([]d.PointType, error) {
	return s.repo.ListPointTypes(ctx, limit, offset)
}
