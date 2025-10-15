package mysql

import (
	d "acto/domain/points"
	uc "acto/points"
	"context"
	"database/sql"
)

type RedemptionRepository struct{ db *sql.DB }

func NewRedemptionRepository(db *sql.DB) *RedemptionRepository { return &RedemptionRepository{db: db} }

var _ uc.RedemptionRepository = (*RedemptionRepository)(nil)

func (r *RedemptionRepository) CreateReward(ctx context.Context, rr d.RedemptionReward) (string, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO redemption_rewards (id,name,description,quantity,enabled,total_redeemed) VALUES (UUID(),?,?,?,?,0)`, rr.Name, rr.Description, rr.Quantity, rr.Enabled)
	if err != nil {
		return "", err
	}
	var id string
	_ = r.db.QueryRowContext(ctx, `SELECT id FROM redemption_rewards WHERE name=? ORDER BY created_at DESC LIMIT 1`, rr.Name).Scan(&id)
	for pt, amt := range rr.Costs {
		if _, err := r.db.ExecContext(ctx, `INSERT INTO redemption_costs (reward_id,point_type_id,amount) VALUES (?,?,?)`, id, pt, amt); err != nil {
			return "", err
		}
	}
	return id, nil
}

func (r *RedemptionRepository) GetRewardByID(ctx context.Context, rewardID string) (*d.RedemptionReward, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id,name,description,quantity,enabled,total_redeemed,created_at FROM redemption_rewards WHERE id=?`, rewardID)
	var rr d.RedemptionReward
	if err := row.Scan(&rr.ID, &rr.Name, &rr.Description, &rr.Quantity, &rr.Enabled, &rr.TotalRedeemed, &rr.CreatedAt); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT point_type_id,amount FROM redemption_costs WHERE reward_id=?`, rewardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rr.Costs = map[string]int64{}
	for rows.Next() {
		var pt string
		var amt int64
		if err := rows.Scan(&pt, &amt); err != nil {
			return nil, err
		}
		rr.Costs[pt] = amt
	}
	return &rr, rows.Err()
}

func (r *RedemptionRepository) DecrementInventory(ctx context.Context, rewardID string, quantity int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE redemption_rewards SET quantity=quantity-? WHERE id=? AND quantity>=?`, quantity, rewardID, quantity)
	return err
}

func (r *RedemptionRepository) CreateRedemptionRecord(ctx context.Context, rec d.RedemptionRecord) (string, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO redemption_records (id,user_id,reward_id,status) VALUES (UUID(),?,?,'completed')`, rec.UserID, rec.RewardID)
	if err != nil {
		return "", err
	}
	var id string
	_ = r.db.QueryRowContext(ctx, `SELECT id FROM redemption_records WHERE user_id=? ORDER BY created_at DESC LIMIT 1`, rec.UserID).Scan(&id)
	return id, nil
}
