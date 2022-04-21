package service

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/d7561985/mongo-ab/pkg/changing"
	"github.com/d7561985/redshift-test/model"
	"github.com/d7561985/tel/v2"
	fuzz "github.com/google/gofuzz"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo interface {
	PlayerInsert(ctx context.Context, p []*model.Player) (string, error)
	CasinoBetInsert(ctx context.Context, p []*model.CBet) (string, error)
}

type Config struct {
	PlayerMaxID uint64
	PlayerRate  int
	CBRate      int

	// WindowTime window making operation
	WindowTime time.Duration
}

type controller struct {
	Config
	repo Repo
}

func New(cfg Config, repo Repo) *controller {
	if cfg.WindowTime == 0 ||
		cfg.PlayerRate == 0 ||
		cfg.CBRate == 0 {
		tel.Global().Fatal("cfg validation", tel.Any("cfg", cfg))
	}

	return &controller{
		Config: cfg,
		repo:   repo,
	}
}

func (c *controller) Run(ctx context.Context) {
	go c.DoPlayer(ctx, c.PlayerRate)
	go c.DoCB(ctx, c.CBRate)

	<-ctx.Done()
}

func (c *controller) DoPlayer(ctx context.Context, rate int) {
	l := tel.FromCtx(ctx)

	start := time.Now()
	last := start
	var out int

	for {
		sleepTime := c.WindowTime - time.Now().Sub(last)
		if sleepTime < 0 {
			sleepTime = 0
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleepTime):
		}

		secs := time.Now().Sub(last)
		num := int(float64(rate) * secs.Minutes())
		if num < 1 {
			continue
		}

		out += num

		last = time.Now()

		var p []*model.Player
		for i := 0; i < num; i++ {
			v := model.NewPlayer(c.PlayerMaxID + uint64(i) + 1)
			v.OK()
			p = append(p, v)
		}

		t := time.Now()
		ms := t.Sub(start)

		f, err := c.repo.PlayerInsert(ctx, p)
		if err != nil {
			l.Fatal("players insert", tel.Error(err))
			return
		}

		eTime := time.Now().Sub(t)

		l.Info("progress",
			tel.String("T", "Player"),
			tel.Float64("comb/sec", float64(out)/ms.Seconds()),
			tel.Duration("duration", ms),
			tel.Int("out", out),
			tel.Duration("last", eTime),
			tel.String("file", f),
		)

		atomic.StoreUint64(&c.PlayerMaxID, c.PlayerMaxID+uint64(num))
	}
}

func (c *controller) DoCB(ctx context.Context, rate int) {
	l := tel.FromCtx(ctx)

	start := time.Now()
	last := start
	var out int

	for {
		sleepTime := c.WindowTime - time.Now().Sub(last)
		if sleepTime < 0 {
			sleepTime = 0
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleepTime):
		}

		secs := time.Now().Sub(last)

		num := int(float64(rate) * secs.Minutes())
		if num < 1 {
			continue
		}

		last = time.Now()
		out += num

		l.Info(">>", tel.Int("num", num))

		var p []*model.CBet
		for i := 0; i < num; i++ {
			ID := rand.Intn(int(atomic.LoadUint64(&c.PlayerMaxID))) + 1
			v := model.NewCBet(ID)
			v.OK()
			p = append(p, v)
		}

		t := time.Now()
		ms := t.Sub(start)

		f, err := c.repo.CasinoBetInsert(ctx, p)
		if err != nil {
			l.Fatal("casino bet insert", tel.Error(err))
			return
		}

		eTime := time.Now().Sub(t)

		l.Info("progress",
			tel.String("T", "CBet"),
			tel.Float64("comb/sec", float64(out)/ms.Seconds()),
			tel.Duration("duration", ms),
			tel.Int("out", out),
			tel.Duration("last", eTime),
			tel.String("file", f),
		)
	}
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
