package migrations

import (
	"github.com/jackc/tern/migrate"
)

const (
	MigrationsTable = "public.migrations"
)

func Init(ms *migrate.Migator) (nsm *migrate.Migator) {
	ms.AppendMigration("001_Create Seasons Table", SeasonsTableUp, SeasonsTableDown)

	return ms
}
