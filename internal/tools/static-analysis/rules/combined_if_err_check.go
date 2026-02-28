// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var _ Rule = CombinedIfErrCheck{}

// matchUnderscoreErr matches lines like "_, err :=" or "_, err =" but NOT
// "_, _, err :=" or "resp, _, err :=" (multi-return) or "for _, err := range" (loops).
var (
	matchUnderscoreErr = regexp.MustCompile(`^\s*_, err :?= `)
	matchMultiReturn   = regexp.MustCompile(`^\s*\w+.*,\s*_, err :?= `)
	matchForRange      = regexp.MustCompile(`^\s*for\s+`)
	matchIfErrNil      = regexp.MustCompile(`^\s*if err != nil`)
)

type CombinedIfErrCheck struct{}

func (r CombinedIfErrCheck) Run() (errors []error) {
	err := filepath.Walk("internal/services", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fileErrors := r.checkFile(path)
		errors = append(errors, fileErrors...)
		return nil
	})
	if err != nil {
		errors = append(errors, fmt.Errorf("walking internal/services: %+v", err))
	}

	return
}

func (r CombinedIfErrCheck) checkFile(path string) (errors []error) {
	f, err := os.Open(path)
	if err != nil {
		return []error{fmt.Errorf("opening %s: %+v", path, err)}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var prevLine string
	var prevLineNum int

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if prevLine != "" {
			if matchIfErrNil.MatchString(line) {
				errors = append(errors, fmt.Errorf(
					"%s:%d: '_, err' assignment should be combined with the following 'if err != nil' into a single 'if' init statement: %s",
					path, prevLineNum, strings.TrimSpace(prevLine),
				))
			}
			prevLine = ""
			prevLineNum = 0
		}

		// Check if this line is a "_, err :=" or "_, err =" that could be combined
		if matchUnderscoreErr.MatchString(line) && !matchMultiReturn.MatchString(line) && !matchForRange.MatchString(line) {
			prevLine = line
			prevLineNum = lineNum
		}
	}

	return
}

func (r CombinedIfErrCheck) Name() string {
	return "combinedIfErr"
}

func (r CombinedIfErrCheck) Description() string {
	return fmt.Sprintf(`
The '%s' check ensures that '_, err := SomeFunc()' followed by 'if err != nil' on
the next line is combined into a single 'if _, err := SomeFunc(); err != nil' statement.
This is the preferred Go style when the non-error return value is discarded.
`, r.Name())
}
