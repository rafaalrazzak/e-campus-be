package dbdef

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

func Migrator(db *sql.DB) error {
	goose.SetBaseFS(EmbedMigrations)
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}
	goose.SetTableName("goose_migrations")

	return goose.Up(db, "migrations")
}
