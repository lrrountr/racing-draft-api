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
