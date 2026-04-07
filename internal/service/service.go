package service

import (
	"context"
	"fmt"

	"github.com/demo-marketplace/internal/model"
	"github.com/demo-marketplace/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	store storage.Store
}

func New(store storage.Store) *Service {
	return &Service{store: store}
}

// --- Applications ---

func (s *Service) CreateApplication(ctx context.Context, app model.Application) (model.Application, error) {
	if app.Name == "" {
		return model.Application{}, fmt.Errorf("application name is required")
	}
	if app.CategoryID == uuid.Nil {
		return model.Application{}, fmt.Errorf("category ID is required")
	}

	// Verify the category exists.
	if _, err := s.store.GetCategory(ctx, app.CategoryID); err != nil {
		return model.Application{}, fmt.Errorf("invalid category: %w", err)
	}

	return s.store.CreateApplication(ctx, app)
}

func (s *Service) GetApplication(ctx context.Context, id uuid.UUID) (model.Application, error) {
	return s.store.GetApplication(ctx, id)
}

func (s *Service) ListApplications(ctx context.Context, activeOnly bool) ([]model.Application, error) {
	return s.store.ListApplications(ctx, activeOnly)
}

func (s *Service) UpdateApplication(ctx context.Context, app model.Application) (model.Application, error) {
	if app.Name == "" {
		return model.Application{}, fmt.Errorf("application name is required")
	}
	return s.store.UpdateApplication(ctx, app)
}

func (s *Service) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteApplication(ctx, id)
}

// --- Categories ---

func (s *Service) CreateCategory(ctx context.Context, cat model.Category) (model.Category, error) {
	if cat.Name == "" {
		return model.Category{}, fmt.Errorf("category name is required")
	}
	return s.store.CreateCategory(ctx, cat)
}

func (s *Service) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.store.ListCategories(ctx)
}

func (s *Service) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteCategory(ctx, id)
}

// --- Connections ---

func (s *Service) CreateConnection(ctx context.Context, conn model.Connection) (model.Connection, error) {
	if conn.ApplicationID == uuid.Nil {
		return model.Connection{}, fmt.Errorf("application ID is required")
	}
	if conn.MerchantID == "" {
		return model.Connection{}, fmt.Errorf("merchant ID is required")
	}

	// Verify the application exists.
	if _, err := s.store.GetApplication(ctx, conn.ApplicationID); err != nil {
		return model.Connection{}, fmt.Errorf("invalid application: %w", err)
	}

	conn.Status = model.ConnectionStatusActive
	return s.store.CreateConnection(ctx, conn)
}

func (s *Service) ListConnectionsByApp(ctx context.Context, appID uuid.UUID) ([]model.Connection, error) {
	return s.store.ListConnectionsByApp(ctx, appID)
}

func (s *Service) UpdateConnectionStatus(ctx context.Context, id uuid.UUID, status string) (model.Connection, error) {
	if status != model.ConnectionStatusActive && status != model.ConnectionStatusInactive {
		return model.Connection{}, fmt.Errorf("invalid status: %s", status)
	}
	return s.store.UpdateConnectionStatus(ctx, id, status)
}

func (s *Service) DeleteConnection(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteConnection(ctx, id)
}

// --- Reviews ---

func (s *Service) CreateReview(ctx context.Context, review model.Review) (model.Review, error) {
	if review.ApplicationID == uuid.Nil {
		return model.Review{}, fmt.Errorf("application ID is required")
	}
	if review.MerchantID == "" {
		return model.Review{}, fmt.Errorf("merchant ID is required")
	}

	// Verify the application exists.
	if _, err := s.store.GetApplication(ctx, review.ApplicationID); err != nil {
		return model.Review{}, fmt.Errorf("invalid application: %w", err)
	}

	return s.store.CreateReview(ctx, review)
}

func (s *Service) ListReviewsByApp(ctx context.Context, appID uuid.UUID) ([]model.Review, error) {
	return s.store.ListReviewsByApp(ctx, appID)
}

func (s *Service) GetAverageRating(ctx context.Context, appID uuid.UUID) (float64, error) {
	return s.store.GetAverageRating(ctx, appID)
}

func (s *Service) DeleteReview(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteReview(ctx, id)
}
