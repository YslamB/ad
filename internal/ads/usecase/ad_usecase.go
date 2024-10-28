package usecase

import (
	"ads/internal/ads/repository"
	"ads/internal/cache"
	"ads/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdUsecase struct {
	repo    repository.AdRepository
	rClient *cache.RedisCache
}

func NewAdUsecase(dbPool *pgxpool.Pool, rClient *cache.RedisCache) *AdUsecase {
	return &AdUsecase{repo: *repository.NewAdRepository(dbPool, rClient), rClient: rClient}
}

func (u *AdUsecase) CreateAd(ctx *context.Context, ad models.CreateAd) models.Response {

	if ad.CreatedAt.IsZero() {
		ad.CreatedAt = time.Now()
	}
	id, err := u.repo.CreateAd(ctx, &ad)

	if err != nil {
		// todo: handle conflict error
		return models.Response{Error: err, Status: 500}
	}
	ad.ID, _ = strconv.Atoi(id)
	u.rClient.CreateAd(ctx, &ad)
	return models.Response{Status: 200, Data: models.DataMessage{Message: "ad created ", ID: id}}
}

func (u *AdUsecase) EditAd(ctx *context.Context, ad models.CreateAd) models.Response {
	id, err := u.repo.EditAd(ctx, ad)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return models.Response{Status: 404, Error: err}
		}
		return models.Response{Error: err, Status: 500}
	}
	u.rClient.CreateAd(ctx, &ad)
	return models.Response{Status: 200, Data: models.DataMessage{Message: "ad Updated", ID: id}}
}

func (u *AdUsecase) GetAdByID(ctx *context.Context, id int) models.Response {
	data, err := u.rClient.GetAd(*ctx, fmt.Sprintf("ad:%d", id))

	if err != nil {
		data, err = u.repo.GetAdByID(ctx, id)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				return models.Response{Status: 404, Error: err}
			}
			return models.Response{Error: err, Status: 500}
		} else {
			u.rClient.SetAd(*ctx, data)
		}
	}
	return models.Response{Status: 200, Data: data}
}

func (u *AdUsecase) DeleteAdByID(ctx *context.Context, id int) models.Response {

	u.rClient.Delete(*ctx, fmt.Sprintf("ad:%d", id))
	err := u.repo.DeleteAdByID(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Response{Status: 404, Error: err}
		}
		return models.Response{Error: err, Status: 400}
	}
	return models.Response{Status: 200, Data: models.Success}
}

func (u *AdUsecase) GetAds(ctx *context.Context, page, limit int, orderBy string) models.Response {
	offset := (page - 1) * limit
	data, err := u.rClient.GetAds(ctx, offset, limit, orderBy)

	if err != nil {
		data, err = u.repo.GetAds(ctx, offset, limit, orderBy)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				return models.Response{Status: 404, Error: err}
			}
			return models.Response{Error: err, Status: 500}
		}
	}

	return models.Response{Status: 200, Data: data}

}
