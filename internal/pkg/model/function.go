package model

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Function describes the properties of a database function.
type Function struct {
	Name         string
	Language     string
	ReturnType   string
	ArgNames     []string
	ArgTypes     []string
	FunctionBody string
}

// Args returns something that can be used in code.
func (f Function) Args(lang string) string {
	switch strings.ToLower(lang) {
	case "go":
		// Return args in the format: "name type, name type"
		vals := lo.Map(lo.Zip2(f.ArgNames, f.ArgTypes), func(t lo.Tuple2[string, string], i int) string {
			return fmt.Sprintf("%s %s", toCamelCase(t.A), toGoType(t.B))
		})

		return strings.Join(vals, ", ")

	default:
		log.Fatalf("unimplemented language: %s", lang)
		return ""
	}
}

func (f Function) QueryArgs(lang string) string {
	switch strings.ToLower(lang) {
	case "go":
		// Return args in the format: "name, name"
		vals := lo.Map(f.ArgNames, func(arg string, i int) string {
			return toCamelCase(arg)
		})

		return strings.Join(vals, ", ")

	default:
		log.Fatalf("unimplemented language: %s", lang)
		return ""
	}
}

func (f Function) LanguageReturnType(lang string) string {
	switch strings.ToLower(lang) {
	case "go":
		return toGoType(f.ReturnType)

	default:
		log.Fatalf("unimplemented language: %s", lang)
		return ""
	}
}

func (f Function) DefaultReturnValue(lang string) string {
	switch strings.ToLower(lang) {
	case "go":
		return defaultValue(f.ReturnType)

	default:
		log.Fatalf("unimplemented language: %s", lang)
		return ""
	}
}

// HasArgs returns true if this function accepts arguments.
func (f Function) HasArgs() bool {
	return len(f.ArgNames) > 0
}

// ReturnsValue returns true if this function returns an argument.
func (f Function) ReturnsValue() bool {
	return strings.ToLower(f.ReturnType) != "void"
}

// ToPascalCase converts a database into PascalCase.
func (f Function) ToPascalCase() string {
	return toPascalCase(f.Name)
}

func toPascalCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und).String(s)
	s = strings.ReplaceAll(s, " ", "")

	return s
}

func (f Function) ToCamelCase() string {
	return toCamelCase(f.ToPascalCase())
}

func toCamelCase(s string) string {
	s = toPascalCase(s)
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func toGoType(dbType string) string {
	switch strings.ToLower(dbType) {
	case "uuid", "varchar":
		return "string"
	case "int", "int2", "int4", "int8":
		return "int64"
	case "date", "timestamp", "timestamptz":
		return "time.Time"
	default:
		log.Fatalf("unimplemented type: %s", dbType)
		return ""
	}
}

func defaultValue(dbType string) string {
	switch strings.ToLower(dbType) {
	case "uuid", "varchar":
		return `""`
	case "int", "int2", "int4", "int8", "float", "decimal":
		return `0`
	case "date", "timestamp", "timestamptz":
		return `time.Time{}`
	default:
		log.Fatalf("unimplemented type: %s", dbType)
		return ""
	}
}
