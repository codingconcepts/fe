package code

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	_ "embed"

	"github.com/codingconcepts/fe/internal/pkg/model"
	"github.com/samber/lo"
)

//go:embed templates/postgres_go.tmpl
var goTemplate string

// GoCodeGenerator is a Go-specific implementation of CodeGenerator.
type GoCodeGenerator struct {
	t *template.Template
}

// NewGoCodeGenerator returns a pointer to a new instance of GoCodeGenerator.
func NewGoCodeGenerator() (*GoCodeGenerator, error) {
	cg := GoCodeGenerator{}
	t, err := template.New("go").Parse(goTemplate)

	if err != nil {
		return nil, fmt.Errorf("parsing go code template: %w", err)
	}

	cg.t = t
	return &cg, nil
}

func (cg *GoCodeGenerator) Generate(functions []model.Function, w io.Writer, pkg string) error {
	data := struct {
		Package           string
		Functions         []model.Function
		AdditionalImports string
	}{
		Package:           pkg,
		Functions:         functions,
		AdditionalImports: additionalImports(functions),
	}

	return cg.t.Execute(w, data)
}

// additionalImports adds supporting imports for types like time.Time etc.
func additionalImports(functions []model.Function) string {
	var imports []string

	for _, f := range functions {
		args := f.ArgTypes
		args = append(args, f.ReturnType)

		for _, t := range args {
			switch strings.ToLower(t) {
			case "date", "timestamp", "timestamptz":
				imports = append(imports, `"time"`)
			}
		}
	}

	return strings.Join(lo.Uniq(imports), "\n")
}
