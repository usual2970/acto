package points

import "time"

type PointType struct {
	ID          string
	Name        string
	DisplayName string
	Description string
	Enabled     bool
	CreatedAt   time.Time
}

type UserBalance struct {
	UserID      string
	PointTypeID string
	Balance     int64
	UpdatedAt   time.Time
}

type TransactionType string

const (
	TransactionCredit TransactionType = "credit"
	TransactionDebit  TransactionType = "debit"
)

type Transaction struct {
	ID          string
	UserID      string
	PointTypeID string
	Amount      int64
	Type        TransactionType
	Reason      string
	Before      int64
	After       int64
	CreatedAt   time.Time
}

type RewardRule struct {
	ID                string
	PointTypeID       string
	MinRank           int
	MaxRank           int
	RewardAmount      int64
	RewardPointTypeID string
	Active            bool
}

type RewardDistributionStatus string

const (
	DistributionPending   RewardDistributionStatus = "pending"
	DistributionCompleted RewardDistributionStatus = "completed"
	DistributionFailed    RewardDistributionStatus = "failed"
)

type RewardDistribution struct {
	ID         string
	SnapshotID string
	ExecutedAt time.Time
	Status     RewardDistributionStatus
}

type RedemptionReward struct {
	ID            string
	Name          string
	Description   string
	Costs         map[string]int64 // pointTypeID -> amount
	Quantity      int
	Enabled       bool
	TotalRedeemed int
	CreatedAt     time.Time
}

type RedemptionStatus string

const (
	RedemptionCompleted RedemptionStatus = "completed"
	RedemptionPending   RedemptionStatus = "pending"
	RedemptionCancelled RedemptionStatus = "cancelled"
)

type RedemptionRecord struct {
	ID        string
	UserID    string
	RewardID  string
	Costs     map[string]int64
	CreatedAt time.Time
	Status    RedemptionStatus
}
