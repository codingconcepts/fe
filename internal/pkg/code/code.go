package code

import (
	"io"

	"github.com/codingconcepts/fe/internal/pkg/model"
)

// CodeGenerator describes the behaviour of a code generator.
type CodeGenerator interface {
	Generate([]model.Function, io.Writer) error
}
