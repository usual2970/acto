package points

// RewardDistributionStatus defines the status of reward distribution
type RewardDistributionStatus string

const (
	DistributionPending   RewardDistributionStatus = "pending"
	DistributionCompleted RewardDistributionStatus = "completed"
	DistributionFailed    RewardDistributionStatus = "failed"
)

// RewardDistribution represents a reward distribution execution record
type RewardDistribution struct {
	ID         string                   `json:"id"`
	SnapshotID string                   `json:"snapshotId"`
	ExecutedAt int64                    `json:"executedAt"`
	Status     RewardDistributionStatus `json:"status"`
}
