package lint

import (
	"github.com/goccy/go-yaml/ast"
)

type Config struct {
	Allow struct {
		SingleQuoteStrings bool
		IdentObjKeys       bool
	}
}

type Warning struct {
	Line   int
	Column int
	Text   string
}

func Run(config *Config, f *ast.File) []Warning {
	c := checker{config: config}
	c.visitFile(f)
	return c.warnings
}
