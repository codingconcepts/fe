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
			name:       "single arg json",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"json"},
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
			name:       "single arg numeric",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"numeric"},
			expArgsOut: "a float64",
		},
		{
			name:       "single arg float4",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"float4"},
			expArgsOut: "a float64",
		},
		{
			name:       "single arg float8",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"float8"},
			expArgsOut: "a float64",
		},
		{
			name:       "single arg money",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"money"},
			expArgsOut: "a float64",
		},
		{
			name:       "single arg bool",
			lang:       "go",
			argNames:   []string{"a"},
			argTypes:   []string{"bool"},
			expArgsOut: "a bool",
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
		name        string
		lang        string
		set         bool
		types       []string
		expDefaults []string
	}{
		{
			name: "go default scalar values",
			lang: "go",
			set:  false,
			types: []string{
				"varchar",
				"uuid",
				"json",
				"int",
				"int2",
				"int4",
				"int8",
				"numeric",
				"float4",
				"float8",
				"money",
				"bool",
				"date",
				"timestamp",
				"timestamptz",
				"record",
			},
			expDefaults: []string{
				`""`,
				`""`,
				`""`,
				`0`,
				`0`,
				`0`,
				`0`,
				`0`,
				`0`,
				`0`,
				`0`,
				`false`,
				`time.Time{}`,
				`time.Time{}`,
				`time.Time{}`,
				`nil`,
			},
		},
		{
			name: "go default set values",
			lang: "go",
			set:  true,
			types: []string{
				"varchar",
				"uuid",
				"json",
				"int",
				"int2",
				"int4",
				"int8",
				"numeric",
				"float4",
				"float8",
				"money",
				"bool",
				"date",
				"timestamp",
				"timestamptz",
				"record",
			},
			expDefaults: []string{
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
				`nil`,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if len(c.types) != len(c.expDefaults) {
				t.Fatalf("mismatched type and assertion lengths")
			}

			for i := 0; i < len(c.types); i++ {
				f := Function{
					ReturnsSet: c.set,
					ReturnType: c.types[i],
				}

				act := f.DefaultReturnValue(c.lang)
				assert.Equal(t, c.expDefaults[i], act)
			}
		})
	}
}
