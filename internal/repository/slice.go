package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudzxx/crypto-match-core/internal/config"
	"github.com/cloudzxx/crypto-match-core/internal/service"
)

type SliceManager struct {
	pg       *PG
	interval int
	keeptime int
}

func NewSliceManager(pg *PG, cfg config.SliceConfig) *SliceManager {
	return &SliceManager{
		pg:       pg,
		interval: cfg.Interval,
		keeptime: cfg.KeepTime,
	}
}

func (sm *SliceManager) DoSlice(mkts []*service.Market, bm *service.BalanceManager) (string, error) {
	ctx := context.Background()

	sliceID := fmt.Sprintf("slice_%d", time.Now().Unix())
	now := time.Now()

	tx, err := sm.pg.Pool().Begin(ctx)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO slice_history (slice_id, begin_time, end_time)
		VALUES ($1, $2, $2)
	`, sliceID, now)
	if err != nil {
		tx.Rollback(ctx)
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return sliceID, nil
}

func (sm *SliceManager) LoadLatestSlice(ctx context.Context) (string, error) {
	var sliceID string
	err := sm.pg.Pool().QueryRow(ctx, `
		SELECT slice_id FROM slice_history
		ORDER BY created_at DESC LIMIT 1
	`).Scan(&sliceID)
	if err != nil {
		return "", err
	}
	return sliceID, nil
}

func (sm *SliceManager) CleanupOldSlices(ctx context.Context) error {
	cutoff := time.Now().Add(-time.Duration(sm.keeptime) * time.Second)
	err := sm.pg.Exec(ctx, `
		DELETE FROM slice_history WHERE created_at < $1
	`, cutoff)
	return err
}

func (sm *SliceManager) Interval() int {
	return sm.interval
}