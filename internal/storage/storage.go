package storage

import (
	"context"

	"github.com/demo-marketplace/internal/model"
	"github.com/google/uuid"
)

type ApplicationStore interface {
	CreateApplication(ctx context.Context, app model.Application) (model.Application, error)
	GetApplication(ctx context.Context, id uuid.UUID) (model.Application, error)
	ListApplications(ctx context.Context, activeOnly bool) ([]model.Application, error)
	UpdateApplication(ctx context.Context, app model.Application) (model.Application, error)
	DeleteApplication(ctx context.Context, id uuid.UUID) error
}

type CategoryStore interface {
	CreateCategory(ctx context.Context, cat model.Category) (model.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (model.Category, error)
	ListCategories(ctx context.Context) ([]model.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type ConnectionStore interface {
	CreateConnection(ctx context.Context, conn model.Connection) (model.Connection, error)
	ListConnectionsByApp(ctx context.Context, appID uuid.UUID) ([]model.Connection, error)
	UpdateConnectionStatus(ctx context.Context, id uuid.UUID, status string) (model.Connection, error)
	DeleteConnection(ctx context.Context, id uuid.UUID) error
}

type Store interface {
	ApplicationStore
	CategoryStore
	ConnectionStore
}
