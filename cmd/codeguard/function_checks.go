package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
)

const minDuplicateStatements = 2

func functionViolations(fset *token.FileSet, path string, file *ast.File) []violation {
	var out []violation
	firstBody := map[string]string{}
	isTest := strings.HasSuffix(path, "_test.go")
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		line := fset.Position(fn.Pos()).Line
		if len(fn.Body.List) == 0 && !isTest {
			out = append(out, violation{path: path, line: line, text: fmt.Sprintf("function %s is empty (Sonar S1186)", fn.Name.Name)})
		}
		if len(fn.Body.List) < minDuplicateStatements {
			continue
		}
		body := signature(fset, fn) + render(fset, fn.Body)
		if original, seen := firstBody[body]; seen {
			out = append(out, violation{path: path, line: line, text: fmt.Sprintf("function %s duplicates %s (Sonar S4144)", fn.Name.Name, original)})
			continue
		}
		firstBody[body] = fn.Name.Name
	}
	return out
}

func signature(fset *token.FileSet, fn *ast.FuncDecl) string {
	if fn.Recv == nil {
		return render(fset, fn.Type)
	}
	return render(fset, fn.Recv.List[0].Type) + render(fset, fn.Type)
}

func render(fset *token.FileSet, node ast.Node) string {
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, fset, node)
	return buf.String()
}
