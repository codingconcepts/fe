package model

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	pg_query "github.com/pganalyze/pg_query_go/v4"
)

// Function describes the properties of a database function.
type Function struct {
	Name         string
	Language     string
	ReturnType   string
	ReturnsSet   bool
	ArgNames     []string
	ArgTypes     []string
	FunctionBody string
}

// ReturnsRecord returns true if this function returns one or move records.
func (f Function) ReturnsRecord() bool {
	return strings.EqualFold(f.ReturnType, "record")
}

func (f Function) SafeFunctionBody() string {
	tree, err := pg_query.Parse(f.FunctionBody)
	if err != nil {
		log.Fatalf("error parsing query: %v", err)
	}

	for _, treeStmt := range tree.Stmts {
		switch stmt := treeStmt.Stmt.Node.(type) {
		case *pg_query.Node_SelectStmt:
			f.subSelect(stmt.SelectStmt)
		case *pg_query.Node_InsertStmt:
			f.subInsert(stmt.InsertStmt)
		case *pg_query.Node_UpdateStmt:
			f.subUpdate(stmt.UpdateStmt)
		case *pg_query.Node_DeleteStmt:
			f.subDelete(stmt.DeleteStmt)
		}
	}

	newStmt, err := pg_query.Deparse(tree)
	if err != nil {
		log.Fatalf("error rebuilding query: %v", err)
	}

	return subArgNumbers(newStmt)
}

func (f Function) subSelect(stmt *pg_query.SelectStmt) {
	if stmt.WhereClause != nil {
		aexpr := stmt.WhereClause.GetAExpr().Rexpr

		switch x := aexpr.Node.(type) {
		case *pg_query.Node_ColumnRef:
			for _, field := range x.ColumnRef.Fields {
				*field.GetString_() = pg_query.String{Sval: "999999999"}
			}

		case *pg_query.Node_List:
			for _, item := range x.List.Items {
				for _, field := range item.GetColumnRef().Fields {
					*field.GetString_() = pg_query.String{Sval: "999999999"}
				}
			}

		default:
			log.Fatalf("unsupported type: %T", aexpr.Node)
		}
	}
}

func (f Function) subInsert(stmt *pg_query.InsertStmt) {
	for _, vl := range stmt.SelectStmt.GetSelectStmt().ValuesLists {
		for _, i := range vl.GetList().Items {
			for _, f := range i.GetColumnRef().Fields {
				*f.GetString_() = pg_query.String{Sval: "999999999"}
			}
		}
	}
}

func (f Function) subUpdate(stmt *pg_query.UpdateStmt) {
	targets := stmt.GetTargetList()
	f.subNodeValues(targets)

	args := stmt.WhereClause.GetBoolExpr().Args
	for i := range args {
		rexpr := args[i].GetAExpr().GetRexpr()

		switch x := rexpr.Node.(type) {
		case *pg_query.Node_AConst:
			x.AConst.GetSval().Sval = "999999999"

		case *pg_query.Node_ColumnRef:
			for _, field := range x.ColumnRef.Fields {
				*field.GetString_() = pg_query.String{Sval: "999999999"}
			}

		case *pg_query.Node_List:
			for _, item := range x.List.Items {
				for _, field := range item.GetColumnRef().Fields {
					*field.GetString_() = pg_query.String{Sval: "999999999"}
				}
			}

		default:
			log.Fatalf("unimplemented node type: %T", rexpr)
		}
	}
}

func (f Function) subDelete(stmt *pg_query.DeleteStmt) {
	log.Fatal("delete statements not yet supported")
}

func (f Function) subNodeValues(nodes []*pg_query.Node) {
	for i := range nodes {
		nodes[i].GetResTarget().Val = pg_query.MakeAConstIntNode(999999999, -1)
	}
}

func subArgNumbers(stmt string) string {
	curr := 1
	next := func(string) string {
		c := curr
		curr++

		return fmt.Sprintf("$%d", c)
	}

	re := regexp.MustCompile(`"?999999999"?`)
	return re.ReplaceAllStringFunc(stmt, next)
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
		return f.defaultValue(f.ReturnType)

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
	case "record":
		return "map[string]any"
	default:
		log.Fatalf("unimplemented type: %s", dbType)
		return ""
	}
}

func (f Function) defaultValue(dbType string) string {
	if f.ReturnsSet {
		return `nil`
	}

	switch strings.ToLower(dbType) {
	case "uuid", "varchar":
		return `""`
	case "int", "int2", "int4", "int8", "float", "decimal":
		return `0`
	case "date", "timestamp", "timestamptz":
		return `time.Time{}`
	case "record":
		return `nil`
	default:
		log.Fatalf("unimplemented type: %s", dbType)
		return ""
	}
}
