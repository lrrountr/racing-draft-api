package migrations

import (
	"github.com/jackc/tern/migrate"
)

const (
	MigrationsTable = "public.migrations"
)

func Init(ms *migrate.Migrator) (nsm *migrate.Migrator) {
	ms.AppendMigration("001_Create Seasons Table", SeasonsTableUp, SeasonsTableDown)

	return ms
}
