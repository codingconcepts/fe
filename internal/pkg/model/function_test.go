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
			inputStatement:     "SELECT a, b, c FROM t WHERE d = 1 AND e BETWEEN 2 AND 3",
			expOutputStatement: "SELECT a, b, c FROM t WHERE d = $1 AND e BETWEEN $2 AND $3",
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
