package points

import (
	d "acto/domain/points"
	"context"
)

type DistributionService struct {
	rewards RewardRepository
	balance BalanceRepository
	ranking RankingRepository
}

func NewDistributionService(rew RewardRepository, bal BalanceRepository, rank RankingRepository) *DistributionService {
	return &DistributionService{rewards: rew, balance: bal, ranking: rank}
}

// Execute runs a distribution for a point type using current ranking top N and active rules.
func (s *DistributionService) Execute(ctx context.Context, pointTypeID string, topN int) error {
	rules, err := s.rewards.ListRules(ctx, pointTypeID)
	if err != nil {
		return err
	}
	if len(rules) == 0 {
		return nil
	}
	users, err := s.ranking.GetTop(ctx, pointTypeID, 0, int64(topN-1))
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
