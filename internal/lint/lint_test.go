package lint

import (
	"fmt"
	"testing"

	"github.com/goccy/go-yaml/parser"
	"github.com/google/go-cmp/cmp"
)

func TestMappingKey(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{
  ? [1, 2]: "y",
  ? true : 'a',
}`
	test.want = []string{
		"2:3: don't use ?-style complex mapping key",
		"3:3: don't use ?-style complex mapping key",
	}
	runLintTest(test)
}

func TestDirective(t *testing.T) {
	test := newSuite(t)
	test.yaml = `%YAML 1.1
---
{}`
	test.want = []string{
		"3:1: found more than one document inside a file",
		"1:1: remove YAML 1.1 directive",
	}
	runLintTest(test)
}

func TestAnchorAndAlias(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{
  "anchor1": &"anchor1" [],
  "alias1": *"anchor1",

  "anchor2": &anchor2 [],
  "alias2": *anchor2,
}`
	test.want = []string{
		`2:13: remove "anchor1" anchor`,
		`3:12: remove "anchor1" alias`,
		"5:13: remove anchor2 anchor",
		"6:12: remove anchor2 alias",
	}
	runLintTest(test)
}

func TestBadObjKey(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{
  true: 2,
  a b: 1,
  false: 3,
    null: 4,
}`
	test.want = []string{
		"2:3: `true` is not a valid ES5.1 object key",
		"3:3: `a b` is not a valid ES5.1 object key",
		"4:3: `false` is not a valid ES5.1 object key",
		"5:5: `null` is not a valid ES5.1 object key",
	}
	runLintTest(test)
}

func TestNestedSinglelineObject(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{"a": {"b": [1]}}`
	test.want = []string{}
	runLintTest(test)
}

func TestMultiDoc(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{}
---
[]`
	test.want = []string{"3:1: found more than one document inside a file"}
	runLintTest(test)
}

func TestKeyValOutsideObj1(t *testing.T) {
	test := newSuite(t)
	test.yaml = `foo: [
  bar: 1,
  baz: 2,
]
`
	test.want = []string{
		"1:4: used a key-value outside of an object",
		"2:6: used a key-value outside of an object",
		"3:6: used a key-value outside of an object",
	}
	runLintTest(test)
}

func TestKeyValOutsideObj2(t *testing.T) {
	test := newSuite(t)
	test.yaml = `bar: 1
baz: 2
`
	test.want = []string{
		"1:4: use a flow object syntax {} instead",
	}
	runLintTest(test)
}

func TestNonFlowArray(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{
  "key":
    - 1
    - 2
}`
	test.want = []string{"3:5: use a flow array syntax [] instead"}
	runLintTest(test)
}

func TestUnquotedString(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{
  ok: "unquoted keys are allowed",
  also_ok: 'single quoted strings are allowed by the default',
  bad: unquoted strings values are bad,
  also_bad: [
    unquoted
  ]
}`
	test.want = []string{
		"4:8: unquoted strings are not allowed",
		"6:5: unquoted strings are not allowed",
	}
	runLintTest(test)
}

func TestTag(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{
  val: !!float 1.5,
}`
	test.want = []string{"2:8: remove !!float tag"}
	runLintTest(test)
}

func TestSpecialVals(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{
  bad1: .NAN,
  bad2: .inf,
  bad3: .Inf,
  bad4: yes,

  good1: null,  # OK
  good2: true,  # OK
  good3: false, # OK
}`
	test.want = []string{
		"2:9: NaN value should not be used",
		"3:9: infinity value should not be used",
		"4:9: infinity value should not be used",
		"5:9: unquoted strings are not allowed",
	}
	runLintTest(test)
}

func TestLiteralBlockScalar(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{
"include_newlines": |
            exactly as you see
            will appear these three
            lines of poetry
,
"fold_newlines": >
            this is really a
            single line of text
            despite appearances
}`
	test.want = []string{
		"2:20: literal block scalar '|' should not be used",
		"7:17: literal block scalar '>' should not be used",
	}
	runLintTest(test)
}

