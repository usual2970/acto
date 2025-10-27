package points

type BalanceCreditRequest struct {
	UserID string `json:"userId"`
	URI    string `json:"uri"`
	Reason string `json:"reason"`
	Amount int64  `json:"amount"`
}

type BalanceDebitRequest struct {
	UserID string `json:"userId"`
	URI    string `json:"uri"`
	Reason string `json:"reason"`
	Amount int64  `json:"amount"`
}

type DistirbutionsExecuteRequest struct {
	URI  string `json:"uri"`
	TopN int    `json:"topN"`
}

type PointTypeCreateRequest struct {
	URI         string `json:"uri"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

// PointTypeUpdateRequest represents the request for updating a point type
type PointTypeUpdateRequest struct {
	DisplayName *string `json:"displayName,omitempty"`
	Description *string `json:"description,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
}

type RedemptionRequest struct {
	UserID   string `json:"userId"`
	RewardID string `json:"rewardId"`
}
