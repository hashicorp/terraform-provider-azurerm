// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// harper-check runs the Harper grammar checker over the repository and reports what it finds against the line that has to be fixed.
//
// Harper only reads comments in Go files, so schema descriptions are pulled out and checked as text, then reported at their own line. It has two steps:
//
//	harper-check code   comments, markdown outside the website docs, shell scripts and schema descriptions
//	harper-check docs   the website docs, reported at the schema description the text came from; text with no description yet is skipped
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// grammarRules are the Harper rules that are right often enough to fail a build on. To try another, run `harper-cli lint --only <rule>` over the
// repository first and read what it reports.
var grammarRules = []string{
	"AnA",           // "a integer", "an Certificate"
	"RepeatedWords", // "the the"
	"ItsPossessive", // "back to it's default"
}

// proseRules are run on schema descriptions and the website docs, which are sentences written for users rather than notes between developers.
var proseRules = append([]string{
	"SpellCheck",           // words not in Harper's dictionary or .harper/dictionary.txt
	"InflectedVerbAfterTo", // "to enabled"
}, grammarRules...)

// generatedMarkers are how generated Go files say so in their first lines; they are skipped, as the fix belongs in the generator.
var generatedMarkers = [][]byte{[]byte("Code generated"), []byte("this file is generated")}

// finding is one problem Harper reported, placed at the file and line to fix.
type finding struct {
	path       string
	line       int
	column     int
	rule       string
	message    string
	text       string
	suggestion string
	from       string // for a description found through the docs, the docs line it was seen on
}

// description is the text of one schema or resource description and where its string starts in the Go file.
type description struct {
	path   string
	line   int
	column int
	text   string
	exact  bool // the source string has no escapes, so a position in text is the same position in the file
}

func main() {
	harper := flag.String("harper", "harper-cli", "path to the harper-cli binary")
	dictionary := flag.String("dictionary", ".harper/dictionary.txt", "file of extra words Harper should accept, one per line")
	ignore := flag.String("ignore", ".harper/ignore.txt", "file of findings to leave alone, each written as `path: rule: the line of the file`")
	all := flag.Bool("all", false, "docs: also report findings in text that did not come from a schema description")
	flag.Parse()

	if flag.NArg() != 1 || (flag.Arg(0) != "code" && flag.Arg(0) != "docs") {
		fmt.Fprintln(os.Stderr, "usage: harper-check [-harper path] [-dictionary path] [-all] code|docs")
		os.Exit(2)
	}

	check := checker{harper: *harper, dictionary: *dictionary}

	var findings []finding
	var err error
	if flag.Arg(0) == "code" {
		findings, err = check.code()
	} else {
		findings, err = check.docs(*all)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "harper-check: %+v\n", err)
		os.Exit(2)
	}

	findings, err = withoutIgnored(findings, *ignore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "harper-check: %+v\n", err)
		os.Exit(2)
	}

	if len(findings) == 0 {
		return
	}

	sort.Slice(findings, func(i, j int) bool {
		a, b := findings[i], findings[j]
		if a.path != b.path {
			return a.path < b.path
		}

		if a.line != b.line {
			return a.line < b.line
		}

		return a.column < b.column
	})

	for _, f := range findings {
		fmt.Println(f)
	}

	fmt.Printf("\n%d problems found. Fix the text, add a correct word Harper does not know to %s, or list a wrong report in %s.\n", len(findings), *dictionary, *ignore)
	os.Exit(1)
}

func (f finding) String() string {
	s := fmt.Sprintf("%s:%d:%d: %s: %s %q", f.path, f.line, f.column, f.rule, f.message, f.text)
	if f.suggestion != "" {
		s += " (" + f.suggestion + ")"
	}

	if f.from != "" {
		s += " [seen in " + f.from + "]"
	}

	return s
}

// withoutIgnored drops the findings listed in the ignore file. An entry names the file, the rule and the whole line the finding is on, so it keeps
// working when the line moves and stops when the line is edited.
func withoutIgnored(findings []finding, ignoreFile string) ([]finding, error) {
	content, err := os.ReadFile(ignoreFile)
	if err != nil {
		if os.IsNotExist(err) {
			return findings, nil
		}

		return nil, err
	}

	ignored := map[string]bool{}
	for _, entry := range strings.Split(string(content), "\n") {
		if entry = strings.TrimSpace(entry); entry != "" && !strings.HasPrefix(entry, "#") {
			ignored[entry] = true
		}
	}

	kept := make([]finding, 0, len(findings))
	lines := map[string][]string{}
	for _, f := range findings {
		if _, ok := lines[f.path]; !ok {
			source, err := os.ReadFile(f.path)
			if err != nil {
				return nil, err
			}

			lines[f.path] = strings.Split(string(source), "\n")
		}

		line := ""
		if f.line >= 1 && f.line <= len(lines[f.path]) {
			line = strings.TrimSpace(lines[f.path][f.line-1])
		}

		if !ignored[fmt.Sprintf("%s: %s: %s", f.path, f.rule, line)] {
			kept = append(kept, f)
		}
	}

	return kept, nil
}

