package model

// Function describes the properties of a database function.
type Function struct {
	Name         string
	Language     string
	ReturnType   string
	ArgNames     []string
	ArtTypes     []string
	FunctionBody string
}
