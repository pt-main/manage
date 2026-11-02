package manager

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pt-main/tap/go/color"
)

const todoMarker = "TODO:"

var lowerTodoMarker = strings.ToLower(todoMarker)

type CodeTask struct {
	Id       int
	About    string
	Done     bool
	CodeLine int
	File     string
}

func (ct *CodeTask) Type() TaskDataType { return CodeTaskT }

func (ct *CodeTask) Render(colored bool) string {
	res := fmt.Sprintf(
		"[?BBK]CODE:[?BYW]TODO [?CN]f:[[?GN]%v:%v[?CN]][?RT]:",
		ct.File, ct.CodeLine,
	)
	if colored {
		res = color.Set(res)
	} else {
		res = color.ReplaceColors(res)
	}
	return res + "\n" + ct.About
}

func (ct *CodeTask) GetId() int       { return ct.Id }
func (ct *CodeTask) IsDone() bool     { return ct.Done }
func (ct *CodeTask) GetAbout() string { return ct.About }

func (ct *CodeTask) GetDate() string { return "" }

func (ct *CodeTask) Add(string) bool      { return false }
func (ct *CodeTask) SetAbout(string) bool { return false }

func (ct *CodeTask) Do() error {
	if ct.File == "" {
		return errors.New("code task: empty file path")
	}

	info, err := os.Stat(ct.File)
	if err != nil {
		return fmt.Errorf("code task: stat %s: %w", ct.File, err)
	}

	data, err := os.ReadFile(ct.File)
	if err != nil {
		return fmt.Errorf("code task: read %s: %w", ct.File, err)
	}

	content := string(data)
	crlf := strings.Contains(content, "\r\n")
	if crlf {
		content = strings.ReplaceAll(content, "\r\n", "\n")
	}

	lines := strings.Split(content, "\n")
	idx := ct.CodeLine - 1
	if idx < 0 || idx >= len(lines) {
		return fmt.Errorf("code task: line %d out of range in %s", ct.CodeLine, ct.File)
	}

	newLine, ok := stripTodoMarker(lines[idx], ct.About)
	if !ok {
		return fmt.Errorf(
			"code task: TODO not found (or content changed) at %s:%d",
			ct.File, ct.CodeLine,
		)
	}
	lines[idx] = newLine

	result := strings.Join(lines, "\n")
	if crlf {
		result = strings.ReplaceAll(result, "\n", "\r\n")
	}
	if err := os.WriteFile(ct.File, []byte(result), info.Mode()); err != nil {
		return fmt.Errorf("code task: write %s: %w", ct.File, err)
	}

	ct.Done = true
	return nil
}

func stripTodoMarker(line, expectedAbout string) (string, bool) {
	positions := findAllTodoPositions(line)
	if len(positions) == 0 {
		return line, false
	}

	pos := positions[0]
	if expectedAbout != "" {
		pos = -1
		for _, p := range positions {
			if strings.TrimSpace(line[p+len(todoMarker):]) == expectedAbout {
				pos = p
				break
			}
		}
		if pos == -1 {
			return line, false
		}
	}

	before := strings.TrimRight(line[:pos], " \t")
	after := strings.TrimLeft(line[pos+len(todoMarker):], " \t")

	switch {
	case before == "":
		return after, true
	case after == "":
		return before, true
	default:
		return before + " " + after, true
	}
}

func findAllTodoPositions(line string) []int {
	lower := strings.ToLower(line)
	var positions []int
	for from := 0; ; {
		i := strings.Index(lower[from:], lowerTodoMarker)
		if i == -1 {
			return positions
		}
		positions = append(positions, from+i)
		from += i + len(lowerTodoMarker)
	}
}

var textExtensions = map[string]struct{}{
	".go": {}, ".py": {}, ".js": {}, ".jsx": {}, ".ts": {}, ".tsx": {},
	".rs": {}, ".c": {}, ".h": {}, ".cc": {}, ".cpp": {}, ".cxx": {},
	".hpp": {}, ".hxx": {}, ".java": {}, ".rb": {}, ".php": {}, ".cs": {},
	".swift": {}, ".kt": {}, ".kts": {}, ".sh": {}, ".bash": {}, ".zsh": {},
	".fish": {}, ".md": {}, ".txt": {}, ".rst": {}, ".org": {},
	".yaml": {}, ".yml": {}, ".toml": {}, ".json": {}, ".jsonc": {}, ".xml": {},
	".html": {}, ".htm": {}, ".css": {}, ".scss": {}, ".sass": {}, ".less": {},
	".sql": {}, ".lua": {}, ".pl": {}, ".pm": {}, ".r": {}, ".scala": {}, ".sc": {},
	".ex": {}, ".exs": {}, ".erl": {}, ".hrl": {}, ".hs": {}, ".lhs": {},
	".ml": {}, ".mli": {}, ".clj": {}, ".cljs": {}, ".vim": {},
	".dart": {}, ".groovy": {}, ".gradle": {}, ".cmake": {},
}

var skipDirs = map[string]struct{}{
	".git": {}, ".svn": {}, ".hg": {},
	"node_modules": {}, "vendor": {}, "target": {}, "build": {}, "dist": {},
	".idea": {}, ".vscode": {}, ".vs": {},
	"__pycache__": {}, ".venv": {}, "venv": {}, "env": {},
	"bin": {}, "obj": {}, ".gradle": {}, ".mypy_cache": {}, ".pytest_cache": {},
}

func FindTodosInDir(dir string) []*CodeTask {
	result := make([]*CodeTask, 0)

	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if _, skip := skipDirs[d.Name()]; skip && path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		if !isTextFile(path) {
			return nil
		}
		scanFileForTodos(path, &result)
		return nil
	})

	sort.Slice(result, func(i, j int) bool {
		if result[i].File != result[j].File {
			return result[i].File < result[j].File
		}
		if result[i].CodeLine != result[j].CodeLine {
			return result[i].CodeLine < result[j].CodeLine
		}
		return result[i].About < result[j].About
	})

	for i, t := range result {
		t.Id = i
	}

	return result
}

func isTextFile(path string) bool {
	if _, ok := textExtensions[strings.ToLower(filepath.Ext(path))]; ok {
		return true
	}
	switch strings.ToLower(filepath.Base(path)) {
	case "makefile", "dockerfile", "rakefile", "gemfile", "cargo.lock",
		"go.mod", "go.sum", "requirements.txt",
		".gitignore", ".dockerignore", ".editorconfig", ".env":
		return true
	}
	return false
}

func scanFileForTodos(path string, out *[]*CodeTask) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 2*1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		positions := findAllTodoPositions(line)
		if len(positions) == 0 {
			continue
		}
		pos := positions[0]
		about := strings.TrimSpace(line[pos+len(todoMarker):])
		*out = append(*out, &CodeTask{
			About:    about,
			CodeLine: lineNum,
			File:     path,
		})
	}
}
