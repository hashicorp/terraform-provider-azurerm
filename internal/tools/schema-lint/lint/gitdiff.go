// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lint

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Changes records, per absolute file path, the set of line numbers added or
// changed relative to a base git ref.
type Changes struct {
	root  string
	files map[string]map[int]bool
}

// Added reports whether the given (absolute) file and line was added/changed.
func (c *Changes) Added(file string, line int) bool {
	if c == nil {
		return false
	}
	abs, err := filepath.Abs(file)
	if err != nil {
		abs = file
	}
	lines, ok := c.files[abs]
	return ok && lines[line]
}

// AddedInRange reports whether any line in [start, end] of the given (absolute)
// file was added/changed.
func (c *Changes) AddedInRange(file string, start, end int) bool {
	if c == nil {
		return false
	}
	abs, err := filepath.Abs(file)
	if err != nil {
		abs = file
	}
	lines, ok := c.files[abs]
	if !ok {
		return false
	}
	if start > end {
		start, end = end, start
	}
	for l := start; l <= end; l++ {
		if lines[l] {
			return true
		}
	}
	return false
}

// ChangedFiles returns the absolute paths of files with added/changed lines.
func (c *Changes) ChangedFiles() []string {
	out := make([]string, 0, len(c.files))
	for f := range c.files {
		out = append(out, f)
	}
	return out
}

// Diff computes the lines added since baseRef within root, restricted to the
// given pathspecs. It diffs the working tree against the merge base so that both
// committed and uncommitted changes on the branch are considered.
func Diff(root, baseRef string, pathspecs []string) (*Changes, error) {
	base := baseRef
	if mb, err := git(root, "merge-base", baseRef, "HEAD"); err == nil {
		if v := strings.TrimSpace(mb); v != "" {
			base = v
		}
	}

	args := []string{"diff", "--no-color", "--unified=0", base, "--"}
	args = append(args, pathspecs...)
	out, err := git(root, args...)
	if err != nil {
		return nil, fmt.Errorf("git diff failed: %w", err)
	}

	return parseDiff(root, out)
}

func parseDiff(root, diff string) (*Changes, error) {
	c := &Changes{root: root, files: map[string]map[int]bool{}}

	var curFile string
	var curLine int
	scanner := bufio.NewScanner(strings.NewReader(diff))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "+++ "):
			curFile = parseDiffPath(root, strings.TrimPrefix(line, "+++ "))
		case strings.HasPrefix(line, "--- "):
			// old-file marker; ignored.
		case strings.HasPrefix(line, "@@"):
			start, ok := parseHunkNewStart(line)
			if ok {
				curLine = start
			}
		case strings.HasPrefix(line, "+"):
			if curFile != "" {
				if c.files[curFile] == nil {
					c.files[curFile] = map[int]bool{}
				}
				c.files[curFile][curLine] = true
			}
			curLine++
		}
	}
	return c, scanner.Err()
}

// parseDiffPath turns a diff "+++" target (e.g. "b/internal/x.go") into an
// absolute path, or "" for /dev/null.
func parseDiffPath(root, p string) string {
	p = strings.TrimSpace(p)
	if p == "/dev/null" {
		return ""
	}
	if i := strings.IndexByte(p, '/'); i >= 0 {
		p = p[i+1:] // strip the a/ or b/ prefix
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(root, p)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

// parseHunkNewStart extracts the new-side start line from a hunk header such as
// "@@ -1,0 +23,4 @@".
func parseHunkNewStart(line string) (int, bool) {
	plus := strings.IndexByte(line, '+')
	if plus < 0 {
		return 0, false
	}
	rest := line[plus+1:]
	end := strings.IndexAny(rest, ", ")
	if end >= 0 {
		rest = rest[:end]
	}
	n, err := strconv.Atoi(rest)
	if err != nil {
		return 0, false
	}
	return n, true
}

func git(root string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return out.String(), nil
}
