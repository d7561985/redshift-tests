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
	Bulk(context.Context, []*postgres.Journal) error
}

type Config struct {
	Size    int
	MaxUser int
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
			default:
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
		for journal := range ch {
			batch = append(batch, &journal)

			if len(batch) < c.Size {
				continue
			}

			t := time.Now()
			ms := t.Sub(start)

			if err := c.repo.Bulk(context.TODO(), batch); err != nil {
				l.Fatal("bulk", tel.Error(errors.WithStack(err)))
			}
			eTime := time.Now().Sub(t)

			batch = batch[:0]
			out += c.Size

			l.Info("progress",
				tel.Float64("comb/sec", float64(out)/ms.Seconds()),
				tel.Duration("duration", ms),
				tel.Int("out", out),
				tel.Duration("last", eTime),
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
