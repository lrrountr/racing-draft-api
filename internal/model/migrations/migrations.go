package migrations

import (
	"github.com/jackc/tern/migrate"
)

const (
	MigrationsTable = "public.migrations"
)

func Init(ms *migrate.Migrator) (nsm *migrate.Migrator) {
	ms.AppendMigration("001_Create PG Extensions", CreatePGExtensionsUp, CreatePGExtensionsDown)
	ms.AppendMigration("002_Create Seasons Table", SeasonsTableUp, SeasonsTableDown)
	ms.AppendMigration("003_Create Users Table", UsersTableUp, UsersTableDown)

	return ms
}
