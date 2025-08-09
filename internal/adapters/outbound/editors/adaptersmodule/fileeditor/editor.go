package fileeditor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nduyhai/gocraft/internal/core/ports"
)

// Editor implements ports.AdaptersModuleEditor by editing
// <root>/internal/adapters/module.go in-place using robust string operations.
// All operations are idempotent.

type Editor struct{ root string }

func New(projectRoot string) *Editor { return &Editor{root: projectRoot} }

func (e *Editor) Ensure(alias, importPath, optionExpr string) error {
	if alias == "" || importPath == "" || optionExpr == "" {
		return fmt.Errorf("invalid ensure args: alias/importPath/optionExpr must be non-empty")
	}
	modFile := filepath.Join(e.root, "internal", "adapters", "module.go")
	b, err := os.ReadFile(modFile)
	if err != nil {
		return fmt.Errorf("read adapters/module.go: %w", err)
	}
	content := string(b)

	importLine := fmt.Sprintf("\t%s \"%s\"", alias, importPath)

	// Ensure import exists
	if !strings.Contains(content, importLine) {
		if strings.Contains(content, "import (") {
			// insert before first closing parenthesis of the import block
			scanner := bufio.NewScanner(bytes.NewReader([]byte(content)))
			var buf bytes.Buffer
			inserted := false
			for scanner.Scan() {
				line := scanner.Text()
				if strings.TrimSpace(line) == ")" && !inserted {
					buf.WriteString(importLine + "\n")
					inserted = true
				}
				buf.WriteString(line + "\n")
			}
			if err := scanner.Err(); err != nil {
				return err
			}
			content = buf.String()
		} else if strings.Contains(content, "import \"") {
			// convert single import to block containing both existing and new
			scanner := bufio.NewScanner(bytes.NewReader([]byte(content)))
			var buf bytes.Buffer
			converted := false
			for scanner.Scan() {
				line := scanner.Text()
				trim := strings.TrimSpace(line)
				if strings.HasPrefix(trim, "import \"") && !converted {
					start := strings.Index(line, "\"")
					end := strings.LastIndex(line, "\"")
					existing := ""
					if start >= 0 && end > start {
						existing = line[start : end+1]
					}
					buf.WriteString("import (\n\t" + existing + "\n" + importLine + "\n)\n")
					converted = true
				} else {
					buf.WriteString(line + "\n")
				}
			}
			if err := scanner.Err(); err != nil {
				return err
			}
			content = buf.String()
		} else {
			// no import block; add after package line
			content = strings.Replace(content, "package adapters\n", "package adapters\n\nimport (\n"+importLine+"\n)\n\n", 1)
		}
	}

	// Ensure fx.Option reference exists
	if !strings.Contains(content, optionExpr) {
		needle := "\t\tfx.Options(Extras...),"
		if strings.Contains(content, needle) {
			content = strings.Replace(content, needle, needle+"\n\t\t"+optionExpr+",", 1)
		} else {
			// Fallback: insert before the closing ')' of return fx.Options(
			blockStart := strings.Index(content, "return fx.Options(")
			if blockStart >= 0 {
				s := content[blockStart:]
				open := 0
				idx := -1
				for i, ch := range s {
					if ch == '(' {
						open++
					}
					if ch == ')' {
						open--
						if open == 0 {
							idx = i
							break
						}
					}
				}
				if idx > 0 {
					insertAt := blockStart + idx
					content = content[:insertAt] + "\n\t\t" + optionExpr + "," + content[insertAt:]
				}
			}
		}
	}

	return os.WriteFile(modFile, []byte(content), 0o644)
}

var _ ports.AdaptersModuleEditor = (*Editor)(nil)
