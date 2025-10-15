# Data Model: Points System

## Entities

### PointType
- Fields: id (uuid), name (string, unique), displayName (string), description (string), enabled (bool), createdAt (timestamptz)
- Rules: name unique; cannot delete if balances exist; enable/disable supported

### UserBalance
- Fields: userId (string), pointTypeId (uuid), balance (int64), updatedAt (timestamptz)
- PK: (userId, pointTypeId)
- Rules: balance is integer, never negative; row locked on updates

### Transaction
- Fields: id (uuid), userId (string), pointTypeId (uuid), amount (int64, positive), type (enum: credit|debit), reason (string), before (int64), after (int64), createdAt (timestamptz)
- Rules: immutable; created for every credit/debit; amounts integers only

### RankingSnapshot
- Fields: id (uuid), pointTypeId (uuid), createdAt (timestamptz)
- Links: materialized from Redis when needed for distribution/audit

### RewardRule
- Fields: id (uuid), pointTypeId (uuid), minRank (int), maxRank (int), rewardAmount (int64), rewardPointTypeId (uuid), active (bool)

### RewardDistribution
- Fields: id (uuid), snapshotId (uuid), executedAt (timestamptz), status (enum: pending|completed|failed)
- Rules: unique per snapshot; links to generated transactions

### RedemptionReward
- Fields: id (uuid), name (string), description (string), costs (map<pointTypeId,int64>), quantity (int), enabled (bool), totalRedeemed (int), createdAt (timestamptz)

### RedemptionRecord
- Fields: id (uuid), userId (string), rewardId (uuid), costs (map<pointTypeId,int64>), createdAt (timestamptz), status (enum: completed|pending|cancelled)
- Rules: immutable once completed

## Relationships
- PointType 1..* UserBalance
- UserBalance 1..* Transaction
- PointType 1..* RewardRule
- RewardDistribution *..* Transaction (via generated bonuses)
- RedemptionReward 1..* RedemptionRecord

## Validation & Invariants
- All point values are integers; no decimals
- Deduction must not make balance negative
- Transaction history is append-only
- Ranking ties share rank; next rank skips accordingly

## State Transitions
- Credit: (before -> before+amount), record Transaction(credit)
- Debit: ensure before>=amount, (before -> before-amount), record Transaction(debit)
- Redemption: validate all required point types sufficient; deduct atomically
- Distribution: select top ranks per rule; create credit transactions atomically per user


