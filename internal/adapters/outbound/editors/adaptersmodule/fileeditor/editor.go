package fileeditor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/nduyhai/gocraft/internal/core/ports"
	"golang.org/x/tools/imports"
)

// Editor implements ports.AdaptersModuleEditor by editing
// <root>/internal/platform/di/root.go using AST and goimports.
// All operations are idempotent and tolerate missing files.

type Editor struct{ root string }

func New(projectRoot string) *Editor { return &Editor{root: projectRoot} }

func (e *Editor) Ensure(alias, importPath, optionExpr string) error {
	if alias == "" || importPath == "" || optionExpr == "" {
		return fmt.Errorf("invalid ensure args: alias/importPath/optionExpr must be non-empty")
	}
	return e.ensureInFile(
		filepath.Join(e.root, "internal", "platform", "di", "root.go"),
		alias, importPath, optionExpr,
	)
}

// ensureInFile ensures an import alias/path and fx option expression exist in the given file using AST.
func (e *Editor) ensureInFile(filePath, alias, importPath, optionExpr string) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		// Missing file: treat as no-op.
		return nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, b, parser.ParseComments)
	if err != nil {
		// If parsing fails, do not modify (avoid corruption); surface error to caller
		return err
	}

	// Ensure import exists with alias
	ensureImport := func() {
		for _, imp := range f.Imports {
			pathVal := strings.Trim(imp.Path.Value, "\"")
			name := ""
			if imp.Name != nil {
				name = imp.Name.Name
			}
			if pathVal == importPath && name == alias {
				return // already present
			}
		}
		// Not present: add
		imp := &ast.ImportSpec{
			Name: &ast.Ident{Name: alias},
			Path: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", importPath)},
		}
		// Create or find import decl
		var gen *ast.GenDecl
		for _, d := range f.Decls {
			if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
				gen = gd
				break
			}
		}
		if gen == nil {
			gen = &ast.GenDecl{Tok: token.IMPORT, Lparen: token.NoPos, Rparen: token.NoPos}
			f.Decls = append([]ast.Decl{gen}, f.Decls...)
		}
		gen.Specs = append(gen.Specs, imp)
		if gen.Lparen == token.NoPos && len(gen.Specs) > 1 {
			// ensure it's a block import
			gen.Lparen = token.Pos(1)
			gen.Rparen = token.Pos(1)
		}
	}
	ensureImport()

	// Ensure optionExpr exists inside return fx.Options(...)
	exprToAdd, err := parser.ParseExpr(optionExpr)
	if err != nil {
		return fmt.Errorf("invalid optionExpr: %w", err)
	}

	foundAndEnsured := false
	ast.Inspect(f, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		for _, res := range ret.Results {
			call, ok := res.(*ast.CallExpr)
			if !ok {
				continue
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel == nil || sel.Sel.Name != "Options" {
				continue
			}
			// naive check: qualifier exists (fx). We don't enforce package name here.
			// ensure not duplicate
			var buf bytes.Buffer
			_ = printer.Fprint(&buf, fset, exprToAdd)
			needle := buf.String()
			already := false
			for _, a := range call.Args {
				buf.Reset()
				_ = printer.Fprint(&buf, fset, a)
				if buf.String() == needle {
					already = true
					break
				}
			}
			if !already {
				call.Args = append(call.Args, exprToAdd)
			}
			foundAndEnsured = true
			return false
		}
		return true
	})
	// If not found, do nothing (idempotent across templates)
	_ = foundAndEnsured

	// Print and run goimports
	var out bytes.Buffer
	if err := printer.Fprint(&out, fset, f); err != nil {
		return err
	}
	processed, err := imports.Process(filePath, out.Bytes(), &imports.Options{Comments: true, TabWidth: 8, Fragment: false})
	if err != nil {
		// Even if imports fails, write the printer output to avoid losing changes
		processed = out.Bytes()
	}
	return os.WriteFile(filePath, processed, 0o644)
}

var _ ports.AdaptersModuleEditor = (*Editor)(nil)
