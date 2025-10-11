package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/advanced-url-shortener/internal/models"
	"github.com/redis/go-redis/v9"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

const UrlsTable = "urls"

func NewUrlRepository(db *pgxpool.Pool, redis *redis.Client) *UrlRepository {
	return &UrlRepository{
		db:    db,
		redis: redis,
	}
}

func (r *UrlRepository) Create(ctx context.Context, url *models.Url) (*models.Url, error) {
	query := fmt.Sprintf(`
		INSERT INTO %v (id, original_url, code)
		VALUES ($1, $2, $3)
		RETURNING *
	`, UrlsTable)

	var newUrl *models.Url
	err := r.db.QueryRow(ctx, query, url.ID, url.OriginalUrl, url.Code).Scan(&newUrl)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrDuplicate
			}
			return nil, err
		}
		return nil, err
	}
	return newUrl, nil
}

func (r *UrlRepository) GetByOriginalUrl(ctx context.Context, originalUrl string) (*models.Url, error) {
	return r.getByColumn(ctx, "original_url", originalUrl)
}

func (r *UrlRepository) GetByCode(ctx context.Context, code string) (*models.Url, error) {
	return r.getByColumn(ctx, "code", code)
}

func (r *UrlRepository) getByColumn(ctx context.Context, columnName, value string) (*models.Url, error) {
	query := fmt.Sprintf(`
        SELECT * 
        FROM %v
        WHERE %s = $1
        LIMIT 1
    `, UrlsTable, columnName)

	var url models.Url
	err := r.db.QueryRow(ctx, query, value).Scan(
		&url.ID,
		&url.OriginalUrl,
		&url.Code,
		&url.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get url by %s: %w", columnName, err)
	}

	return &url, nil
}
