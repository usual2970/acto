package points

// TransactionType defines the type of point transaction
type TransactionType string

const (
	TransactionCredit TransactionType = "credit"
	TransactionDebit  TransactionType = "debit"
)

// Transaction represents a point transaction record
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
