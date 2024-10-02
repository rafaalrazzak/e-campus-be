package migrations

import (
	"context"
	"github.com/go-bun/bun-starter-kit/ecampus"

	"github.com/go-bun/bun-starter-kit/bunapp"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		db.RegisterModel((*ecampus.User)(nil))

		fixture := dbfixture.New(db, dbfixture.WithRecreateTables())
		return fixture.Load(ctx, bunapp.FS(), "fixture/fixture.yml")
	}, nil)
}
