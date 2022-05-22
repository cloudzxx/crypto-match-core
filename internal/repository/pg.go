package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/cloudzxx/crypto-match-core/internal/config"
)

type PG struct {
	pool *pgxpool.Pool
}

func NewPG(cfg config.PGConfig) (*PG, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.MaxConns)

	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PG{pool: pool}, nil
}

func (p *PG) Close() {
	p.pool.Close()
}

func (p *PG) Pool() *pgxpool.Pool {
	return p.pool
}

func (p *PG) Exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := p.pool.Exec(ctx, sql, args...)
	return err
}

func (p *PG) Query(ctx context.Context, sql string, args ...interface{}) (*Rows, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{rows: rows, closable: rows}, nil
}

type Rows struct {
	rows    interface{ Next() bool; Scan(...interface{}) error }
	closable interface{ Close() }
}

func (r *Rows) Close() error {
	if r.closable != nil {
		r.closable.Close()
	}
	return nil
}

func (r *Rows) Next() bool {
	if r.rows == nil {
		return false
	}
	return r.rows.Next()
}

func (r *Rows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}