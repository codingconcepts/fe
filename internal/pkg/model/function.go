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
		subNode(stmt.WhereClause)
	}
}

func (f Function) subInsert(stmt *pg_query.InsertStmt) {
	for _, vl := range stmt.SelectStmt.GetSelectStmt().ValuesLists {
		for _, i := range vl.GetList().Items {
			subNode(i)
		}
	}
}

func (f Function) subUpdate(stmt *pg_query.UpdateStmt) {
	targets := stmt.GetTargetList()
	subNodes(targets)

	if stmt.WhereClause != nil {
		subNode(stmt.WhereClause)
	}
}

func (f Function) subDelete(stmt *pg_query.DeleteStmt) {
	if stmt.WhereClause != nil {
		subNode(stmt.WhereClause)
	}
}

func subNode(n *pg_query.Node) {
	switch x := n.Node.(type) {
	case *pg_query.Node_BoolExpr:
		for _, arg := range x.BoolExpr.Args {
			subNode(arg.GetAExpr().Rexpr)
		}
	case *pg_query.Node_ColumnRef:
		subNodes(x.ColumnRef.Fields)
	case *pg_query.Node_ResTarget:
		subNode(x.ResTarget.Val)
	case *pg_query.Node_AConst:
		switch x.AConst.Val.(type) {
		case *pg_query.A_Const_Sval:
			x.AConst.Val = &pg_query.A_Const_Sval{Sval: &pg_query.String{Sval: "999999999"}}
		case *pg_query.A_Const_Ival:
			x.AConst.Val = &pg_query.A_Const_Ival{Ival: &pg_query.Integer{Ival: 999999999}}
		case *pg_query.A_Const_Fval:
			x.AConst.Val = &pg_query.A_Const_Fval{Fval: &pg_query.Float{Fval: "999999999"}}
		}
	case *pg_query.Node_AExpr:
		subNode(x.AExpr.Rexpr)
	case *pg_query.Node_List:
		subNodes(x.List.Items)
	case *pg_query.Node_String_:
		*n.GetString_() = pg_query.String{Sval: "999999999"}
	case *pg_query.Node_Integer:
		*n.GetInteger() = pg_query.Integer{Ival: 999999999}
	}
}

func subNodes(nodes []*pg_query.Node) {
	for _, n := range nodes {
		subNode(n)
	}
}

func subArgNumbers(stmt string) string {
	curr := 1
	next := func(string) string {
		c := curr
		curr++

		return fmt.Sprintf("$%d", c)
	}

	re := regexp.MustCompile(`["']?999999999["']?`)
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
		return f.defaultValueGo(f.ReturnType)

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
	return toCamelCase(f.Name)
}

func toCamelCase(s string) string {
	s = toPascalCase(s)
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func toGoType(dbType string) string {
	switch strings.ToLower(dbType) {
	case "uuid", "varchar", "json":
		return "string"
	case "int", "int2", "int4", "int8":
		return "int64"
	case "numeric", "float4", "float8", "money":
		return "float64"
	case "bool":
		return "bool"
	case "date", "timestamp", "timestamptz":
		return "time.Time"
	case "record":
		return "map[string]any"
	default:
		log.Fatalf("unimplemented type: %s", dbType)
		return ""
	}
}

func (f Function) defaultValueGo(dbType string) string {
	if f.ReturnsSet {
		return `nil`
	}

	switch strings.ToLower(dbType) {
	case "uuid", "varchar", "json":
		return `""`
	case "int", "int2", "int4", "int8", "float4", "float8", "numeric", "money":
		return `0`
	case "bool":
		return `false`
	case "date", "timestamp", "timestamptz":
		return `time.Time{}`
	case "record":
		return `nil`
	default:
		log.Fatalf("unimplemented type: %s", dbType)
		return ""
	}
}
