package points

import (
	"context"

	d "github.com/usual2970/acto/domain/points"
)

type BalanceService struct {
	repo       BalanceRepository
	ranking    RankingRepository
	pointTypes PointTypeRepository
}

func NewBalanceService(repo BalanceRepository, ranking RankingRepository, pts PointTypeRepository) *BalanceService {
	return &BalanceService{repo: repo, ranking: ranking, pointTypes: pts}
}

func (s *BalanceService) Credit(ctx context.Context, req BalanceCreditRequest) error {
	if req.Amount <= 0 {
		return nil
	}
	pt, err := s.pointTypes.GetPointTypeByName(ctx, req.PointTypeName)
	if err != nil {
		return err
	}
	return s.repo.WithTx(ctx, func(ctx context.Context) error {
		ub, err := s.repo.GetUserBalanceForUpdate(ctx, req.UserID, pt.ID)
		if err != nil {
			return err
		}
		before := ub.Balance
		ub.Balance += req.Amount
		if err := s.repo.UpsertUserBalance(ctx, *ub); err != nil {
			return err
		}
		_, err = s.repo.InsertTransaction(ctx, d.Transaction{UserID: req.UserID, PointTypeID: pt.ID, Amount: req.Amount, Type: d.TransactionCredit, Reason: req.Reason, Before: before, After: ub.Balance})
		if err != nil {
			return err
		}
		if s.ranking != nil {
			_ = s.ranking.UpdateUserScore(ctx, pt.ID, req.UserID, ub.Balance)
		}
		return nil
	})
}

func (s *BalanceService) Debit(ctx context.Context, req BalanceDebitRequest) error {
	if req.Amount <= 0 {
		return nil
	}
	pt, err := s.pointTypes.GetPointTypeByName(ctx, req.PointTypeName)
	if err != nil {
		return err
	}
	return s.repo.WithTx(ctx, func(ctx context.Context) error {
		ub, err := s.repo.GetUserBalanceForUpdate(ctx, req.UserID, pt.ID)
		if err != nil {
			return err
		}
		if ub.Balance < req.Amount {
			return d.ErrInsufficientBalance
		}
		before := ub.Balance
		ub.Balance -= req.Amount
		if err := s.repo.UpsertUserBalance(ctx, *ub); err != nil {
			return err
		}
		_, err = s.repo.InsertTransaction(ctx, d.Transaction{UserID: req.UserID, PointTypeID: pt.ID, Amount: req.Amount, Type: d.TransactionDebit, Reason: req.Reason, Before: before, After: ub.Balance})
		if err != nil {
			return err
		}
		if s.ranking != nil {
			_ = s.ranking.UpdateUserScore(ctx, pt.ID, req.UserID, ub.Balance)
		}
		return nil
	})
}

// ListTransactions returns transactions for a user with optional filters
func (s *BalanceService) ListTransactions(ctx context.Context, userID, pointTypeName, op string, startTime, endTime int64, limit, offset int) ([]d.Transaction, int, error) {
	var pointTypeID string
	if pointTypeName != "" {
		pt, err := s.pointTypes.GetPointTypeByName(ctx, pointTypeName)
		if err != nil {
			return nil, 0, err
		}
		pointTypeID = pt.ID
	}
	filter := TransactionFilter{PointTypeID: pointTypeID, OperationType: op, StartTime: startTime, EndTime: endTime, Limit: limit, Offset: offset}
	return s.repo.ListTransactions(ctx, userID, filter)
}
