package mysql

import (
	d "acto/domain/points"
	uc "acto/points"
	"context"
	"database/sql"
	"fmt"
)

type BalanceTxRepository struct {
	db *sql.DB
}

func NewBalanceTxRepository(db *sql.DB) *BalanceTxRepository { return &BalanceTxRepository{db: db} }

var _ uc.BalanceRepository = (*BalanceTxRepository)(nil)

func (r *BalanceTxRepository) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)
	if err := fn(ctxWithTx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

type txKey struct{}

func getTx(ctx context.Context, db *sql.DB) executor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok && tx != nil {
		return tx
	}
	return db
}

type executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *BalanceTxRepository) GetUserBalanceForUpdate(ctx context.Context, userID, pointTypeID string) (*d.UserBalance, error) {
	ex := getTx(ctx, r.db)
	row := ex.QueryRowContext(ctx, `SELECT user_id, point_type_id, balance, updated_at FROM user_balances WHERE user_id=? AND point_type_id=? FOR UPDATE`, userID, pointTypeID)
	var ub d.UserBalance
	if err := row.Scan(&ub.UserID, &ub.PointTypeID, &ub.Balance, &ub.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return &d.UserBalance{UserID: userID, PointTypeID: pointTypeID, Balance: 0}, nil
		}
		return nil, err
	}
	return &ub, nil
}

func (r *BalanceTxRepository) UpsertUserBalance(ctx context.Context, ub d.UserBalance) error {
	ex := getTx(ctx, r.db)
	_, err := ex.ExecContext(ctx, `INSERT INTO user_balances (user_id, point_type_id, balance) VALUES (?,?,?) ON DUPLICATE KEY UPDATE balance=VALUES(balance)`, ub.UserID, ub.PointTypeID, ub.Balance)
	return err
}

func (r *BalanceTxRepository) InsertTransaction(ctx context.Context, tx d.Transaction) (string, error) {
	ex := getTx(ctx, r.db)
	_, err := ex.ExecContext(ctx, `INSERT INTO transactions (id,user_id,point_type_id,amount,type,reason,before_balance,after_balance) VALUES (UUID(),?,?,?,?,?,?,?)`, tx.UserID, tx.PointTypeID, tx.Amount, string(tx.Type), tx.Reason, tx.Before, tx.After)
	if err != nil {
		return "", err
	}
	var id string
	_ = r.db.QueryRowContext(ctx, `SELECT id FROM transactions WHERE user_id=? ORDER BY created_at DESC LIMIT 1`, tx.UserID).Scan(&id)
	return id, nil
}

func (r *BalanceTxRepository) ListTransactions(ctx context.Context, userID string, filter uc.TransactionFilter) ([]d.Transaction, int, error) {
	where := "WHERE user_id=?"
	args := []any{userID}
	if filter.PointTypeID != "" {
		where += " AND point_type_id=?"
		args = append(args, filter.PointTypeID)
	}
	if filter.OperationType == "credit" || filter.OperationType == "debit" {
		where += " AND type=?"
		args = append(args, filter.OperationType)
	}
	if filter.StartTimeISO != "" {
		where += " AND created_at>=?"
		args = append(args, filter.StartTimeISO)
	}
	if filter.EndTimeISO != "" {
		where += " AND created_at<?"
		args = append(args, filter.EndTimeISO)
	}

	// total count
	var total int
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(1) FROM transactions %s", where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// pagination
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT id,user_id,point_type_id,amount,type,reason,before_balance,after_balance,created_at FROM transactions %s ORDER BY created_at DESC LIMIT ? OFFSET ?", where), append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var res []d.Transaction
	for rows.Next() {
		var t d.Transaction
		var typ string
		if err := rows.Scan(&t.ID, &t.UserID, &t.PointTypeID, &t.Amount, &typ, &t.Reason, &t.Before, &t.After, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		if typ == "credit" {
			t.Type = d.TransactionCredit
		} else {
			t.Type = d.TransactionDebit
		}
		res = append(res, t)
	}
	return res, total, rows.Err()
}
