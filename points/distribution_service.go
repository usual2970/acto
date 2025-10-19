package points

import (
	"context"

	d "github.com/usual2970/acto/domain/points"
)

type DistributionService struct {
	rewards RewardRepository
	balance BalanceRepository
	ranking RankingRepository
	points  PointTypeRepository
}

func NewDistributionService(rew RewardRepository, bal BalanceRepository, rank RankingRepository, pts PointTypeRepository) *DistributionService {
	return &DistributionService{rewards: rew, balance: bal, ranking: rank, points: pts}
}

// rankingsService implements RankingsService using repositories.
type rankingsService struct {
	ranking RankingRepository
	points  PointTypeRepository
}

func NewRankingsService(rank RankingRepository, pts PointTypeRepository) RankingsService {
	return &rankingsService{ranking: rank, points: pts}
}

func (s *rankingsService) GetTop(ctx context.Context, pointTypeName string, limit, offset int) ([]string, error) {
	var ptID string
	if pointTypeName != "" && s.points != nil {
		if pt, err := s.points.GetPointTypeByName(ctx, pointTypeName); err == nil && pt != nil {
			ptID = pt.ID
		}
	}
	if limit <= 0 {
		limit = 100
	}
	start := int64(offset)
	stop := int64(offset + limit - 1)
	return s.ranking.GetTop(ctx, ptID, start, stop)
}

// Execute runs a distribution for a point type using current ranking top N and active rules.
func (s *DistributionService) Execute(ctx context.Context, req DistirbutionsExecuteRequest) error {
	pt, err := s.points.GetPointTypeByName(ctx, req.PointTypeName)
	if err != nil {
		return err
	}
	pointTypeID := pt.ID
	rules, err := s.rewards.ListRules(ctx, pointTypeID)
	if err != nil {
		return err
	}
	if len(rules) == 0 {
		return nil
	}
	users, err := s.ranking.GetTop(ctx, pointTypeID, 0, int64(req.TopN-1))
	if err != nil {
		return err
	}
	distID, err := s.rewards.CreateDistribution(ctx, d.RewardDistribution{SnapshotID: "now", Status: d.DistributionPending})
	if err != nil {
		return err
	}
	// naive application: apply rewards by index rank (1-based)
	rank := 1
	for _, user := range users {
		for _, rule := range rules {
			if rank >= rule.MinRank && rank <= rule.MaxRank {
				// credit reward to user in rule.RewardPointTypeID
				_ = s.balance.WithTx(ctx, func(ctx context.Context) error {
					ub, err := s.balance.GetUserBalanceForUpdate(ctx, user, rule.RewardPointTypeID)
					if err != nil {
						return err
					}
					before := ub.Balance
					ub.Balance += rule.RewardAmount
					if err := s.balance.UpsertUserBalance(ctx, *ub); err != nil {
						return err
					}
					_, err = s.balance.InsertTransaction(ctx, d.Transaction{UserID: user, PointTypeID: rule.RewardPointTypeID, Amount: rule.RewardAmount, Type: d.TransactionCredit, Reason: "rank reward", Before: before, After: ub.Balance})
					return err
				})
				break
			}
		}
		rank++
	}
	return s.rewards.MarkDistributionCompleted(ctx, distID)
}
