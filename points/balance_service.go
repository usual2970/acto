package points

import (
	d "acto/domain/points"
	"context"
)

type BalanceService struct {
	repo       BalanceRepository
	ranking    RankingRepository
	pointTypes PointTypeRepository
}

func NewBalanceService(repo BalanceRepository, ranking RankingRepository, pts PointTypeRepository) *BalanceService {
	return &BalanceService{repo: repo, ranking: ranking, pointTypes: pts}
}

func (s *BalanceService) Credit(ctx context.Context, userID, pointTypeName, reason string, amount int64) error {
	if amount <= 0 {
		return nil
	}
	pt, err := s.pointTypes.GetPointTypeByName(ctx, pointTypeName)
	if err != nil {
		return err
	}
	return s.repo.WithTx(ctx, func(ctx context.Context) error {
		ub, err := s.repo.GetUserBalanceForUpdate(ctx, userID, pt.ID)
		if err != nil {
			return err
		}
		before := ub.Balance
		ub.Balance += amount
		if err := s.repo.UpsertUserBalance(ctx, *ub); err != nil {
			return err
		}
		_, err = s.repo.InsertTransaction(ctx, d.Transaction{UserID: userID, PointTypeID: pt.ID, Amount: amount, Type: d.TransactionCredit, Reason: reason, Before: before, After: ub.Balance})
		if err != nil {
			return err
		}
		if s.ranking != nil {
			_ = s.ranking.UpdateUserScore(ctx, pt.ID, userID, ub.Balance)
		}
		return nil
	})
}

func (s *BalanceService) Debit(ctx context.Context, userID, pointTypeName, reason string, amount int64) error {
	if amount <= 0 {
		return nil
	}
	pt, err := s.pointTypes.GetPointTypeByName(ctx, pointTypeName)
	if err != nil {
		return err
	}
	return s.repo.WithTx(ctx, func(ctx context.Context) error {
		ub, err := s.repo.GetUserBalanceForUpdate(ctx, userID, pt.ID)
		if err != nil {
			return err
		}
		if ub.Balance < amount {
			return d.ErrInsufficientBalance
		}
		before := ub.Balance
		ub.Balance -= amount
		if err := s.repo.UpsertUserBalance(ctx, *ub); err != nil {
			return err
		}
		_, err = s.repo.InsertTransaction(ctx, d.Transaction{UserID: userID, PointTypeID: pt.ID, Amount: amount, Type: d.TransactionDebit, Reason: reason, Before: before, After: ub.Balance})
		if err != nil {
			return err
		}
		if s.ranking != nil {
			_ = s.ranking.UpdateUserScore(ctx, pt.ID, userID, ub.Balance)
		}
		return nil
	})
}

// ListTransactions returns transactions for a user with optional filters
func (s *BalanceService) ListTransactions(ctx context.Context, userID, pointTypeName, op, startISO, endISO string, limit, offset int) ([]d.Transaction, int, error) {
	var pointTypeID string
	if pointTypeName != "" {
		pt, err := s.pointTypes.GetPointTypeByName(ctx, pointTypeName)
		if err != nil {
			return nil, 0, err
		}
		pointTypeID = pt.ID
	}
	filter := TransactionFilter{PointTypeID: pointTypeID, OperationType: op, StartTimeISO: startISO, EndTimeISO: endISO, Limit: limit, Offset: offset}
	return s.repo.ListTransactions(ctx, userID, filter)
}
