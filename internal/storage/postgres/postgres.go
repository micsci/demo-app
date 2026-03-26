package postgres

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/demo-marketplace/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Store struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// --- Applications ---

func (s *Store) CreateApplication(ctx context.Context, app model.Application) (model.Application, error) {
	app.ID = uuid.New()

	query, args, err := psql.Insert("application").
		Columns("id", "name", "display_name", "description", "active", "visible", "category_id").
		Values(app.ID, app.Name, app.DisplayName, app.Description, app.Active, app.Visible, app.CategoryID).
		Suffix("RETURNING id, name, display_name, description, active, visible, category_id, created_at, updated_at").
		ToSql()
	if err != nil {
		return model.Application{}, fmt.Errorf("building query: %w", err)
	}

	var result model.Application
	if err := s.db.QueryRowxContext(ctx, query, args...).StructScan(&result); err != nil {
		return model.Application{}, fmt.Errorf("creating application: %w", err)
	}
	return result, nil
}

func (s *Store) GetApplication(ctx context.Context, id uuid.UUID) (model.Application, error) {
	query, args, err := psql.Select("id", "name", "display_name", "description", "active", "visible", "category_id", "created_at", "updated_at").
		From("application").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Application{}, fmt.Errorf("building query: %w", err)
	}

	var app model.Application
	if err := s.db.QueryRowxContext(ctx, query, args...).StructScan(&app); err != nil {
		if err == sql.ErrNoRows {
			return model.Application{}, fmt.Errorf("application not found: %s", id)
		}
		return model.Application{}, fmt.Errorf("getting application: %w", err)
	}
	return app, nil
}

func (s *Store) ListApplications(ctx context.Context, activeOnly bool) ([]model.Application, error) {
	qb := psql.Select("id", "name", "display_name", "description", "active", "visible", "category_id", "created_at", "updated_at").
		From("application").
		OrderBy("created_at DESC")

	if activeOnly {
		qb = qb.Where(sq.Eq{"active": true})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	var apps []model.Application
	if err := s.db.SelectContext(ctx, &apps, query, args...); err != nil {
		return nil, fmt.Errorf("listing applications: %w", err)
	}
	return apps, nil
}

func (s *Store) UpdateApplication(ctx context.Context, app model.Application) (model.Application, error) {
	query, args, err := psql.Update("application").
		Set("name", app.Name).
		Set("display_name", app.DisplayName).
		Set("description", app.Description).
		Set("active", app.Active).
		Set("visible", app.Visible).
		Set("category_id", app.CategoryID).
		Where(sq.Eq{"id": app.ID}).
		Suffix("RETURNING id, name, display_name, description, active, visible, category_id, created_at, updated_at").
		ToSql()
	if err != nil {
		return model.Application{}, fmt.Errorf("building query: %w", err)
	}

	var result model.Application
	if err := s.db.QueryRowxContext(ctx, query, args...).StructScan(&result); err != nil {
		return model.Application{}, fmt.Errorf("updating application: %w", err)
	}
	return result, nil
}

func (s *Store) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Delete("application").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("building query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("deleting application: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("application not found: %s", id)
	}
	return nil
}

// --- Categories ---

func (s *Store) CreateCategory(ctx context.Context, cat model.Category) (model.Category, error) {
	cat.ID = uuid.New()

	query, args, err := psql.Insert("category").
		Columns("id", "name").
		Values(cat.ID, cat.Name).
		Suffix("RETURNING id, name, created_at, updated_at").
		ToSql()
	if err != nil {
		return model.Category{}, fmt.Errorf("building query: %w", err)
	}

	var result model.Category
	if err := s.db.QueryRowxContext(ctx, query, args...).StructScan(&result); err != nil {
		return model.Category{}, fmt.Errorf("creating category: %w", err)
	}
	return result, nil
}

func (s *Store) GetCategory(ctx context.Context, id uuid.UUID) (model.Category, error) {
	query, args, err := psql.Select("id", "name", "created_at", "updated_at").
		From("category").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Category{}, fmt.Errorf("building query: %w", err)
	}

	var cat model.Category
	if err := s.db.QueryRowxContext(ctx, query, args...).StructScan(&cat); err != nil {
		if err == sql.ErrNoRows {
			return model.Category{}, fmt.Errorf("category not found: %s", id)
		}
		return model.Category{}, fmt.Errorf("getting category: %w", err)
	}
	return cat, nil
}

func (s *Store) ListCategories(ctx context.Context) ([]model.Category, error) {
	query, args, err := psql.Select("id", "name", "created_at", "updated_at").
		From("category").
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	var cats []model.Category
	if err := s.db.SelectContext(ctx, &cats, query, args...); err != nil {
		return nil, fmt.Errorf("listing categories: %w", err)
	}
	return cats, nil
}

func (s *Store) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Delete("category").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("building query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("deleting category: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("category not found: %s", id)
	}
	return nil
}

// --- Connections ---

func (s *Store) CreateConnection(ctx context.Context, conn model.Connection) (model.Connection, error) {
	conn.ID = uuid.New()

	query, args, err := psql.Insert("connection").
		Columns("id", "application_id", "merchant_id", "status").
		Values(conn.ID, conn.ApplicationID, conn.MerchantID, conn.Status).
		Suffix("RETURNING id, application_id, merchant_id, status, created_at, updated_at").
		ToSql()
	if err != nil {
		return model.Connection{}, fmt.Errorf("building query: %w", err)
	}

	var result model.Connection
	if err := s.db.QueryRowxContext(ctx, query, args...).StructScan(&result); err != nil {
		return model.Connection{}, fmt.Errorf("creating connection: %w", err)
	}
	return result, nil
}

func (s *Store) ListConnectionsByApp(ctx context.Context, appID uuid.UUID) ([]model.Connection, error) {
	query, args, err := psql.Select("id", "application_id", "merchant_id", "status", "created_at", "updated_at").
		From("connection").
		Where(sq.Eq{"application_id": appID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	var conns []model.Connection
	if err := s.db.SelectContext(ctx, &conns, query, args...); err != nil {
		return nil, fmt.Errorf("listing connections: %w", err)
	}
	return conns, nil
}

func (s *Store) UpdateConnectionStatus(ctx context.Context, id uuid.UUID, status string) (model.Connection, error) {
	query, args, err := psql.Update("connection").
		Set("status", status).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, application_id, merchant_id, status, created_at, updated_at").
		ToSql()
	if err != nil {
		return model.Connection{}, fmt.Errorf("building query: %w", err)
	}

	var result model.Connection
	if err := s.db.QueryRowxContext(ctx, query, args...).StructScan(&result); err != nil {
		return model.Connection{}, fmt.Errorf("updating connection: %w", err)
	}
	return result, nil
}

func (s *Store) DeleteConnection(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Delete("connection").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("building query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("deleting connection: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("connection not found: %s", id)
	}
	return nil
}
