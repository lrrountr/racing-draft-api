//go:build test

package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/lrrountr/racing-draft-api/internal/config"
)

func EnsureDatabase(conf config.DBConf) (err error) {
	ctx, cls := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cls()

	dbConf, err := pgxpool.ParseConfig(fmt.Sprintf("postgresql://%s:%s@%s:%d", conf.User, conf.Password, conf.Host, conf.Port))
	if err != nil {
		return fmt.Errorf("could not parse config: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, dbConf)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer pool.Close()

	p, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer p.Release()

	row := p.QueryRow(ctx, "SELECT COUNT(*) FROM pg_database WHERE datname = $1", conf.DBName)
	count := 0
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Creae database cannot be parametrized - but this is just a test function
	_, err = p.Exec(ctx, fmt.Sprintf("CREATE DATABASE \"%s\"", conf.DBName))
	if err != nil {
		return err
	}
	return nil
}
