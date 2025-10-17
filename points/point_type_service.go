package points

import (
	"context"
	"strings"

	d "github.com/usual2970/acto/domain/points"
)

// UpdatePointTypeRequest represents the request for updating a point type
type UpdatePointTypeRequest struct {
	DisplayName *string `json:"displayName,omitempty"`
	Description *string `json:"description,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
}

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

func (s *PointTypeService) Update(ctx context.Context, name string, updates UpdatePointTypeRequest) error {
	// 首先验证积分类型是否存在
	existing, err := s.repo.GetPointTypeByName(ctx, name)
	if err != nil {
		return err
	}
	if existing == nil {
		return d.ErrPointTypeNotFound
	}

	// 构建更新对象，只更新提供的字段
	updated := *existing
	if updates.DisplayName != nil {
		updated.DisplayName = *updates.DisplayName
	}
	if updates.Description != nil {
		updated.Description = *updates.Description
	}
	if updates.Enabled != nil {
		updated.Enabled = *updates.Enabled
	}

	return s.repo.UpdatePointType(ctx, updated)
}

func (s *PointTypeService) Delete(ctx context.Context, name string) error {
	// 首先验证积分类型是否存在
	existing, err := s.repo.GetPointTypeByName(ctx, name)
	if err != nil {
		return err
	}
	if existing == nil {
		return d.ErrPointTypeNotFound
	}

	// 检查是否已经被软删除
	if existing.DeletedAt != nil {
		return d.ErrPointTypeAlreadyDeleted
	}

	// 检查是否有余额
	hasBalances, err := s.repo.HasBalances(ctx, existing.ID)
	if err != nil {
		return err
	}
	if hasBalances {
		return d.ErrPointTypeInUse
	}

	// 执行软删除
	return s.repo.SoftDeletePointType(ctx, name)
}

func (s *PointTypeService) GetByID(ctx context.Context, id string) (*d.PointType, error) {
	return s.repo.GetPointTypeByID(ctx, id)
}

func (s *PointTypeService) List(ctx context.Context, limit, offset int) ([]d.PointType, error) {
	return s.repo.ListPointTypes(ctx, limit, offset)
}
