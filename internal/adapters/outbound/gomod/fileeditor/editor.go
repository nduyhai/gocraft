package fileeditor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Editor provides a naive text-based go.mod editor that works without invoking external tools.
// It supports adding require entries and replace directives in an idempotent way.
// Tidy is a no-op.
//
// This is intentionally simple and resilient for generator use; it won't fully parse go.mod semantics.

type Editor struct {
	root string
}

func New(projectRoot string) *Editor { return &Editor{root: projectRoot} }

func (e *Editor) goModPath() string { return filepath.Join(e.root, "go.mod") }

// Add ensures the require entry exists with the specified version.
func (e *Editor) Add(module, version string) error {
	if module == "" {
		return fmt.Errorf("module path is empty")
	}
	path := e.goModPath()
	b, _ := os.ReadFile(path)
	content := string(b)

	// Check if already present
	if hasRequire(content, module) {
		// If present but without version match, we leave as-is to avoid breaking constraints.
		return nil
	}

	line := fmt.Sprintf("\t%s %s", module, version)
	if strings.Contains(content, "require (") {
		// Insert inside the first require block before the closing ')'
		idx := strings.Index(content, "require (")
		end := strings.Index(content[idx:], ")")
		if end > 0 {
			insertPos := idx + end
			content = content[:insertPos] + "\n" + line + "\n" + content[insertPos:]
		} else {
			// malformed; append a new block at end
			if !strings.HasSuffix(content, "\n") {
				content += "\n"
			}
			content += "require (\n" + line + "\n)\n"
		}
	} else {
		// No require block; try find a single-line require or add a block
		scanner := bufio.NewScanner(strings.NewReader(content))
		var hasSingle bool
		for scanner.Scan() {
			lineTxt := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(lineTxt, "require ") && strings.Count(lineTxt, " ") >= 2 {
				hasSingle = true
				break
			}
		}
		if hasSingle {
			// Convert to block by appending new block at end (simplest approach)
			if !strings.HasSuffix(content, "\n") {
				content += "\n"
			}
			content += "require (\n" + line + "\n)\n"
		} else {
			if !strings.HasSuffix(content, "\n") {
				content += "\n"
			}
			content += "require (\n" + line + "\n)\n"
		}
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// Replace adds or updates a replace directive.
func (e *Editor) Replace(oldPath, newPath string) error {
	if oldPath == "" || newPath == "" {
		return fmt.Errorf("replace paths must be non-empty")
	}
	path := e.goModPath()
	b, _ := os.ReadFile(path)
	content := string(b)

	re := regexp.MustCompile(`(?m)^replace\s+` + regexp.QuoteMeta(oldPath) + `\s*=>.*$`)
	if re.MatchString(content) {
		content = re.ReplaceAllString(content, "replace "+oldPath+" => "+newPath)
	} else {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += "replace " + oldPath + " => " + newPath + "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// Tidy is a no-op placeholder.
func (e *Editor) Tidy() error { return nil }

func hasRequire(content, module string) bool {
	re := regexp.MustCompile(`(?m)^\s*require\s*\(|^\s*require\s+`)
	if !re.MatchString(content) {
		return false
	}
	// simple contains check for module in a require line
	scanner := bufio.NewScanner(bytes.NewReader([]byte(content)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "require ") {
			// single-line form: require <mod> <ver>
			if strings.Contains(line, " "+module+" ") || strings.HasSuffix(line, " "+module) {
				return true
			}
		}
		if strings.HasPrefix(line, ")") {
			continue
		}
		// inside block: <mod> <ver>
		if strings.HasPrefix(line, module+" ") || strings.HasSuffix(line, " "+module) {
			return true
		}
	}
	return false
}
