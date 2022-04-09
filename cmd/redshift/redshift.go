package redshift

import (
	"github.com/d7561985/redshift-test/internal/config"
	"github.com/d7561985/redshift-test/pkg/service"
	"github.com/d7561985/redshift-test/store/postgres"
	"github.com/d7561985/redshift-test/store/s3"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const defMaxUserID = 100_000
const defThreads = 100
const cSize = 500

var dbConnect = "redshift://localhost:5439/dev?sslmode=disable"

const (
	Transaction = "tx"
	Insert      = "insert"
)

const (
	fThreads = "threads"
	fMaxUser = "maxUser"
	fOpt     = "operation"

	fStore = "store"

	fAddr = "addr"
	fS3   = "s3"
)

const (
	EnvThreads   = "THREADS"
	EnvMaxUser   = "MAX_USER"
	EnvOperation = "OPERATION"
	EnvAddr      = "REDSHIFT_ADDR"
	EnvS3Bucket  = "S3_BUCKET"
	EnvStore     = "STORE"
)

const (
	StoreRedshift = "redshift"
	StoreS3       = "s3"
)

var errNoStore = errors.New("unsupported store type")

type postgresCommand struct{}

func New() *cli.Command {
	c := new(postgresCommand)

	return &cli.Command{
		Name:        "redshift",
		Description: "run postgres compliance test which runs transactions",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: fThreads, Value: defThreads, Aliases: []string{"t"}, EnvVars: []string{EnvThreads}},
			&cli.IntFlag{Name: fMaxUser, Value: defMaxUserID, Aliases: []string{"m"}, EnvVars: []string{EnvMaxUser}},
			&cli.StringFlag{Name: fOpt, Value: Insert, Usage: "What test start: tx - transaction intense, insert - only insert", Aliases: []string{"o"}, EnvVars: []string{EnvOperation}},

			&cli.StringFlag{Name: fStore, Value: StoreS3, EnvVars: []string{EnvStore}},
			&cli.StringFlag{Name: fS3, EnvVars: []string{EnvS3Bucket}, Required: true},
			&cli.StringFlag{Name: fAddr, Value: dbConnect, EnvVars: []string{EnvAddr}},
		},
		Action: c.Action,
	}
}

func (m *postgresCommand) Action(c *cli.Context) error {
	repo, err := repoFactory(c)
	if err != nil {
		return errors.WithStack(err)
	}

	svc := service.New(service.Config{
		Size:    cSize,
		MaxUser: c.Int(fMaxUser),
	}, repo)

	svc.Run(c.Context)

	return nil
}

func repoFactory(c *cli.Context) (service.Repo, error) {
	switch c.String(fStore) {
	case StoreRedshift:
		repo, err := postgres.New(c.Context, config.Postgres{
			Addr: c.String(fAddr),
		})

		return repo, errors.WithStack(err)
	case StoreS3:
		return s3.New(c.String(fS3)), nil
	default:
		return nil, errors.WithStack(errNoStore)
	}
}
