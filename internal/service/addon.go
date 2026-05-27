package service

import (
	"context"
	"pizza-backend/internal/model"
)

type AddonRepository interface {
	GetAddonListByCategoryID(ctx context.Context, categoryID int64) ([]model.AddonRow, error)
}

type AddonService struct {
	repo AddonRepository
}

func NewAddonService(r AddonRepository) *AddonService {
	return &AddonService{repo: r}
}

func (s *AddonService) GetAddonListByCategoryID(ctx context.Context, categoryID int64) ([]model.Addon, error) {
	rows, err := s.repo.GetAddonListByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	result := make([]model.Addon, 0, len(rows))

	for _, row := range rows {
		result = append(result, model.Addon{
			ID:    row.ID,
			Name:  row.Name,
			Price: float64(row.Price) / 100,
		})
	}

	return result, nil
}
