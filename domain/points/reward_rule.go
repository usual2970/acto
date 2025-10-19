package points

// RewardRule defines rules for distributing rewards based on ranking
type RewardRule struct {
	ID                string `json:"id"`
	PointTypeID       string `json:"pointTypeId"`
	MinRank           int    `json:"minRank"`
	MaxRank           int    `json:"maxRank"`
	RewardAmount      int64  `json:"rewardAmount"`
	RewardPointTypeID string `json:"rewardPointTypeId"`
	Active            bool   `json:"active"`
}