type checker struct {
	harper     string
	dictionary string
}

// code checks everything except the website docs: comments in Go files, markdown, shell scripts and schema descriptions.
func (c checker) code() ([]finding, error) {
	tracked, err := trackedFiles()
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(tracked))
	for _, path := range tracked {
		if checkedAsCode(path) {
			files = append(files, path)
		}
	}

	findings, err := c.lint(files, grammarRules, "us")
	if err != nil {
		return nil, err
	}

	descriptions, err := schemaDescriptions(tracked)
	if err != nil {
		return nil, err
	}

	inDescriptions, err := c.lintDescriptions(descriptions)
	if err != nil {
		return nil, err
	}

	return append(findings, inDescriptions...), nil
}

// docs checks the website docs. A finding in a sentence that is also a schema description is reported at that description, as that is where the docs
// will be generated from. Other findings are skipped unless all is set, since most descriptions have not been written yet.
func (c checker) docs(all bool) ([]finding, error) {
	tracked, err := trackedFiles()
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, path := range tracked {
		if strings.HasPrefix(path, "website/docs/") && strings.HasSuffix(path, ".markdown") {
			files = append(files, path)
		}
	}

	inDocs, err := c.lintProse(files)
	if err != nil {
		return nil, err
	}

	descriptions, err := schemaDescriptions(tracked)
	if err != nil {
		return nil, err
	}

	bySentence := map[string][]sentenceOf{}
	for _, d := range descriptions {
		for _, s := range sentences(d.text) {
			bySentence[s.text] = append(bySentence[s.text], sentenceOf{description: d, start: s.start})
		}
	}

	findings := make([]finding, 0)
	seen := map[string]bool{}
	skipped := 0
	lines := map[string][]string{}
	for _, f := range inDocs {
		if _, ok := lines[f.path]; !ok {
			content, err := os.ReadFile(f.path)
			if err != nil {
				return nil, err
			}

			lines[f.path] = strings.Split(string(content), "\n")
		}

		sentence, offset, ok := sentenceAt(lines[f.path], f.line, f.column)
		matches := bySentence[sentence]
		if !ok || len(matches) == 0 {
			if all {
				findings = append(findings, f)
			} else {
				skipped++
			}

			continue
		}

		for _, m := range matches {
			traced := f
			traced.path, traced.line, traced.column = m.description.path, m.description.line, m.description.column
			if m.description.exact {
				traced.column += 1 + m.start + offset
			}

			traced.from = fmt.Sprintf("%s:%d", f.path, f.line)

			// the same sentence is often on several docs pages, so report each description once
			key := fmt.Sprintf("%s:%d:%d:%s", traced.path, traced.line, traced.column, traced.rule)
			if !seen[key] {
				seen[key] = true
				findings = append(findings, traced)
			}
		}
	}

	if skipped > 0 {
		fmt.Printf("%d problems in docs text that is not a schema description yet were skipped; run with -all to list them.\n", skipped)
	}

	return findings, nil
}

// trackedFiles lists the files git tracks, so build output and downloaded tools are never checked.
func trackedFiles() ([]string, error) {
	out, err := exec.Command("git", "ls-files", "-z").Output()
	if err != nil {
		return nil, fmt.Errorf("listing tracked files: %+v", err)
	}

	return strings.Split(strings.TrimRight(string(out), "\x00"), "\x00"), nil
}

// checkedAsCode reports whether the code step reads this file: Go, markdown and shell that people write by hand. Harper treats files it has no
// parser for, such as YAML and Terraform, as plain English and reports their syntax, so those are left out.
func checkedAsCode(path string) bool {
	if strings.HasPrefix(path, "vendor/") || strings.Contains(path, "/vendor/") || strings.Contains(path, "/testdata/") {
		return false
	}

	switch filepath.Ext(path) {
	case ".go":
		return !generated(path)
	case ".sh":
		return true
	case ".md", ".markdown":
		// the docs have their own step, and released changelog entries are history that is not rewritten
		return !strings.HasPrefix(path, "website/docs/") && !strings.HasPrefix(filepath.Base(path), "CHANGELOG")
	}

	return false
}

