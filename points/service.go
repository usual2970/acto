package points

import (
	"context"

	d "github.com/usual2970/acto/domain/points"
)

// Repository interfaces declared at consumer side (use case layer)
// so implementations can live in outer layers.

type PointTypeRepository interface {
	CreatePointType(ctx context.Context, pt d.PointType) (string, error)
	UpdatePointType(ctx context.Context, pt d.PointType) error
	DeletePointType(ctx context.Context, pointTypeID string) error
	SoftDeletePointType(ctx context.Context, name string) error
	GetPointTypeByID(ctx context.Context, pointTypeID string) (*d.PointType, error)
	GetPointTypeByName(ctx context.Context, name string) (*d.PointType, error)
	ListPointTypes(ctx context.Context, limit, offset int) ([]d.PointType, error)
	HasBalances(ctx context.Context, pointTypeID string) (bool, error)
}

type BalanceRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	GetUserBalanceForUpdate(ctx context.Context, userID, pointTypeID string) (*d.UserBalance, error)
	UpsertUserBalance(ctx context.Context, ub d.UserBalance) error
	InsertTransaction(ctx context.Context, tx d.Transaction) (string, error)
	ListTransactions(ctx context.Context, userID string, filter TransactionFilter) ([]d.Transaction, int, error)
}

type RankingRepository interface {
	UpdateUserScore(ctx context.Context, pointTypeID, userID string, score int64) error
	GetTop(ctx context.Context, pointTypeID string, start, stop int64) ([]string, error)
}

// RankingsService provides read-only ranking queries for delivery layer
type RankingsService interface {
	GetTop(ctx context.Context, pointTypeName string, limit, offset int) ([]string, error)
}

type RewardRepository interface {
	CreateRule(ctx context.Context, rr d.RewardRule) (string, error)
	ListRules(ctx context.Context, pointTypeID string) ([]d.RewardRule, error)
	CreateDistribution(ctx context.Context, rd d.RewardDistribution) (string, error)
	MarkDistributionCompleted(ctx context.Context, distributionID string) error
}

type RedemptionRepository interface {
	CreateReward(ctx context.Context, r d.RedemptionReward) (string, error)
	GetRewardByID(ctx context.Context, rewardID string) (*d.RedemptionReward, error)
	DecrementInventory(ctx context.Context, rewardID string, quantity int) error
	CreateRedemptionRecord(ctx context.Context, rr d.RedemptionRecord) (string, error)
}

// TransactionFilter defines optional filters and pagination for listing transactions
type TransactionFilter struct {
	PointTypeID   string
	OperationType string // "credit" | "debit" | ""
	StartTime     int64  // Unix timestamp or 0
	EndTime       int64  // Unix timestamp or 0
	Limit         int
	Offset        int
}
