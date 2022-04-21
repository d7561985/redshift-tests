package redshift

import (
	"time"

	"github.com/d7561985/redshift-test/pkg/repo/test"
	"github.com/d7561985/redshift-test/pkg/s3copy"
	"github.com/d7561985/redshift-test/pkg/service"
	"github.com/d7561985/redshift-test/store/pgxx"
	"github.com/d7561985/redshift-test/store/s3"
	"github.com/d7561985/tel/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const defMaxUserID = 1
const defCbRate = 100 * 60
const defaultWindowTime = time.Second * 10

var dbConnect = "redshift://localhost:5439/dev?sslmode=disable"

const (
	fUserRPM = "userRatePerMinute"
	fCbRPM   = "cbRatePerMinute"
	fWindow  = "windowTime"

	fStore = "store"

	fAddr  = "addr"
	fS3    = "s3"
	fIRole = "imRole"
)

const (
	EnvMaxUser    = "USER_RATE_PER_MINUTE"
	EnvCbRate     = "CB_RATE_PER_MINUTE"
	EnvWindowTime = "WINDOW_TIME"

	EnvAddr     = "REDSHIFT_ADDR"
	EnvS3Bucket = "S3_BUCKET"
	EnvIMRole   = "RED_SHIFT_ATTACHED_IAM_ROLE" // for S3 copy
	EnvStore    = "STORE"
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
			&cli.IntFlag{Name: fUserRPM, Value: defMaxUserID, Aliases: []string{"ur"}, EnvVars: []string{EnvMaxUser}},
			&cli.IntFlag{Name: fCbRPM, Value: defCbRate, Aliases: []string{"cbr"}, EnvVars: []string{EnvCbRate}},
			&cli.DurationFlag{Name: fWindow, Value: defaultWindowTime, Aliases: []string{"w"}, EnvVars: []string{EnvWindowTime}},

			&cli.StringFlag{Name: fStore, Value: StoreS3, EnvVars: []string{EnvStore}},
			&cli.StringFlag{Name: fS3, EnvVars: []string{EnvS3Bucket}, Required: true},
			&cli.StringFlag{Name: fIRole, EnvVars: []string{EnvIMRole}, Required: true},
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

	r, err := pgxx.New(pgxx.Cfg{Addr: c.String(fAddr)})
	if err != nil {
		return errors.WithStack(err)
	}

	redshift := test.New(r)
	PlayerMaxID, err := redshift.PlayerGetLastID(c.Context)
	if err != nil {
		return errors.WithStack(err)
	}

	tel.FromCtx(c.Context).Info("max players", tel.Int("val", PlayerMaxID))

	cp := s3copy.New(s3copy.Cfg{
		Bucket: c.String(fS3),
		IRole:  c.String(fIRole),
	}, r)

	svc := service.New(service.Config{
		PlayerMaxID: uint64(PlayerMaxID),
		PlayerRate:  c.Int(fUserRPM),
		CBRate:      c.Int(fCbRPM),
		WindowTime:  c.Duration(fWindow),
	}, repo, cp)

	svc.Run(c.Context)

	return nil
}

func repoFactory(c *cli.Context) (service.Repo, error) {
	switch c.String(fStore) {
	case StoreRedshift:
		r, err := pgxx.New(pgxx.Cfg{
			Addr: c.String(fAddr),
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return test.New(r), errors.WithStack(err)
	case StoreS3:
		return s3.New(c.String(fS3)), nil
	default:
		return nil, errors.WithStack(errNoStore)
	}
}
