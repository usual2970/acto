package handlers

import (
	"net/http"

	d "github.com/usual2970/acto/domain/points"
)

func WriteDomainError(w http.ResponseWriter, err error) {
	// Map domain errors to numeric codes and human messages. HTTP status is
	// always 200 per API contract; the envelope's Code signals success/failure.
	switch err {
	case d.ErrInsufficientBalance:
		WriteError(w, 1001, "insufficient balance")
	case d.ErrPointTypeInUse:
		WriteError(w, 1002, "point type in use")
	case d.ErrRewardOutOfStock:
		WriteError(w, 1003, "reward out of stock")
	case d.ErrUnauthorizedOperation:
		WriteError(w, 1004, "forbidden")
	default:
		WriteError(w, 1500, "internal error")
	}
}
