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
