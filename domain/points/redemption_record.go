package points

// RedemptionStatus defines the status of a redemption record
type RedemptionStatus string

const (
	RedemptionCompleted RedemptionStatus = "completed"
	RedemptionPending   RedemptionStatus = "pending"
	RedemptionCancelled RedemptionStatus = "cancelled"
)

// RedemptionRecord represents a record of user redeeming rewards
type RedemptionRecord struct {
	ID        string           `json:"id"`
	UserID    string           `json:"userId"`
	RewardID  string           `json:"rewardId"`
	Costs     map[string]int64 `json:"costs"`
	CreatedAt int64            `json:"createdAt"`
	Status    RedemptionStatus `json:"status"`
}
