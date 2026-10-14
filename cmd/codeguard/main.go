package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const maxFileLines = 600

var (
	allowedCommentPrefixes = []string{"BREAKING:", "LICENSE:", "NOSONAR", "go:"}
	skippedDirs            = []string{".git", ".claude", "vendor", "tmp", "bin", "node_modules"}
)

type violation struct {
	path string
	line int
	text string
}

func main() {
	violations, err := scan(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	for _, v := range violations {
		fmt.Printf("%s:%d  %s\n", v.path, v.line, v.text)
	}
	if len(violations) > 0 {
		fmt.Fprintf(os.Stderr, "\ncodeguard failed: %d violation(s). Comments are allowed only with %s; files stay at or below %d lines. See AGENTS.md.\n",
			len(violations), strings.Join(allowedCommentPrefixes, ", "), maxFileLines)
		os.Exit(1)
	}
	fmt.Println("codeguard passed: no comments, every Go file within the line limit.")
}

func scan(root string) ([]violation, error) {
	var out []violation
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if slices.Contains(skippedDirs, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		found, scanErr := scanFile(path)
		out = append(out, found...)
		return scanErr
	})
	return out, err
}

func scanFile(path string) ([]violation, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []violation
	if lines := bytes.Count(source, []byte("\n")); lines > maxFileLines {
		out = append(out, violation{path: path, line: lines, text: fmt.Sprintf("file has %d lines (max %d)", lines, maxFileLines)})
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, source, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	out = append(out, functionViolations(fset, path, file)...)
	for _, group := range file.Comments {
		for _, comment := range group.List {
			if !allowed(comment.Text) {
				out = append(out, violation{path: path, line: fset.Position(comment.Pos()).Line, text: preview(comment.Text)})
			}
		}
	}
	return out, nil
}

func allowed(text string) bool {
	body := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(text, "//"), "/*"))
	if strings.HasPrefix(text, "//go:") {
		return true
	}
	for _, prefix := range allowedCommentPrefixes {
		if strings.HasPrefix(body, prefix) {
			return true
		}
	}
	return false
}

func preview(text string) string {
	const limit = 80
	first, _, _ := strings.Cut(text, "\n")
	if len(first) > limit {
		return first[:limit-3] + "..."
	}
	return first
}
