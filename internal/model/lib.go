package model

type DBClient struct {
	pool *pgxpool.Pool
}
