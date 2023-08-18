package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPascalCase(t *testing.T) {
	cases := []struct {
		name         string
		lang         string
		functionName string
		expName      string
	}{
		{
			name:         "snake case to go",
			lang:         "go",
			functionName: "transactions_by_day",
			expName:      "TransactionsByDay",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := Function{
				Name: c.functionName,
			}

			act := f.ToPascalCase()
			assert.Equal(t, c.expName, act)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	cases := []struct {
		name         string
		lang         string
		functionName string
		expName      string
	}{
		{
			name:         "snake case to go",
			lang:         "go",
			functionName: "transactions_by_day",
			expName:      "transactionsByDay",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := Function{
				Name: c.functionName,
			}

			act := f.ToCamelCase()
			assert.Equal(t, c.expName, act)
		})
	}
}

func TestSafeFunctionBody(t *testing.T) {
	cases := []struct {
		name               string
		inputStatement     string
		expOutputStatement string
	}{
		{
			name:               "select",
			inputStatement:     `SELECT a, b, c FROM t WHERE d = '1' AND e BETWEEN 2 AND 3`,
			expOutputStatement: `SELECT a, b, c FROM t WHERE d = $1 AND e BETWEEN $2 AND $3`,
		},
		{
			name:               "insert single-row dml",
			inputStatement:     "INSERT INTO t (a, b, c) VALUES (1, 2, 3)",
			expOutputStatement: "INSERT INTO t (a, b, c) VALUES ($1, $2, $3)",
		},
		{
			name:               "insert multi-row dml",
			inputStatement:     "INSERT INTO t (a, b, c) VALUES (1, 2, 3), (4, 5, 6)",
			expOutputStatement: "INSERT INTO t (a, b, c) VALUES ($1, $2, $3), ($4, $5, $6)",
		},
		{
			name:               "update",
			inputStatement:     "UPDATE t SET a = 1, b = 2, c = '3' WHERE d = 4",
			expOutputStatement: "UPDATE t SET a = $1, b = $2, c = $3 WHERE d = $4",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := Function{
				FunctionBody: c.inputStatement,
			}

			act := f.SafeFunctionBody()
			assert.Equal(t, c.expOutputStatement, act)
		})
	}
}

func TestArgs(t *testing.T) {
	cases := []struct {
		name       string
		lang       string
		argNames   []string
		argTypes   []string
		expArgsOut string
	}{
		{
			name:       "single arg varchar",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"varchar"},
			expArgsOut: "a string",
		},
		{
			name:       "single arg uuid",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"uuid"},
			expArgsOut: "a string",
		},
		{
			name:       "single arg int",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"int"},
			expArgsOut: "a int64",
		},
		{
			name:       "single arg int2",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"int2"},
			expArgsOut: "a int64",
		},
		{
			name:       "single arg int4",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"int4"},
			expArgsOut: "a int64",
		},
		{
			name:       "single arg int8",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"int8"},
			expArgsOut: "a int64",
		},
		{
			name:       "single arg date",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"date"},
			expArgsOut: "a time.Time",
		},
		{
			name:       "single arg timestamp",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"timestamp"},
			expArgsOut: "a time.Time",
		},
		{
			name:       "single arg timestamptz",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"timestamptz"},
			expArgsOut: "a time.Time",
		},
		{
			name:       "single arg record",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"record"},
			expArgsOut: "a map[string]any",
		},
		{
			name:       "multiple args",
			lang:       "go",
			argNames:   []string{"a", "b", "c"},
			argTypes:   []string{"uuid", "int", "date"},
			expArgsOut: "a string, b int64, c time.Time",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := Function{
				ArgNames: c.argNames,
				ArgTypes: c.argTypes,
			}

			act := f.Args(c.lang)
			assert.Equal(t, c.expArgsOut, act)
		})
	}
}

func TestDefaultValue(t *testing.T) {
	cases := []struct {
		name       string
		lang       string
		argNames   []string
		argTypes   []string
		expArgsOut string
	}{
		{
			name:       "single arg varchar",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"varchar"},
			expArgsOut: "a string",
		},
		{
			name:       "single arg uuid",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"uuid"},
			expArgsOut: "a string",
		},
		{
			name:       "single arg int",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"int"},
			expArgsOut: "a int64",
		},
		{
			name:       "single arg int2",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"int2"},
			expArgsOut: "a int64",
		},
		{
			name:       "single arg int4",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"int4"},
			expArgsOut: "a int64",
		},
		{
			name:       "single arg int8",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"int8"},
			expArgsOut: "a int64",
		},
		{
			name:       "single arg date",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"date"},
			expArgsOut: "a time.Time",
		},
		{
			name:       "single arg timestamp",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"timestamp"},
			expArgsOut: "a time.Time",
		},
		{
			name:       "single arg timestamptz",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"timestamptz"},
			expArgsOut: "a time.Time",
		},
		{
			name:       "single arg record",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"record"},
			expArgsOut: "a map[string]any",
		},
		{
			name:       "multiple args",
			lang:       "go",
			argNames:   []string{"a", "b", "c"},
			argTypes:   []string{"uuid", "int", "date"},
			expArgsOut: "a string, b int64, c time.Time",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := Function{
				ArgNames: c.argNames,
				ArgTypes: c.argTypes,
			}

			act := f.Args(c.lang)
			assert.Equal(t, c.expArgsOut, act)
		})
	}
}
