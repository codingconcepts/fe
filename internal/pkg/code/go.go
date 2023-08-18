package code

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strings"
	"text/template"

	_ "embed"

	"github.com/codingconcepts/fe/internal/pkg/model"
	"github.com/samber/lo"
)

//go:embed templates/postgres_go.go.tmpl
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

	var buf bytes.Buffer
	if err := cg.t.Execute(&buf, data); err != nil {
		return fmt.Errorf("executing code template: %w", err)
	}

	code, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting code output: %w", err)
	}

	if _, err = w.Write(code); err != nil {
		return fmt.Errorf("writing code to output file: %w", err)
	}

	return nil
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
