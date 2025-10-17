package handlers

import (
	"encoding/json"
	"net/http"

	d "github.com/usual2970/acto/domain/points"
)

type errorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func writeDomainError(w http.ResponseWriter, err error, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	var status int
	var code string
	switch err {
	case d.ErrInsufficientBalance:
		status, code = http.StatusConflict, "INSUFFICIENT_BALANCE"
	case d.ErrPointTypeInUse:
		status, code = http.StatusConflict, "CANNOT_DELETE_ACTIVE_POINT_TYPE"
	case d.ErrRewardOutOfStock:
		status, code = http.StatusConflict, "REWARD_OUT_OF_STOCK"
	case d.ErrUnauthorizedOperation:
		status, code = http.StatusForbidden, "FORBIDDEN"
	default:
		status, code = http.StatusInternalServerError, "INTERNAL_ERROR"
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Code: code, Message: err.Error(), Details: details})
}
