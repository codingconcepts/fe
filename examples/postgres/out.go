// GENERATED CODE! Don't modify.

package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func (df *DatabaseFunctions) GetOldestPerson(ctx context.Context) (string, error) {
	const stmt = `SELECT full_name FROM person ORDER BY date_of_birth DESC LIMIT 1`

	row := df.db.QueryRow(ctx, stmt)

	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("calling get_oldest_person: %w", err)
	}

	return result, nil
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

	results, err := scan(rows)
	if err != nil {
		return nil, fmt.Errorf("calling people_between: %w", err)
	}

	return results, nil
}

func (df *DatabaseFunctions) PersonById(ctx context.Context, id string) (interface{}, error) {
	const stmt = `SELECT id, country, full_name, date_of_birth FROM person WHERE id = $1`

	row := df.db.QueryRow(ctx, stmt, id)

	var result interface{}
	if err := row.Scan(&result); err != nil {
		return nil, fmt.Errorf("calling person_by_id: %w", err)
	}

	return result, nil
}

func (df *DatabaseFunctions) AddPerson(ctx context.Context, fullName string, dateOfBirth time.Time, country string) error {
	const stmt = `INSERT INTO person (full_name, date_of_birth, country) VALUES ($1, $2, $3)`

	_, err := df.db.Exec(ctx, stmt, fullName, dateOfBirth, country)
	if err != nil {
		return fmt.Errorf("calling add_person: %w", err)
	}

	return nil
}

func scan(rows pgx.Rows) ([]map[string]any, error) {
	fields := rows.FieldDescriptions()

	var values []map[string]any
	for rows.Next() {
		scans := make([]any, len(fields))
		row := make(map[string]any)

		for i := range scans {
			scans[i] = &scans[i]
		}

		if err := rows.Scan(scans...); err != nil {
			return nil, fmt.Errorf("scaning values: %w", err)
		}

		for i, v := range scans {
			if v != nil {
				// Convert UUID into a string.
				if fields[i].DataTypeOID == 2950 {
					b := v.([16]byte)
					v = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
				}
				row[fields[i].Name] = v
			}
		}
		values = append(values, row)
	}

	return values, nil
}
