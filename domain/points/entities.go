package points

// Time fields are now represented as Unix timestamps (int64)

type PointType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	DeletedAt   *int64 `json:"deletedAt,omitempty"`
	CreatedAt   int64  `json:"createdAt"`
}

type UserBalance struct {
	UserID      string `json:"userId"`
	PointTypeID string `json:"pointTypeId"`
	Balance     int64  `json:"balance"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type TransactionType string

const (
	TransactionCredit TransactionType = "credit"
	TransactionDebit  TransactionType = "debit"
)

type Transaction struct {
	ID          string          `json:"id"`
	UserID      string          `json:"userId"`
	PointTypeID string          `json:"pointTypeId"`
	Amount      int64           `json:"amount"`
	Type        TransactionType `json:"type"`
	Reason      string          `json:"reason"`
	Before      int64           `json:"before"`
	After       int64           `json:"after"`
	CreatedAt   int64           `json:"createdAt"`
}

type RewardRule struct {
	ID                string `json:"id"`
	PointTypeID       string `json:"pointTypeId"`
	MinRank           int    `json:"minRank"`
	MaxRank           int    `json:"maxRank"`
	RewardAmount      int64  `json:"rewardAmount"`
	RewardPointTypeID string `json:"rewardPointTypeId"`
	Active            bool   `json:"active"`
}

type RewardDistributionStatus string

const (
	DistributionPending   RewardDistributionStatus = "pending"
	DistributionCompleted RewardDistributionStatus = "completed"
	DistributionFailed    RewardDistributionStatus = "failed"
)

type RewardDistribution struct {
	ID         string                   `json:"id"`
	SnapshotID string                   `json:"snapshotId"`
	ExecutedAt int64                    `json:"executedAt"`
	Status     RewardDistributionStatus `json:"status"`
}

type RedemptionReward struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Costs         map[string]int64 `json:"costs"` // pointTypeID -> amount
	Quantity      int              `json:"quantity"`
	Enabled       bool             `json:"enabled"`
	TotalRedeemed int              `json:"totalRedeemed"`
	CreatedAt     int64            `json:"createdAt"`
}

type RedemptionStatus string

const (
	RedemptionCompleted RedemptionStatus = "completed"
	RedemptionPending   RedemptionStatus = "pending"
	RedemptionCancelled RedemptionStatus = "cancelled"
)

type RedemptionRecord struct {
	ID        string           `json:"id"`
	UserID    string           `json:"userId"`
	RewardID  string           `json:"rewardId"`
	Costs     map[string]int64 `json:"costs"`
	CreatedAt int64            `json:"createdAt"`
	Status    RedemptionStatus `json:"status"`
}
