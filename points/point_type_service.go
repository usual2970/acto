package points

import (
	"context"
	"errors"
	"strings"

	d "github.com/usual2970/acto/domain/points"
)

type PointTypeService struct {
	repo PointTypeRepository
}

func NewPointTypeService(repo PointTypeRepository) *PointTypeService {
	return &PointTypeService{repo: repo}
}

func (s *PointTypeService) Create(ctx context.Context, req PointTypeCreateRequest) (string, error) {
	uri := strings.TrimSpace(req.URI)
	if uri == "" {
		return "", errors.New("uri cannot be empty")
	}
	rs, err := s.repo.CreatePointType(ctx, d.PointType{
		URI:         uri,
		DisplayName: strings.TrimSpace(req.DisplayName),
		Description: strings.TrimSpace(req.Description),
		Enabled:     true,
	})
	if err != nil {
		return "", errors.New("create point type failed")
	}
	return rs, nil
}

func (s *PointTypeService) Update(ctx context.Context, uri string, updates PointTypeUpdateRequest) error {
	// 首先验证积分类型是否存在
	existing, err := s.repo.GetPointTypeByURI(ctx, uri)
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

func (s *PointTypeService) Delete(ctx context.Context, uri string) error {
	// 首先验证积分类型是否存在
	existing, err := s.repo.GetPointTypeByURI(ctx, uri)
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
	return s.repo.SoftDeletePointType(ctx, uri)
}

func (s *PointTypeService) GetByID(ctx context.Context, id int64) (*d.PointType, error) {
	return s.repo.GetPointTypeByID(ctx, id)
}

func (s *PointTypeService) List(ctx context.Context, limit, offset int) ([]d.PointType, error) {
	return s.repo.ListPointTypes(ctx, limit, offset)
}
