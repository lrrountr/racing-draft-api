package migrations

const (
	CreatePGExtensionsUp = `
		CREATE EXTENSION pg_trgm;
		CREATE EXTENSION fuzzystrmatch;
	`
	CreatePGExtensionsDown = `
		DROP EXTENSION pg_trgm;
		DROP EXTENSION fuzzystrmatch;
	`
)

const (
	SeasonsTableUp = `
		BEGIN;
			CREATE TABLE seasons(
				uuid TEXT NOT NULL,
				name TEXT NOT NULL,
				racing_series TEXT NOT NULL DEFAULT 'F1',
				year INT NOT NULL,
				active BOOL NOT NULL
			);
			CREATE UNIQUE INDEX unique_season_uuid ON seasons(uuid); 
		COMMIT;
	`
	SeasonsTableDown = `
		BEGIN;
			DROP INDEX unique_season_uuid;
			DROP TABLE seasons;
		COMMIT;
	`
)

const (
	UsersTableUp = `
		BEGIN;
			CREATE TABLE users(
				uuid TEXT NOT NULL,
				name TEXT NOT NULL,
				email TEXT NOT NULL
			);
			CREATE UNIQUE INDEX unique_user_email ON users(email);
			CREATE INDEX users_search_terms_idx ON users USING GIN(to_tsvector('english', name || ' ' || email));
		COMMIT;
	`
	UsersTableDown = `
		BEGIN;
			DROP INDEX unique_user_email;
			DROP INDEX users_search_terms_idx;
			DROP TABLE users;
		COMMIT;
	`
)
