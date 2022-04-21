package test

import (
	"context"

	"github.com/d7561985/redshift-test/model"
	"github.com/d7561985/redshift-test/store/pgxx"
	"github.com/d7561985/redshift-test/store/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Repo struct {
	*pgxx.Storage
}

const batchSize = 500

func (r *Repo) Bulk(ctx context.Context, p []*postgres.Journal) (string, error) {
	sql := `INSERT INTO journal("id","transactionId", "accountId", "balance", "change","currency","created_at",
"pincoinBalance","pincoinChange","project","revert","type"
) VALUES (:id,:transactionId,:accountId,:balance,:change,:currency,:createdAt,
					:pincoinBalance,:pincoinChange,:project,:revert,:type)`

	return "", errors.WithStack(
		insert(ctx, r.DB, sql, p),
	)
}

func (r *Repo) PlayerInsert(ctx context.Context, p []*model.Player) (string, error) {
	sql := `INSERT INTO players(id, guid, license, playerid, clickid, registerdate, language, email, 
                    isemailverify, phone, isphoneverify, ismultiaccount, birthday, accountverifytime,
                    lastlogintime, country, city, currency, sex, istest, isbot, project, activatestatus,
                    depositstatus, smsstatus, domain, webview, ipaddress, useragent, createunixnano, updateunixnano)
                     VALUES (
                          :id, :guid, :license, :playerId, :clickId, :registerDate, :language, :email, 
                    :isEmailVerify, :phone, :isPhoneVerify, :isMultiAccount, :birthday, :accountVerifyTime,
                    :lastLoginTime, :country, :city, :currency, :sex, :isTest, :isBot, :project, :activateStatus,
                    :depositStatus, :smsStatus, :domain, :webview, :ipAddress, :userAgent, :createUnixNano, :updateUnixNano 
                     )`

	return "", errors.WithStack(
		insert(ctx, r.DB, sql, p),
	)
}

func (r *Repo) CasinoBetInsert(ctx context.Context, p []*model.CBet) (string, error) {
	sql := `INSERT INTO cb(id, license, playerid, gamename, gametype, gameid, 
               bonusid, bet, winlose, purse, currencycode, gameprovider, gameroundid, 
               tranid, date, createunixnano, updateunixnano, rollback, status, error, 
               hall, "system", betinfo, agent, domain, webview, istournament) 
			VALUES (:id, :license, :playerId, :gameName, :gameType, :gameId, 
               :bonusId, :bet, :winLose, :purse, :currencyCode, :gameProvider, :gameRoundId, 
               :tranId, :date, :createUnixNano, :updateUnixNano, :rollback, :status, :error, 
               :hall, :system, :betInfo, :agent, :domain, :webview, :isTournament)`

	return "", errors.WithStack(
		insert(ctx, r.DB, sql, p),
	)
}

func insert[Q any](ctx context.Context, db *sqlx.DB, sql string, args []Q) error {
	wg := errgroup.Group{}

	for i := 0; i < len(args); i += batchSize {
		end := i + batchSize
		if end > len(args) {
			end = len(args)
		}

		i := i
		wg.Go(func() error {
			_, err := db.NamedExecContext(ctx, sql, args[i:end])
			if err != nil {
				return errors.WithStack(err)
			}
			return nil
		})
		_, err := db.NamedExecContext(ctx, sql, args[i:end])
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func New(r *pgxx.Storage) *Repo {
	return &Repo{Storage: r}
}

//===== SERVICE ======//

func (r *Repo) PlayerGetLastID(ctx context.Context) (int, error) {
	var res = new(int)
	err := r.Storage.QueryRow("SELECT MAX(id) FROM players").Scan(&res)

	if res == nil {
		return 0, errors.WithStack(err)
	}

	return *res, errors.WithStack(err)
}
