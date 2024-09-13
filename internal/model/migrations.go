package model

import (
	"context"

	"github.com/jackc/tern/migrate"
	"github.com/lrrountr/racing-draft-api/internal/model/migrations"
)

func (c *DBClient) Reinit() error {
	ctx := context.Background()
	p, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer p.Release()

	conn := p.Conn()
	ms, err := migrate.NewMigrator(ctx, conn, migrations.MigrationsTable)
	if err != nil {
		return err
	}
	ms = migrations.Init(ms)
	err = ms.MigrateTo(ctx, 0)
	if err != nil {
		return err
	}
	return ms.Migrate(ctx)
}
