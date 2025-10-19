package points

// UserBalance represents the balance of a specific point type for a user
type UserBalance struct {
	UserID      string `json:"userId"`
	PointTypeID string `json:"pointTypeId"`
	Balance     int64  `json:"balance"`
	UpdatedAt   int64  `json:"updatedAt"`
}
