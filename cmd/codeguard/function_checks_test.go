package main

import (
	"go/parser"
	"go/token"
	"slices"
	"testing"
)

const functionSource = `package sample

type box struct{}

func empty() {}

func first(a int) int {
	b := a + 1
	return b
}

func second(a int) int {
	b := a + 1
	return b
}

func other(a string) string {
	b := a + "1"
	return b
}

func (box) first(a int) int {
	b := a + 1
	return b
}

func short(a int) int { return a }

func shortCopy(a int) int { return a }
`

func violationsOf(t *testing.T, path string) []string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, functionSource, 0)
	if err != nil {
		t.Fatal(err)
	}
	var texts []string
	for _, v := range functionViolations(fset, path, file) {
		texts = append(texts, v.text)
	}
	return texts
}

func expectViolations(t *testing.T, path string, want ...string) {
	t.Helper()
	if got := violationsOf(t, path); !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFunctionViolations(t *testing.T) {
	expectViolations(t, "sample.go",
		"function empty is empty (Sonar S1186)",
		"function second duplicates first (Sonar S4144)",
	)
}

func TestEmptyFunctionsAllowedInTests(t *testing.T) {
	expectViolations(t, "sample_test.go", "function second duplicates first (Sonar S4144)")
}
