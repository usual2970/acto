package mysql

import (
	d "acto/domain/points"
	uc "acto/points"
	"context"
	"database/sql"
	"time"
)

type RewardsRepository struct{ db *sql.DB }

func NewRewardsRepository(db *sql.DB) *RewardsRepository { return &RewardsRepository{db: db} }

var _ uc.RewardRepository = (*RewardsRepository)(nil)

func (r *RewardsRepository) CreateRule(ctx context.Context, rr d.RewardRule) (string, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO reward_rules (id,point_type_id,min_rank,max_rank,reward_amount,reward_point_type_id,active) VALUES (UUID(),?,?,?,?,?,?)`, rr.PointTypeID, rr.MinRank, rr.MaxRank, rr.RewardAmount, rr.RewardPointTypeID, rr.Active)
	if err != nil {
		return "", err
	}
	var id string
	_ = r.db.QueryRowContext(ctx, `SELECT id FROM reward_rules WHERE point_type_id=? ORDER BY id DESC LIMIT 1`, rr.PointTypeID).Scan(&id)
	return id, nil
}

func (r *RewardsRepository) ListRules(ctx context.Context, pointTypeID string) ([]d.RewardRule, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,point_type_id,min_rank,max_rank,reward_amount,reward_point_type_id,active FROM reward_rules WHERE point_type_id=? AND active=1`, pointTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []d.RewardRule
	for rows.Next() {
		var rr d.RewardRule
		if err := rows.Scan(&rr.ID, &rr.PointTypeID, &rr.MinRank, &rr.MaxRank, &rr.RewardAmount, &rr.RewardPointTypeID, &rr.Active); err != nil {
			return nil, err
		}
		res = append(res, rr)
	}
	return res, rows.Err()
}

func (r *RewardsRepository) CreateDistribution(ctx context.Context, rd d.RewardDistribution) (string, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO reward_distributions (id,snapshot_id,status) VALUES (UUID(),?,?)`, rd.SnapshotID, rd.Status)
	if err != nil {
		return "", err
	}
	var id string
	_ = r.db.QueryRowContext(ctx, `SELECT id FROM reward_distributions ORDER BY id DESC LIMIT 1`).Scan(&id)
	return id, nil
}

func (r *RewardsRepository) MarkDistributionCompleted(ctx context.Context, distributionID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE reward_distributions SET status='completed', executed_at=? WHERE id=?`, time.Now().Unix(), distributionID)
	return err
}
