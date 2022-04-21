package s3copy

import (
	"context"
	"fmt"

	"github.com/d7561985/redshift-test/model"
	"github.com/d7561985/redshift-test/store/pgxx"
	"github.com/d7561985/tel/v2"
	"github.com/pkg/errors"
)

type Cfg struct {
	Bucket string
	IRole  string
}

type Repo struct {
	cfg Cfg
	*pgxx.Storage
}

func New(cfg Cfg, r *pgxx.Storage) *Repo {
	return &Repo{cfg: cfg, Storage: r}
}

func (r *Repo) Copy(ctx context.Context, dst model.Copy) error {
	sql := `COPY %s(%s)
FROM 's3://%s/%s'
iam_role '%s'
REGION 'eu-central-1'
CSV
IGNOREHEADER as 1
TIMEFORMAT AS 'epochmillisecs' --'DD.MM.YYYY HH:MI:SS'
GZIP;`

	sql = fmt.Sprintf(sql, dst.Table, dst.Fields, r.cfg.Bucket, dst.Path, r.cfg.IRole)
	tel.FromCtx(ctx).Debug("copy", tel.String("sql", sql))

	_, err := r.ExecContext(ctx, sql)

	return errors.WithStack(err)
}
