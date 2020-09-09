package lint

import (
	"fmt"
	"regexp"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
)

type checker struct {
	config   *Config
	warnings []Warning
	path     []ast.Node
}

func (c *checker) warn(n ast.Node, format string, args ...interface{}) {
	tok := n.GetToken()
	c.warnings = append(c.warnings, Warning{
		Line:   tok.Position.Line,
		Column: tok.Position.Column,
		Text:   fmt.Sprintf(format, args...),
	})
}

func (c *checker) getParent(n int) ast.Node {
	index := len(c.path) - n - 1
	if index >= 0 && index < len(c.path) {
		return c.path[index]
	}
	return nil
}

func (c *checker) isInsideObject() bool {
	_, ok := c.getParent(1).(*ast.MappingNode)
	return ok
}

func (c *checker) visitNode(n ast.Node) {
	c.path = append(c.path, n)
	defer func() { c.path = c.path[:len(c.path)-1] }()

	switch n := n.(type) {
	case *ast.NullNode, *ast.IntegerNode, *ast.FloatNode, *ast.BoolNode:
		// OK.

	case *ast.NanNode:
		c.warn(n, "NaN value should not be used")
	case *ast.InfinityNode:
		c.warn(n, "infinity value should not be used")
	case *ast.LiteralNode:
		c.warn(n, "literal block scalar '%s' should not be used", n.Start.Value)

	case *ast.AliasNode:
		c.warn(n, "remove %s alias", n.Value)
		c.visitNode(n.Value)
	case *ast.AnchorNode:
		c.warn(n, "remove %s anchor", n.Name)
		c.visitNode(n.Name)
		c.visitNode(n.Value)
	case *ast.TagNode:
		c.warn(n, "remove %s tag", n.Start.Value)
		c.visitNode(n.Value)

	case *ast.SequenceNode:
		if !n.IsFlowStyle {
			c.warn(n, "use a flow array syntax instead")
		}
		for _, v := range n.Values {
			c.visitNode(v)
		}

	case *ast.StringNode:
		c.visitString(n, false)

	case *ast.MappingValueNode:
		if !c.isInsideObject() {
			c.warn(n, "used a key-value outside of an object")
		}
		c.visitObjectKey(n.Key)
		c.visitNode(n.Value)

	case *ast.MappingNode:
		for _, v := range n.Values {
			c.visitNode(v)
		}

	default:
		panic(fmt.Sprintf("unhandled node %T", n))
	}
}

func (c *checker) visitFile(f *ast.File) {
	if len(f.Docs) != 1 {
		c.warn(f.Docs[1], "found more than one document inside a file")
	}

	for _, doc := range f.Docs {
		c.visitNode(doc.Body)
	}
}

func (c *checker) visitString(s *ast.StringNode, isKey bool) {
	switch c.getParent(1).(type) {
	case *ast.AnchorNode, *ast.AliasNode:
		return
	}

	if s.Token.Type == token.SingleQuoteType {
		if !c.config.Allow.SingleQuoteStrings {
			c.warn(s, "single quote strings are not allowed")
		}
		return
	}

	if s.Token.Type == token.DoubleQuoteType {
		return // It's always OK to have ""-literals
	}

	if !isKey {
		c.warn(s, "unquoted strings are not allowed")
		return
	}

	if !c.config.Allow.IdentObjKeys {
		c.warn(s, "unquoted object keys are not allowed")
		return
	}

	if !c.isValidIdent(s.Value) {
		c.warn(s, "`%s` is not a valid ES5.1 object key", s.Value)
		return
	}
}

func (c *checker) visitObjectKey(k ast.Node) {
	switch k := k.(type) {
	case *ast.StringNode:
		c.visitString(k, true)
	case *ast.BoolNode:
		switch k.Token.Value {
		case "true":
			c.warn(k, "`true` is not a valid ES5.1 object key")
		case "false":
			c.warn(k, "`false` is not a valid ES5.1 object key")
		default:
			c.visitNode(k)
		}
	case *ast.NullNode:
		if k.Token.Value == "null" {
			c.warn(k, "`null` is not a valid ES5.1 object key")
		} else {
			c.visitNode(k)
		}
	default:
		c.visitNode(k)
	}
}

func (c *checker) isValidIdent(s string) bool {
	// See https://www.ecma-international.org/ecma-262/5.1/#sec-7.6
	if _, isReserved := es5reserved[s]; isReserved {
		return false
	}
	// TODO: allow non-ascii letters in identifiers.
	// Right now we're using a simple regexp that covers only ASCII identifiers.
	return es5identRE.MatchString(s)
}

var es5identRE = regexp.MustCompile(`^[$_A-Za-z][\w$]*$`)

var es5reserved = newStringSet(
	"null",
	"true",
	"false",

	// FutureReservedWord
	"class",
	"enum",
	"extends",
	"super",
	"const",
	"export",
	"import",

	// Keyword
	"break",
	"case",
	"catch",
	"continue",
	"debugger",
	"default",
	"delete",
	"do",
	"else",
	"finally",
	"for",
	"function",
	"if",
	"in",
	"instanceof",
	"new",
	"return",
	"switch",
	"this",
	"throw",
	"try",
	"typeof",
	"var",
	"void",
	"while",
	"with",
)

func newStringSet(keys ...string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
