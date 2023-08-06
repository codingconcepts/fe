package repo

import (
	"database/sql"
	"fmt"

	"github.com/codingconcepts/fe/internal/pkg/model"
)

// PostgresRepo provides the necessary functionality for working against a
// Postgres database.
type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo returns a pointer to a new instance of PostgresRepo.
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

// GetFunctions returns all functions in the Postgres database.
func (r *PostgresRepo) GetFunctions() ([]model.Function, error) {
	const stmt = `
	SELECT
		pp.proname AS function_name,
		pl.lanname AS function_language,
		pt.typname AS function_return_type,
		pp.proargnames AS function_argument_names,
		ARRAY (
			SELECT pt.typname
			FROM ROWS FROM (unnest(pp.proargtypes))
			WITH ORDINALITY AS a (arg_id, ord)
			JOIN pg_type AS pt ON pt.oid = a.arg_id
			ORDER BY a.ord
		) AS function_argument_types,
		pp.prosrc AS function_sql
	FROM
		pg_proc AS pp
		INNER JOIN pg_namespace AS pn ON pp.pronamespace = pn.oid
		INNER JOIN pg_language AS pl ON pp.prolang = pl.oid
		INNER JOIN pg_type AS pt ON pp.prorettype = pt.oid
	WHERE
		pl.lanname NOT IN ('c', 'internal')
		AND pn.nspname NOT LIKE 'pg_%'
		AND pn.nspname != 'information_schema'`

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("querying functions: %w", err)
	}

	var functions []model.Function
	for rows.Next() {
		var f model.Function
		if err = rows.Scan(&f.Name, &f.Language, &f.ReturnType, &f.ArgNames, &f.ArtTypes, &f.FunctionBody); err != nil {
			return nil, fmt.Errorf("scanning function: %w", err)
		}
		functions = append(functions, f)
	}

	return functions, nil
}
