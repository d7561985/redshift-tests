package pgxx

import (
	"sync"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stoewer/go-strcase"
)

const (
	driverName = "pgx"
)

const connectionLifetime = 15 * time.Minute

var once sync.Once //nolint:gochecknoglobals // used for registering the driver

type Cfg struct {
	Addr string

	MaxOpenConn int
	MaxIdleConn int
}

type Storage struct {
	*sqlx.DB

	cfg Cfg
}

func New(cfg Cfg) (*Storage, error) {
	sqlx.NameMapper = strcase.LowerCamelCase

	//once.Do(func() {
	//	sql.Register(driverName,
	//		instrumentedsql.WrapDriver(&pq.Driver{},
	//			instrumentedsql.WithTracer(opentracing.NewTracer(true)),
	//			instrumentedsql.WithOpsExcluded(
	//				instrumentedsql.OpSQLConnectorConnect,
	//				instrumentedsql.OpSQLPing,
	//				instrumentedsql.OpSQLPrepare,
	//				instrumentedsql.OpSQLRowsNext,
	//				instrumentedsql.OpSQLStmtClose,
	//			),
	//		))
	//})

	db, err := sqlx.Open(driverName, cfg.Addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open PostgreSQL")
	}

	if cfg.MaxOpenConn != 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConn)
	}

	if cfg.MaxIdleConn != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConn)
	}

	db.SetConnMaxLifetime(connectionLifetime)

	if err = db.Ping(); err != nil {
		return nil, errors.Wrapf(err, "failed ping")
	}

	return &Storage{DB: db, cfg: cfg}, nil
}

// MustRegisterMetrics registers prometheus metrics for given Database.
// Should be called only once for same Database, otherwise it panics on the second call.
func (s *Storage) MustRegisterMetrics() {
	collector := sqlstats.NewStatsCollector("xxx", s.DB)
	prometheus.MustRegister(collector)
}
