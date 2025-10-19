package points

// RedemptionReward represents a reward that users can redeem with points
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
