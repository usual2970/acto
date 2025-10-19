package points

import (
	"context"
	"errors"

	d "github.com/usual2970/acto/domain/points"
)

type RedemptionService struct {
	rewards RedemptionRepository
	balance BalanceRepository
}

func NewRedemptionService(rew RedemptionRepository, bal BalanceRepository) *RedemptionService {
	return &RedemptionService{rewards: rew, balance: bal}
}

// Redeem performs a redemption by deducting required point type costs and creating a record.
func (s *RedemptionService) Redeem(ctx context.Context, req RedemptionRequest) error {
	reward, err := s.rewards.GetRewardByID(ctx, req.RewardID)
	if err != nil {
		return err
	}
	if !reward.Enabled {
		return d.ErrUnauthorizedOperation
	}
	if len(reward.Costs) == 0 {
		return nil
	}

	return s.balance.WithTx(ctx, func(ctx context.Context) error {
		// Check all balances first
		for ptID, cost := range reward.Costs {
			ub, err := s.balance.GetUserBalanceForUpdate(ctx, req.UserID, ptID)
			if err != nil {
				return err
			}
			if ub.Balance < cost {
				return d.ErrInsufficientBalance
			}
		}
		// Deduct for each cost
		for ptID, cost := range reward.Costs {
			ub, err := s.balance.GetUserBalanceForUpdate(ctx, req.UserID, ptID)
			if err != nil {
				return err
			}
			before := ub.Balance
			ub.Balance = before - cost
			if err := s.balance.UpsertUserBalance(ctx, *ub); err != nil {
				return err
			}
			if _, err := s.balance.InsertTransaction(ctx, d.Transaction{UserID: req.UserID, PointTypeID: ptID, Amount: cost, Type: d.TransactionDebit, Reason: "redemption", Before: before, After: ub.Balance}); err != nil {
				return err
			}
		}
		// Create record (note: inventory decrement not in same tx due to different repo)
		if err := s.rewards.DecrementInventory(ctx, req.RewardID, 1); err != nil {
			return err
		}
		if _, err := s.rewards.CreateRedemptionRecord(ctx, d.RedemptionRecord{UserID: req.UserID, RewardID: req.RewardID, Status: d.RedemptionCompleted}); err != nil {
			return err
		}
		return nil
	})
}

var ErrInvalidRequest = errors.New("invalid redemption request")
