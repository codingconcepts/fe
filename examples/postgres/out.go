// GENERATED CODE! Don't modify.

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseFunctions contains your database functions.
type DatabaseFunctions struct {
	db *pgxpool.Pool
}

func (df *DatabaseFunctions) GetOldestPerson(ctx context.Context) (string, error) {
	const stmt = `
  SELECT full_name
  FROM person
  ORDER BY age DESC
  LIMIT 1;
`

	row := df.db.QueryRow(ctx, stmt)

	var val string
	if err := row.Scan(&val); err != nil {
		return "", fmt.Errorf("calling get_oldest_person: %w", err)
	}

	return val, nil
}

func (df *DatabaseFunctions) PeopleBornOn(ctx context.Context, d time.Time) (int64, error) {
	const stmt = `
  SELECT count(*)
  FROM person
  WHERE date_of_birth = d;
`

	row := df.db.QueryRow(ctx, stmt, d)

	var val int64
	if err := row.Scan(&val); err != nil {
		return 0, fmt.Errorf("calling people_born_on: %w", err)
	}

	return val, nil
}

func (df *DatabaseFunctions) AddPerson(ctx context.Context, fullName string, dateOfBirth time.Time, country string) error {
	const stmt = `
  INSERT INTO person (full_name, date_of_birth, country) VALUES
    (full_name, date_of_birth, country);
`

	_, err := df.db.Exec(ctx, stmt, fullName, dateOfBirth, country)
	if err != nil {
		return fmt.Errorf("calling add_person: %w", err)
	}

	return nil
}
