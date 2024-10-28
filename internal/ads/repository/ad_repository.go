// internal/ads/repository/ad_repository.go
package repository

import (
	"ads/internal/cache"
	"ads/internal/models"
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdRepository struct {
	dbPool  *pgxpool.Pool
	rClient *cache.RedisCache
}

func NewAdRepository(dbPool *pgxpool.Pool, rClient *cache.RedisCache) *AdRepository {
	return &AdRepository{dbPool: dbPool, rClient: rClient}
}

func (r *AdRepository) CreateAd(ctx *context.Context, ad *models.CreateAd) (string, error) {
	var id string
	query := `INSERT INTO ads (name, description, price, created_at, is_active) VALUES ($1, $2, $3, $4, $5) returning id`

	err := r.dbPool.QueryRow(*ctx, query, ad.Name, ad.Description, ad.Price, ad.CreatedAt, ad.IsActive).Scan(&id)
	return id, err
}

func (r *AdRepository) EditAd(ctx *context.Context, ad models.CreateAd) (string, error) {
	var id string
	query := `update ads set name = $1, description = $2, price = $3, created_at = $4, is_active = $5 where id = $6 returning id`

	err := r.dbPool.QueryRow(*ctx, query, ad.Name, ad.Description, ad.Price, ad.CreatedAt, ad.IsActive, ad.ID).Scan(&id)
	return id, err
}

func (r *AdRepository) GetAdByID(ctx *context.Context, id int) (*models.Ad, error) {
	var ad models.Ad
	query := `SELECT id, name, description, price, created_at, is_active FROM ads WHERE id = $1`

	err := r.dbPool.QueryRow(*ctx, query, id).Scan(&ad.ID, &ad.Name, &ad.Description, &ad.Price, &ad.CreatedAt, &ad.IsActive)
	return &ad, err
}

func (r *AdRepository) GetAds(ctx *context.Context, offset, limit int, orderBy string) (*[]models.Ad, error) {
	var ads = make([]models.Ad, 0)
	query := `SELECT id, name, description, price, created_at, is_active FROM ads ` + orderBy + ` offset $1 limit $2`

	err := pgxscan.Select(*ctx, r.dbPool, &ads, query, offset, limit)
	return &ads, err
}

func (r *AdRepository) DeleteAdByID(ctx *context.Context, id int) error {
	query := `DELETE FROM ads WHERE id = $1 returning id`
	err := r.dbPool.QueryRow(*ctx, query, id).Scan(&id)
	return err
}

// Other repository methods as needed
