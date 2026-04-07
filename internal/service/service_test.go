package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/demo-marketplace/internal/model"
	"github.com/demo-marketplace/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockStore implements storage.Store for testing.
type mockStore struct {
	mock.Mock
}

func (m *mockStore) CreateApplication(ctx context.Context, app model.Application) (model.Application, error) {
	args := m.Called(ctx, app)
	return args.Get(0).(model.Application), args.Error(1)
}

func (m *mockStore) GetApplication(ctx context.Context, id uuid.UUID) (model.Application, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Application), args.Error(1)
}

func (m *mockStore) ListApplications(ctx context.Context, activeOnly bool) ([]model.Application, error) {
	args := m.Called(ctx, activeOnly)
	return args.Get(0).([]model.Application), args.Error(1)
}

func (m *mockStore) UpdateApplication(ctx context.Context, app model.Application) (model.Application, error) {
	args := m.Called(ctx, app)
	return args.Get(0).(model.Application), args.Error(1)
}

func (m *mockStore) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockStore) CreateCategory(ctx context.Context, cat model.Category) (model.Category, error) {
	args := m.Called(ctx, cat)
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *mockStore) GetCategory(ctx context.Context, id uuid.UUID) (model.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *mockStore) ListCategories(ctx context.Context) ([]model.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *mockStore) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockStore) CreateConnection(ctx context.Context, conn model.Connection) (model.Connection, error) {
	args := m.Called(ctx, conn)
	return args.Get(0).(model.Connection), args.Error(1)
}

func (m *mockStore) ListConnectionsByApp(ctx context.Context, appID uuid.UUID) ([]model.Connection, error) {
	args := m.Called(ctx, appID)
	return args.Get(0).([]model.Connection), args.Error(1)
}

func (m *mockStore) UpdateConnectionStatus(ctx context.Context, id uuid.UUID, status string) (model.Connection, error) {
	args := m.Called(ctx, id, status)
	return args.Get(0).(model.Connection), args.Error(1)
}

func (m *mockStore) DeleteConnection(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockStore) CreateReview(ctx context.Context, review model.Review) (model.Review, error) {
	args := m.Called(ctx, review)
	return args.Get(0).(model.Review), args.Error(1)
}

func (m *mockStore) ListReviewsByApp(ctx context.Context, appID uuid.UUID) ([]model.Review, error) {
	args := m.Called(ctx, appID)
	return args.Get(0).([]model.Review), args.Error(1)
}

func (m *mockStore) GetAverageRating(ctx context.Context, appID uuid.UUID) (float64, error) {
	args := m.Called(ctx, appID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockStore) DeleteReview(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func TestCreateApplication(t *testing.T) {
	ctx := context.Background()
	catID := uuid.New()

	t.Run("success", func(t *testing.T) {
		store := new(mockStore)
		svc := service.New(store)

		app := model.Application{
			Name:       "Test App",
			CategoryID: catID,
		}

		store.On("GetCategory", ctx, catID).Return(model.Category{ID: catID}, nil)
		store.On("CreateApplication", ctx, app).Return(model.Application{
			ID:         uuid.New(),
			Name:       "Test App",
			CategoryID: catID,
		}, nil)

		result, err := svc.CreateApplication(ctx, app)
		assert.NoError(t, err)
		assert.Equal(t, "Test App", result.Name)
		store.AssertExpectations(t)
	})

	t.Run("empty name", func(t *testing.T) {
		store := new(mockStore)
		svc := service.New(store)

		_, err := svc.CreateApplication(ctx, model.Application{CategoryID: catID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("missing category", func(t *testing.T) {
		store := new(mockStore)
		svc := service.New(store)

		_, err := svc.CreateApplication(ctx, model.Application{Name: "Test"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category ID is required")
	})

	t.Run("invalid category", func(t *testing.T) {
		store := new(mockStore)
		svc := service.New(store)

		app := model.Application{Name: "Test", CategoryID: catID}
		store.On("GetCategory", ctx, catID).Return(model.Category{}, fmt.Errorf("not found"))

		_, err := svc.CreateApplication(ctx, app)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid category")
	})
}

func TestCreateConnection(t *testing.T) {
	ctx := context.Background()
	appID := uuid.New()

	t.Run("success", func(t *testing.T) {
		store := new(mockStore)
		svc := service.New(store)

		conn := model.Connection{
			ApplicationID: appID,
			MerchantID:    "merchant-123",
		}

		store.On("GetApplication", ctx, appID).Return(model.Application{ID: appID}, nil)
		store.On("CreateConnection", ctx, mock.MatchedBy(func(c model.Connection) bool {
			return c.Status == model.ConnectionStatusActive
		})).Return(model.Connection{
			ID:            uuid.New(),
			ApplicationID: appID,
			MerchantID:    "merchant-123",
			Status:        model.ConnectionStatusActive,
		}, nil)

		result, err := svc.CreateConnection(ctx, conn)
		assert.NoError(t, err)
		assert.Equal(t, model.ConnectionStatusActive, result.Status)
		store.AssertExpectations(t)
	})

	t.Run("invalid status update", func(t *testing.T) {
		store := new(mockStore)
		svc := service.New(store)

		_, err := svc.UpdateConnectionStatus(ctx, uuid.New(), "invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})
}

func TestCreateReview(t *testing.T) {
	ctx := context.Background()
	appID := uuid.New()

	t.Run("success", func(t *testing.T) {
		store := new(mockStore)
		svc := service.New(store)

		review := model.Review{
			ApplicationID: appID,
			MerchantID:    "merchant-123",
			Rating:        5,
			Body:          "Great app!",
		}

		store.On("GetApplication", ctx, appID).Return(model.Application{ID: appID}, nil)
		store.On("CreateReview", ctx, review).Return(model.Review{
			ID:            uuid.New(),
			ApplicationID: appID,
			MerchantID:    "merchant-123",
			Rating:        5,
			Body:          "Great app!",
		}, nil)

		result, err := svc.CreateReview(ctx, review)
		assert.NoError(t, err)
		assert.Equal(t, 5, result.Rating)
		assert.Equal(t, "Great app!", result.Body)
		store.AssertExpectations(t)
	})
}
