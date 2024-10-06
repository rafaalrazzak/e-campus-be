package database

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database/dbdef"
	"go.uber.org/fx"
)

type ECampusDB struct {
	conn *sqlx.DB
	gqdb *goqu.Database
}

func NewECampusDBImpl(conn *sqlx.DB) *ECampusDB {
	gqdb := goqu.New("postgres", conn)
	return &ECampusDB{conn, gqdb}
}

func NewDatabaseConn(lc fx.Lifecycle, cfg config.Config) (conn *sqlx.DB, err error) {
	if conn, err = sqlx.Open("pgx", cfg.Database.Url); err != nil {
		return
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return conn.Ping()
		},
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return
}

func Migrator(lc fx.Lifecycle, db *sqlx.DB) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return dbdef.Migrator(db.DB)
		},
	})
}
