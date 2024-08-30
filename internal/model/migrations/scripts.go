package migrations

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
		COMMIT;
	`
	SeasonsTableDown = `
		BEGIN;
			DROP TABLE seasons;
		COMMIT;
	`
)