func generated(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	head := make([]byte, 1024)
	n, _ := f.Read(head)
	for _, marker := range generatedMarkers {
		if bytes.Contains(head[:n], marker) {
			return true
		}
	}

	return false
}

// harperFile and harperLint are the parts of Harper's JSON output that are used.
type harperFile struct {
	File  string       `json:"file"`
	Lints []harperLint `json:"lints"`
}

type harperLint struct {
	Rule        string   `json:"rule"`
	Line        int      `json:"line"`
	Column      int      `json:"column"`
	Message     string   `json:"message"`
	Suggestions []string `json:"suggestions"`
	MatchedText string   `json:"matched_text"`
}

// lintProse runs the prose rules over files. The repository mixes British and American spelling, and Harper only accepts one at a time, so spelling
// is checked both ways and a word is only reported when neither accepts it.
func (c checker) lintProse(files []string) ([]finding, error) {
	american, err := c.lint(files, proseRules, "us")
	if err != nil {
		return nil, err
	}

	british, err := c.lint(files, []string{"SpellCheck"}, "uk")
	if err != nil {
		return nil, err
	}

	misspeltInBritish := map[string]bool{}
	for _, f := range british {
		misspeltInBritish[fmt.Sprintf("%s:%d:%d", f.path, f.line, f.column)] = true
	}

	findings := make([]finding, 0, len(american))
	for _, f := range american {
		if f.rule != "SpellCheck" || misspeltInBritish[fmt.Sprintf("%s:%d:%d", f.path, f.line, f.column)] {
			findings = append(findings, f)
		}
	}

	return findings, nil
}

