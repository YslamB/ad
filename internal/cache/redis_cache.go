// internal/cache/redis_cache.go
package cache

import (
	"ads/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) SetAd(ctx context.Context, ad *models.Ad) error {
	data, err := json.Marshal(ad)

	if err != nil {
		return err
	}
	key := "ad:" + strconv.Itoa(ad.ID) // Use a unique key, e.g., "ad:1"

	if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
		return err
	}

	// Add to sorted sets for ordering by created_at and price
	if err := r.client.ZAdd(ctx, "ads_by_created_at", &redis.Z{Score: float64(ad.CreatedAt.Unix()), Member: key}).Err(); err != nil {
		return err
	}

	if err := r.client.ZAdd(ctx, "ads_by_price", &redis.Z{Score: ad.Price, Member: key}).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) CreateAd(ctx *context.Context, ad *models.CreateAd) error {
	data, err := json.Marshal(ad)

	if err != nil {
		return err
	}
	key := "ad:" + strconv.Itoa(ad.ID)

	if err := r.client.Set(*ctx, key, data, 0).Err(); err != nil {
		return err
	}

	// Add to sorted sets for ordering by created_at and price
	if err := r.client.ZAdd(*ctx, "ads_by_created_at", &redis.Z{Score: float64(ad.CreatedAt.Unix()), Member: key}).Err(); err != nil {
		return err
	}

	if err := r.client.ZAdd(*ctx, "ads_by_price", &redis.Z{Score: ad.Price, Member: key}).Err(); err != nil {
		return err
	}
	return nil
}

// Retrieve ads with pagination and orderBy options
func (r *RedisCache) GetAds(ctx *context.Context, offset, limit int, orderBy string) (*[]models.Ad, error) {
	var sortSet string
	var orderFunc func(start, stop int64) *redis.StringSliceCmd

	// Determine which sorted set to use and order direction
	switch orderBy {
	case "order by created_at":
		sortSet = "ads_by_created_at"
		orderFunc = func(start, stop int64) *redis.StringSliceCmd {
			return r.client.ZRange(*ctx, sortSet, start, stop)
		}
	case "order by created_at desc":
		sortSet = "ads_by_created_at"
		orderFunc = func(start, stop int64) *redis.StringSliceCmd {
			return r.client.ZRevRange(*ctx, sortSet, start, stop)
		}
	case "order by price":
		sortSet = "ads_by_price"
		orderFunc = func(start, stop int64) *redis.StringSliceCmd {
			return r.client.ZRange(*ctx, sortSet, start, stop)
		}
	case "order by price desc":
		sortSet = "ads_by_price"
		orderFunc = func(start, stop int64) *redis.StringSliceCmd {
			return r.client.ZRevRange(*ctx, sortSet, start, stop)
		}
	default:
		return nil, errors.New("invalid orderBy parameter")
	}

	// Fetch a range of keys based on offset and limit
	keys, err := orderFunc(int64(offset), int64(offset+limit-1)).Result()

	if err != nil {
		return nil, err
	}

	// Retrieve each ad by its key
	ads := make([]models.Ad, 0, len(keys))

	for _, key := range keys {
		result, err := r.client.Get(*ctx, key).Result()

		if err == redis.Nil {
			continue // Skip missing ads
		} else if err != nil {
			return nil, err
		}

		var ad models.Ad

		if err := json.Unmarshal([]byte(result), &ad); err != nil {
			return nil, err
		}

		ads = append(ads, ad)
	}

	return &ads, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) SetAllADs(ctx context.Context, dbPool *pgxpool.Pool) error {
	var ads []models.Ad
	rows, _ := dbPool.Query(ctx, `SELECT id, name, description, price, created_at, is_active FROM ads`)
	for rows.Next() {
		var ad models.Ad
		rows.Scan(&ad.ID, &ad.Name, &ad.Description, &ad.Price, &ad.CreatedAt, &ad.IsActive)
		ads = append(ads, ad)
	}
	fmt.Println("ads")
	fmt.Println(len(ads))
	for i := range ads {
		r.SetAd(ctx, &ads[i])
	}
	return nil
}

func (r *RedisCache) GetAd(ctx context.Context, key string) (*models.Ad, error) {
	// Get the JSON data from Redis
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("ad with ID %s not found", key)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get ad: %w", err)
	}

	// Unmarshal the JSON data into an Ad struct
	var ad models.Ad
	if err := json.Unmarshal([]byte(result), &ad); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ad data: %w", err)
	}

	return &ad, nil
}
