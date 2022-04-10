package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/d7561985/mongo-ab/pkg/changing"
	"github.com/d7561985/redshift-test/store/postgres"
	"github.com/d7561985/tel/v2"
	fuzz "github.com/google/gofuzz"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo interface {
	Bulk(context.Context, []*postgres.Journal) (string, error)
}

type Config struct {
	Size int
	//
	MaxBatch  int
	TimeLimit time.Duration
	MaxUser   int
}

type controller struct {
	Config
	repo Repo
}

func New(cfg Config, repo Repo) *controller {
	return &controller{
		Config: cfg,
		repo:   repo,
	}
}

func (c *controller) Run(ctx context.Context) {
	l := tel.FromCtx(ctx)
	ch := make(chan postgres.Journal, c.Size)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-time.After(time.Microsecond * 5):
				//default:
			}

			tx := genRequest(uint64(rand.Int()%c.MaxUser), 100)
			ch <- postgres.NewJournal(postgres.Balance{
				AccountID:      tx.AccountID,
				Balance:        float64(rand.Int63() % 10000),
				PincoinBalance: float64(rand.Int63() % 10000),
			}, tx)
		}
	}()

	go func() {
		batch := make([]*postgres.Journal, 0, c.Size)
		var out int

		start := time.Now()
		last := start

		for journal := range ch {
			j := journal
			batch = append(batch, &j)

			if c.MaxBatch >= 0 && len(batch) < c.MaxBatch &&
				time.Since(last) < c.TimeLimit {
				continue
			}

			t := time.Now()
			ms := t.Sub(start)

			f, err := c.repo.Bulk(context.TODO(), batch)
			if err != nil {
				l.Fatal("bulk", tel.Error(errors.WithStack(err)))
			}

			last = time.Now()
			eTime := last.Sub(t)

			out += len(batch)
			batch = batch[:0]

			l.Info("progress",
				tel.Float64("comb/sec", float64(out)/ms.Seconds()),
				tel.Duration("duration", ms),
				tel.Int("out", out),
				tel.Duration("last", eTime),
				tel.String("file", f),
			)
		}
	}()

	<-ctx.Done()
}

func genRequest(usr uint64, add float64) changing.Transaction {
	tx := changing.Transaction{}
	fuzz.New().Fuzz(&tx)

	tx.Inc = changing.Inc{
		Balance:        add,
		DepositAllSum:  100,
		DepositCount:   1,
		PincoinBalance: 100,
		PincoinsAllSum: 1,
	}

	tx.AccountID = usr
	tx.Currency = 123
	tx.Change = add

	switch rand.Int63() % 2 {
	case 0:
		tx.TransactionID = uint64(rand.Int63())
	case 1:
		tx.TransactionIDBson = primitive.NewObjectID()
		tx.TransactionID = 0
	}

	tx.Project = []string{"casino", "sport", primitive.NewObjectID().String()}[rand.Int63()%3]
	tx.Type = []string{
		"None", "Add Deposit", "Write bet", "FreebetWin", "Withdraw", "LotteryWin", "Welcome deposit", "Revert",
	}[rand.Int63()%8]

	return tx
}