// lint runs Harper over files with only the given rules on, reading them as the given dialect of English, and returns what it found.
func (c checker) lint(files []string, rules []string, dialect string) ([]finding, error) {
	if len(files) == 0 {
		return nil, nil
	}

	dir, err := os.MkdirTemp("", "harper-check")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	// Harper names a file in its output by its last part only, so each is linked under a number that leads back to the real path
	links := make([]string, 0, len(files))
	paths := map[string]string{}
	for i, path := range files {
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		name := fmt.Sprintf("%07d%s", i, filepath.Ext(path))
		if err := os.Symlink(abs, filepath.Join(dir, name)); err != nil {
			return nil, err
		}

		links = append(links, filepath.Join(dir, name))
		paths[name] = path
	}

	// an empty folder for per-file word lists, and the repository's own word list, so nothing is picked up from the machine it runs on
	fileDictionaries := filepath.Join(dir, "file-dictionaries")
	if err := os.Mkdir(fileDictionaries, 0o755); err != nil {
		return nil, err
	}

	findings := make([]finding, 0)
	const batch = 500
	for start := 0; start < len(links); start += batch {
		end := min(start+batch, len(links))

		args := []string{"lint", "--format", "json", "--quiet", "--dialect", dialect, "--only", strings.Join(rules, ","), "--user-dict-path", c.dictionary, "--file-dict-path", fileDictionaries}
		cmd := exec.Command(c.harper, append(args, links[start:end]...)...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		// Harper exits non-zero when it finds anything, so only output that cannot be read is a failure
		out, runErr := cmd.Output()
		var results []harperFile
		if err := json.Unmarshal(out, &results); err != nil {
			return nil, fmt.Errorf("running %s: %v: %s", c.harper, errors.Join(runErr, err), stderr.String())
		}

		for _, result := range results {
			for _, l := range result.Lints {
				f := finding{path: paths[result.File], line: l.Line, column: l.Column, rule: l.Rule, message: l.Message, text: l.MatchedText}
				if len(l.Suggestions) > 0 {
					f.suggestion = l.Suggestions[0]
				}

				findings = append(findings, f)
			}
		}
	}

	return findings, nil
}

// lintDescriptions checks schema descriptions by writing them out as a markdown file, one to a paragraph, and placing each finding back at the
// description's own line in the Go file.
func (c checker) lintDescriptions(descriptions []description) ([]finding, error) {
	if len(descriptions) == 0 {
		return nil, nil
	}

	dir, err := os.MkdirTemp("", "harper-check-descriptions")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	var text strings.Builder
	for _, d := range descriptions {
		text.WriteString(strings.Join(strings.Fields(d.text), " "))
		text.WriteString("\n\n")
	}

	file := filepath.Join(dir, "descriptions.md")
	if err := os.WriteFile(file, []byte(text.String()), 0o644); err != nil {
		return nil, err
	}

	found, err := c.lintProse([]string{file})
	if err != nil {
		return nil, err
	}

	findings := make([]finding, 0, len(found))
	for _, f := range found {
		// each description takes two lines of the file: its text and a blank line
		d := descriptions[(f.line-1)/2]
		f.path, f.line = d.path, d.line
		if d.exact {
			f.column += d.column
		} else {
			f.column = d.column
		}

		findings = append(findings, f)
	}

	return findings, nil
}

// schemaDescriptions finds every description written as a plain string in non-test Go code under internal: a `Description: "..."` field, and a
// typed resource's `Description()` method that returns one.
func schemaDescriptions(tracked []string) ([]description, error) {
	descriptions := make([]description, 0)
	fset := token.NewFileSet()
	for _, path := range tracked {
		if !strings.HasPrefix(path, "internal/") || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || strings.Contains(path, "/vendor/") {
			continue
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if !bytes.Contains(src, []byte("Description")) {
			continue
		}

		file, err := parser.ParseFile(fset, path, src, parser.SkipObjectResolution)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %+v", path, err)
		}

		ast.Inspect(file, func(n ast.Node) bool {
			var lit *ast.BasicLit
			switch n := n.(type) {
			case *ast.KeyValueExpr:
				if key, ok := n.Key.(*ast.Ident); ok && key.Name == "Description" {
					lit, _ = n.Value.(*ast.BasicLit)
				}
			case *ast.FuncDecl:
				if n.Name.Name == "Description" && n.Recv != nil && n.Body != nil && len(n.Body.List) == 1 {
					if ret, ok := n.Body.List[0].(*ast.ReturnStmt); ok && len(ret.Results) == 1 {
						lit, _ = ret.Results[0].(*ast.BasicLit)
					}
				}
			}

			if lit == nil || lit.Kind != token.STRING {
				return true
			}

			text, err := strconv.Unquote(lit.Value)
			if err != nil || strings.TrimSpace(text) == "" {
				return true
			}

			pos := fset.Position(lit.Pos())
			descriptions = append(descriptions, description{path: path, line: pos.Line, column: pos.Column, text: text, exact: lit.Value[1:len(lit.Value)-1] == text})

			return true
		})
	}

	return descriptions, nil
}

// sentence is one sentence of a longer text and where it starts in it.
type sentence struct {
	text  string
	start int
}

// sentenceOf is a sentence's place in the description it is part of.
type sentenceOf struct {
	description description
	start       int
}

// sentences splits text after each full stop, question mark or exclamation mark that is followed by a space. It does not need to be right about
// what a sentence is, only to cut a description and the docs line copied from it in the same places.
func sentences(text string) []sentence {
	out := make([]sentence, 0)
	start := 0
	for i := 0; i < len(text); i++ {
		last := i == len(text)-1
		if !last && !(strings.ContainsRune(".?!", rune(text[i])) && text[i+1] == ' ') {
			continue
		}

		// leading spaces are not part of the sentence, so its start moves past them
		raw := text[start : i+1]
		trimmed := strings.TrimLeft(raw, " ")
		if s := strings.TrimSpace(trimmed); s != "" {
			out = append(out, sentence{text: s, start: start + len(raw) - len(trimmed)})
		}

		start = i + 1
	}

	return out
}

// argumentPrefixRegex matches the start of a docs line that documents a property, up to where its description begins. In the line
// "* `name` - (Required) The name of the Example." it matches everything before "The name".
var argumentPrefixRegex = regexp.MustCompile("^\\s*[*-]\\s+`[^`]+`\\s+-\\s+(\\((Required|Optional)\\)\\s+)?")

// sentenceAt returns the sentence of a docs line that holds the given column, and how far into that sentence the column is.
func sentenceAt(lines []string, line int, column int) (string, int, bool) {
	if line < 1 || line > len(lines) {
		return "", 0, false
	}

	text := lines[line-1]
	prefix := len(argumentPrefixRegex.FindString(text))
	at := column - 1 - prefix
	if at < 0 {
		return "", 0, false
	}

	for _, s := range sentences(text[prefix:]) {
		if at >= s.start && at < s.start+len(s.text) {
			return s.text, at - s.start, true
		}
	}

	return "", 0, false
}
