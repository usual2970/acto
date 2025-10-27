package points

// PointType represents a type of points in the system
type PointType struct {
	ID          int64  `json:"id"`
	URI         string `json:"uri"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	DeletedAt   *int64 `json:"deletedAt,omitempty"`
	CreatedAt   int64  `json:"createdAt"`
}
