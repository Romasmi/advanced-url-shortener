package repository

import (
	"context"
	"errors"
	"fmt"
	"shorturl/internal/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlRepository struct {
	db *pgxpool.Pool
}

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

const UrlsTable = "urls"

func NewUrlRepository(db *pgxpool.Pool) *UrlRepository {
	return &UrlRepository{
		db: db,
	}
}

func (r *UrlRepository) Create(ctx context.Context, url *models.Url) (*models.Url, error) {
	query := fmt.Sprintf(`
		INSERT INTO %v (id, original_url, short_url)
		VALUES ($1, $2, $3)
		RETURNING *
	`, UrlsTable)

	var newUrl *models.Url
	err := r.db.QueryRow(ctx, query, url.ID, url.OriginalUrl, url.ShortUrl).Scan(&newUrl)
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
	query := fmt.Sprintf(`
		SELECT * 
		FROM %v
		WHERE original_url = $1
		LIMIT 1
	`, UrlsTable)
	var url *models.Url
	err := r.db.QueryRow(ctx, query, originalUrl).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return url, nil
}
