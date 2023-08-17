// GENERATED CODE! Don't modify.

package repo

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

// NewDatabaseFunctions returns a pointer to a new instance of DatabaseFunctions.
func NewDatabaseFunctions(db *pgxpool.Pool) *DatabaseFunctions {
	return &DatabaseFunctions{
		db: db,
	}
}

func (df *DatabaseFunctions) PeopleBornOn(ctx context.Context, d time.Time) (int64, error) {
	const stmt = `SELECT count(*) FROM person WHERE date_of_birth = $1`

	row := df.db.QueryRow(ctx, stmt, d)

	var result int64
	if err := row.Scan(&result); err != nil {
		return 0, fmt.Errorf("calling people_born_on: %w", err)
	}

	return result, nil
}

func (df *DatabaseFunctions) PeopleBetween(ctx context.Context, idFrom string, idTo string) (interface{}, error) {
	const stmt = `SELECT id, country, full_name, date_of_birth FROM person WHERE id BETWEEN $1 AND $2`

	rows, err := df.db.Query(ctx, stmt, idFrom, idTo)
	if err != nil {
		return nil, fmt.Errorf("calling people_between: %w", err)
	}

	var results []interface{}
	for rows.Next() {
		var result interface{}
		if err := rows.Scan(&result); err != nil {
			return nil, fmt.Errorf("calling people_between: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

func (df *DatabaseFunctions) AddPerson(ctx context.Context, fullName string, dateOfBirth time.Time, country string) error {
	const stmt = `INSERT INTO person (full_name, date_of_birth, country) VALUES ($1, $2, $3)`

	_, err := df.db.Exec(ctx, stmt, fullName, dateOfBirth, country)
	if err != nil {
		return fmt.Errorf("calling add_person: %w", err)
	}

	return nil
}

func (df *DatabaseFunctions) GetOldestPerson(ctx context.Context) (string, error) {
	const stmt = `SELECT full_name FROM person ORDER BY date_of_birth DESC LIMIT 1`

	row := df.db.QueryRow(ctx, stmt)

	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("calling get_oldest_person: %w", err)
	}

	return result, nil
}