func TestBadCase1(t *testing.T) {
	test := newSuite(t)
	test.yaml = `# An employee record
name: Martin D'vloper
job: Developer
skill: Elite
employed: True
foods:
    - Apple
    - Orange
    - Strawberry
    - Mango
languages:
    perl: Elite
    python: Elite
    pascal: Lame
education: |
    4 GCSEs
    3 A-Levels
    BSc in the Internet of Things`
	test.want = []string{
		"2:5: use a flow object syntax {} instead",
		"2:7: unquoted strings are not allowed",
		"3:6: unquoted strings are not allowed",
		"4:8: unquoted strings are not allowed",
		"7:5: use a flow array syntax [] instead",
		"7:7: unquoted strings are not allowed",
		"8:7: unquoted strings are not allowed",
		"9:7: unquoted strings are not allowed",
		"10:7: unquoted strings are not allowed",
		"12:9: use a flow object syntax {} instead",
		"12:11: unquoted strings are not allowed",
		"13:13: unquoted strings are not allowed",
		"14:13: unquoted strings are not allowed",
		"15:12: literal block scalar '|' should not be used",
	}
	runLintTest(test)
}

func TestBadCase2(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{   Title: Blank lines denote

   paragraph breaks,
content: |-
   Or we
   can auto
   convert line breaks
   to save space
}
`
	test.want = []string{
		"1:12: unquoted strings are not allowed",
		"3:4: `paragraph breaks` is not a valid ES5.1 object key",
		"4:10: literal block scalar '|-' should not be used",
	}
	runLintTest(test)
}

func TestBadCase3(t *testing.T) {
	test := newSuite(t)
	test.yaml = `   Blank lines denote

   paragraph breaks
content: |-
   Or we
   can auto
   convert line breaks
   to save space
`
	test.want = []string{
		"3:4: found more than one document inside a file",
		"1:4: unquoted strings are not allowed",
		"3:4: unquoted strings are not allowed",
		"4:8: used a key-value outside of an object",
		"4:10: literal block scalar '|-' should not be used",
	}
	runLintTest(test)
}

func TestGoodCase1(t *testing.T) {
	test := newSuite(t)
	test.yaml = `{
  "foo": 3.5,
    # The indentation is broken on purpose/
      bar: {
        "nested": [],
      },
}`
	test.want = []string{}
	runLintTest(test)
}

func TestGoodCase2(t *testing.T) {
	test := newSuite(t)
	// Taken from the https://json5.org/.
	test.yaml = `{
  # comments
  unquoted: 'and you can quote me on that',
  singleQuotes: 'I can use "double quotes" here',
  lineBreaks: "Look, Mom! \
No \\n's!",
  hexadecimal: 0xdecaf,
  leadingDecimalPoint: .8675309, andTrailing: 8675309.,
  positiveSign: +1,
  trailingComma: 'in objects', andIn: ['arrays',],
  "backwardsCompatible": "with JSON",
}
`
	test.want = []string{}
	runLintTest(test)
}

type testSuite struct {
	t      *testing.T
	yaml   string
	config Config
	want   []string
}

func newSuite(t *testing.T) *testSuite {
	config := Config{}
	config.Allow.IdentObjKeys = true
	config.Allow.SingleQuoteStrings = true
	return &testSuite{t: t, config: config}
}

func runLintTest(suite *testSuite) {
	suite.t.Helper()
	const parseMode = 0
	f, err := parser.ParseBytes([]byte(suite.yaml), parseMode)
	if err != nil {
		suite.t.Fatalf("parse YAML: %v", err)
	}
	warnings := Run(&suite.config, f)
	have := make([]string, len(warnings))
	for i, w := range warnings {
		have[i] = fmt.Sprintf("%d:%d: %s", w.Line, w.Column, w.Text)
	}
	if diff := cmp.Diff(have, suite.want); diff != "" {
		suite.t.Errorf("output mismatch:\n%s", diff)
	}
}
