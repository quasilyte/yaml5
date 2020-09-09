package lintcmd

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/goccy/go-yaml/parser"
	"github.com/quasilyte/yaml5/internal/lint"
)

func Main() (int, error) {
	var args arguments
	flag.BoolVar(&args.singleQuoteStrings, "allow-single-quote-strings", true,
		`Whether to allow single quoted strings like 'abc'`)
	flag.BoolVar(&args.unquotedKeys, "allow-unquoted-keys", true,
		`Whether to allow identifier (unquoted) object keys like {key: "val"}`)

	flag.Parse()
	args.targets = flag.Args()

	code, err := mainNoExit(&args)
	return int(code), err
}

type arguments struct {
	targets []string

	singleQuoteStrings bool
	unquotedKeys       bool
}

type exitCode int

const (
	codeOK exitCode = iota
	codeWarning
	codeErr
)

func mainNoExit(args *arguments) (exitCode, error) {
	numReported := 0
	for _, target := range args.targets {
		warnings, err := checkTarget(args, target)
		if err != nil {
			return codeErr, err
		}
		numReported += len(warnings)
		for _, w := range warnings {
			fmt.Printf("%s:%d:%d: %s\n", target, w.Line, w.Column, w.Text)
		}
	}

	if numReported != 0 {
		return codeWarning, nil
	}
	return codeOK, nil
}

func checkTarget(args *arguments, target string) ([]lint.Warning, error) {
	yamlBytes, err := ioutil.ReadFile(target)
	if err != nil {
		return nil, err
	}

	const parseMode = 0
	f, err := parser.ParseBytes(yamlBytes, parseMode)
	if err != nil {
		return nil, fmt.Errorf("parse YAML: %v", err)
	}

	config := lint.Config{}
	config.Allow.IdentObjKeys = args.unquotedKeys
	return lint.Run(&config, f), nil
}
