# YAML5 - a JSON5* that is as widespread as YAML

[YAML5](https://github.com/quasilyte/yaml5) is a way of writing [YAML](https://yaml.org/) files that makes them look like [JSON5](https://json5.org/). In other words: it's not a new format.

This repository provides a useful tooling for you:

* `yaml5 lint` checks whether the YAML file complies with the YAML5 rules
* `yaml5 fmt` reads a YAML file and pretty-prints it as a YAML5 file (not implemented yet)

## YAML5 rules (tl;dr: it's JSON5 with a different comment syntax)

There is only one rule: your YAML file needs to be a valid JSON5 document.

The only difference is the single line comment syntax: use `#` instead of `//`.

Let's take an example from the [json5.org](https://json5.org/):

```yaml
{
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
```

The `yaml5 lint` tool does catch some corner cases, for example, it does check that you've used a valid unquoted object key following the `ES5.1` rules.

All YAML documents that use features that are not part of the JSON5 format will be reported as errors:

```
$ cat bad.yml
foo: 
  - .inf
  - 10

$ yaml5 lint bad.yml
bad.yml:1:4: used a key-value outside of an object
bad.yml:2:3: use a flow array syntax instead
bad.yml:2:5: infinity value should not be used
```

## Why you may want to use YAML5

* You were searching for a "JSON with comments and trailing commas", but the [JSON5](https://json5.org/) is not so popular
* You have a lot of YAML files but would like to keep them strict and free of the indentation-sensitive features

## Installation

```bash
# The easiest way to build it from sources:
go get -u -v github.com/quasilyte/yaml5/cmd/yaml5

# Build yaml5 binary from sources + embed the version info so the "yaml5 version"
# can give an appropriate output. Note that GOBIN folder is set to the
# current directory (the cloned yaml5 repository in this example); you can
# choose any other binary destination.
git clone https://github.com/quasilyte/yaml5.git
cd yaml5
GOBIN=$(pwd) make
```

> TODO: make a binary release.
