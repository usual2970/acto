package lib

import (
	restHandlers "acto/internal/rest/handlers"
	"acto/points"
)

// NewPointTypesHTTPHandler wraps existing point type handler using provided service.
func NewPointTypesHTTPHandler(s *Services) *restHandlers.PointTypesHandler {
	return restHandlers.NewPointTypesHandler(s.PointTypeService)
}

// NewBalancesHTTPHandler wraps existing balance handler using provided service.
func NewBalancesHTTPHandler(s *Services) *restHandlers.BalancesHandler {
	return restHandlers.NewBalancesHandler(s.BalanceService)
}

// Distributions
func NewDistributionsHTTPHandler(s *Services) *restHandlers.DistributionsHandler {
	return restHandlers.NewDistributionsHandler(s.DistributionService)
}

// Rankings
func NewRankingHTTPHandler(s *Services, svc points.RankingsService) *restHandlers.RankingsHandler {
	return restHandlers.NewRankingsHandler(svc)
}

// Redemptions
func NewRedemptionsHTTPHandler(s *Services) *restHandlers.RedemptionsHandler {
	return restHandlers.NewRedemptionsHandler(s.RedemptionService)
}

// Ensure imports kept if only used by adapter
var _ = points.TransactionFilter{}
