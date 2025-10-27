package redis

import (
	"context"
	"fmt"

	uc "github.com/usual2970/acto/points"

	goRedis "github.com/redis/go-redis/v9"
)

type RankingRepository struct {
	client *goRedis.Client
}

func NewRankingRepository(client *goRedis.Client) *RankingRepository {
	return &RankingRepository{client: client}
}

var _ uc.RankingRepository = (*RankingRepository)(nil)

func key(pointTypeID int64) string { return "ranking:" + fmt.Sprint(pointTypeID) }

func (r *RankingRepository) UpdateUserScore(ctx context.Context, pointTypeID int64, userID string, score int64) error {
	return r.client.ZAdd(ctx, key(pointTypeID), goRedis.Z{Member: userID, Score: float64(score)}).Err()
}

func (r *RankingRepository) GetTop(ctx context.Context, pointTypeID int64, start, stop int64) ([]string, error) {
	vals, err := r.client.ZRevRange(ctx, key(pointTypeID), start, stop).Result()
	if err != nil {
		return nil, err
	}
	return vals, nil
}
