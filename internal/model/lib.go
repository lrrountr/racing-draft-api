package model

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"

	"github.com/lrrountr/racing-draft-api/internal/config"
	"github.com/lrrountr/racing-draft-api/internal/model/migrations"
)

type DBClient struct {
	pool *pgxpool.Pool
}

type PageInfo struct {
	Total int
	Next  int
}

func (i *PageInfo) Update(totalItems, itemsInPage, offset int) {
	i.Total = totalItems
	i.Next = offset + itemsInPage
	if i.Next > totalItems {
		i.Next = totalItems
	}
}

func (c DBClient) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	tx, err = c.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start TX: %w", err)
	}
	return tx, nil
}

func NewClient(conf config.Config) (client DBClient, err error) {
	ctx, cls := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cls()

	dbConfig := conf.DB
	tlsMode := "prefer"
	var rootCAs *x509.CertPool = nil
	if conf.DB.Region != "" {
		rootCAs, err = setupTLSPEMBundle(dbConfig)
		if err != nil {
			return client, err
		}
		tlsMode = "verify-full"
	}

	dbConf, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.DBName,
			tlsMode,
		),
	)
	if err != nil {
		return client, fmt.Errorf("could not parse config: %w", err)
	}
	if rootCAs != nil {
		dbConf.ConnConfig.TLSConfig.RootCAs = rootCAs
	}

	pool, err := pgxpool.ConnectConfig(ctx, dbConf)
	if err != nil {
		return client, fmt.Errorf("could not connect to database: %w", err)
	}
	client = DBClient{
		pool: pool,
	}
	p, err := client.pool.Acquire(ctx)
	if err != nil {
		return client, err
	}
	defer p.Release()

	conn := p.Conn()
	ms, err := migrate.NewMigrator(ctx, conn, migrations.MigrationsTable)
	if err != nil {
		return client, fmt.Errorf("could not create migrator: %w", err)
	}

	ms = migrations.Init(ms)
	err = ms.Migrate(ctx)
	if err != nil {
		return client, fmt.Errorf("could not migrate table: %w", err)
	}

	return client, nil
}

func setupTLSPEMBundle(conf config.DBConf) (*x509.CertPool, error) {
	url := fmt.Sprintf("https://truststore.pki.rds.amazonaws.com/%s/%s-bundle.pem", conf.Region, conf.Region)
	c := &http.Client{
		Timeout: 10 * time.Minute,
	}
	resp, err := c.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get TLS PEM bundle: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code downloading PEM bundle: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	chain, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rootCAs := x509.NewCertPool()
	ok := rootCAs.AppendCertsFromPEM(chain)
	if !ok {
		return nil, errors.New("AWS returned invalid Cert PEMs")
	}
	return rootCAs, nil
}

func (c DBClient) Stat() *pgxpool.Stat {
	return c.pool.Stat()
}
