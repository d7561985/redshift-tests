package postgres

import (
	"context"

	"github.com/d7561985/redshift-test/internal/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const maxCon = 100

type Repo struct {
	cfg config.Postgres

	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Postgres) (*Repo, error) {
	c, err := pgxpool.ParseConfig(cfg.Addr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.MaxConns = maxCon

	dbpool, err := pgxpool.ConnectConfig(ctx, c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err = dbpool.Ping(ctx); err != nil {
		return nil, errors.WithStack(err)
	}

	return &Repo{
		cfg:  cfg,
		pool: dbpool,
	}, nil
}

func (s *Repo) Insert(ctx context.Context, j Journal) error {
	sq := `INSERT INTO journal("id","transactionId", "accountId", "balance", "change","currency","created_at",
"pincoinBalance","pincoinChange","project","revert","type"
                    ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := s.pool.Exec(ctx, sq,
		j.ID, j.TransactionID, j.AccountID, j.Balance, j.Change, j.Currency, j.CreatedAt,
		j.PincoinBalance, j.PincoinChange, j.Project, j.Revert, j.Type,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil

}

func (s *Repo) Bulk(ctx context.Context, list []*Journal) (string, error) {
	sq := `INSERT INTO journal("id","transactionId", "accountId", "balance", "change","currency","created_at",
"pincoinBalance","pincoinChange","project","revert","type"
                    ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	batch := &pgx.Batch{}

	for _, j := range list {
		batch.Queue(sq, j.ID, j.TransactionID, j.AccountID, j.Balance, j.Change, j.Currency, j.CreatedAt,
			j.PincoinBalance, j.PincoinChange, j.Project, j.Revert, j.Type)
	}

	br := s.pool.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return "direct", nil
}
