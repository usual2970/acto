package points

import "errors"

// Domain-level errors for points system. These errors are intended to be
// checked by upper layers and mapped to transport-specific error responses.

var (
	ErrInsufficientBalance     = errors.New("insufficient balance")
	ErrPointTypeInUse          = errors.New("cannot delete point type with active balances")
	ErrDuplicatePointTypeName  = errors.New("duplicate point type name")
	ErrPointTypeNotFound       = errors.New("point type not found")
	ErrPointTypeAlreadyDeleted = errors.New("point type already deleted")
	ErrRewardOutOfStock        = errors.New("reward out of stock")
	ErrDistributionAlreadyDone = errors.New("distribution already executed for period")
	ErrUnauthorizedOperation   = errors.New("unauthorized operation for role")
)
