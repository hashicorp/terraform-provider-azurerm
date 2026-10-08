// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// teamcity-test-reporter reads the output of `go test -json` on stdin and reports each test to
// TeamCity using Service Messages:
//
//	go test -v ./internal/services/foo/... -json | go run ./internal/tools/teamcity-test-reporter
//
// Tests are reported under their bare name (`TestAccFoo_basic`, with no package) so that a test
// keeps its history in TeamCity when it moves between packages. The output of each test is held
// back until that test ends, so that it shows in the build log as one collapsible block rather than
// interleaved with the tests running alongside it. Subtests are folded into their top-level test.
//
// Because the package is dropped, a test name used in two packages fails the build.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func main() {
	r := newReporter(os.Stdout)
	if !r.run(os.Stdin) {
		os.Exit(1)
	}
}

// event is a line of `go test -json` output, see `go doc test2json`.
type event struct {
	Action      string  `json:"Action"`
	Package     string  `json:"Package"`
	Test        string  `json:"Test"`
	Elapsed     float64 `json:"Elapsed"`
	Output      string  `json:"Output"`
	FailedBuild string  `json:"FailedBuild"`
}

type testKey struct {
	pkg  string
	name string
}

type reporter struct {
	out io.Writer

	// the output so far of each top-level test which has started but not yet ended
	running map[testKey]*strings.Builder
	ended   map[testKey]bool

	// tests are reported by name alone, so a name must only ever be used by one package: these are
	// the package which first used each name, and the tests found reusing another package's name
	packageForName map[string]string
	duplicateOf    map[testKey]string
	failedTests    map[string]int

	events int
	failed bool
}

func newReporter(out io.Writer) *reporter {
	return &reporter{
		out:            out,
		running:        make(map[testKey]*strings.Builder),
		ended:          make(map[testKey]bool),
		packageForName: make(map[string]string),
		duplicateOf:    make(map[testKey]string),
		failedTests:    make(map[string]int),
	}
}

// run reports every test read from in, returning whether the test run passed.
func (r *reporter) run(in io.Reader) bool {
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadBytes('\n')
		if len(bytes.TrimSpace(line)) > 0 {
			r.handleLine(line)
		}
		if err != nil {
			break
		}
	}

	// `go test` was cut short, so no package result arrived to close these off
	for _, key := range r.unfinished("") {
		r.report(key, "fail", -1, "Test did not finish: `go test` exited before it reported a result")
	}

	if r.events == 0 {
		r.buildProblem("no test results were read from `go test`")
	}

	return !r.failed
}

func (r *reporter) handleLine(line []byte) {
	var e event
	if err := json.Unmarshal(line, &e); err != nil || e.Action == "" {
		// not a test event (e.g. a message from the go toolchain) so pass it through untouched
		r.print(string(line))
		return
	}
	r.events++

	switch {
	case e.Action == "build-output":
		r.print(e.Output)
	case e.Action == "build-fail":
		// reported when the package fails, below
	case e.Test == "":
		r.handlePackage(e)
	default:
		r.handleTest(e)
	}
}

func (r *reporter) handlePackage(e event) {
	switch e.Action {
	case "output":
		r.print(e.Output)

	case "pass", "skip", "fail":
		// the test binary has exited, so anything still running was cut short by a panic or timeout
		for _, key := range r.unfinished(e.Package) {
			r.report(key, "fail", -1, "Test did not finish: the test binary exited before it reported a result (another test panicked, or the run timed out)")
		}

		if e.Action != "fail" {
			return
		}

		switch {
		case e.FailedBuild != "":
			r.buildProblem(fmt.Sprintf("%s failed to build", e.Package))
		case r.failedTests[e.Package] == 0:
			r.buildProblem(fmt.Sprintf("%s failed without a failing test", e.Package))
		}
	}
}

func (r *reporter) handleTest(e event) {
	// subtests (`TestFoo/bar`) are folded into their top-level test
	name, _, isSubtest := strings.Cut(e.Test, "/")
	key := testKey{pkg: e.Package, name: name}

	if r.ended[key] {
		if e.Action == "run" && !isSubtest {
			// the test is being run again (`-count`)
			delete(r.ended, key)
		} else {
			// output arriving after the test's result has nothing to be attached to
			if e.Action == "output" {
				r.print(e.Output)
			}
			return
		}
	}

	if _, ok := r.running[key]; !ok {
		r.running[key] = &strings.Builder{}
		r.checkNameIsUnique(key)
	}

	switch e.Action {
	case "output":
		r.running[key].WriteString(e.Output)

	case "pass", "fail", "skip":
		if !isSubtest {
			r.report(key, e.Action, e.Elapsed, "Test failed")
		}
	}
}

// checkNameIsUnique fails the build when a test uses a name already used by a test in another
// package, since TeamCity would otherwise merge the two into one test. The run is left to finish
// rather than stopped, so that the tests already running still clean up after themselves.
func (r *reporter) checkNameIsUnique(key testKey) {
	pkg, ok := r.packageForName[key.name]
	if !ok {
		r.packageForName[key.name] = key.pkg
		return
	}
	if pkg == key.pkg {
		return
	}

	r.duplicateOf[key] = pkg
	r.buildProblem(fmt.Sprintf("%s is defined in both %s and %s - test names must be unique", key.name, pkg, key.pkg))
}

// unfinished returns the tests which have started but not ended, for one package or ("") for all.
func (r *reporter) unfinished(pkg string) []testKey {
	keys := make([]testKey, 0)
	for key := range r.running {
		if pkg == "" || key.pkg == pkg {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].pkg != keys[j].pkg {
			return keys[i].pkg < keys[j].pkg
		}
		return keys[i].name < keys[j].name
	})
	return keys
}

// report writes a finished test and its output as a single block of Service Messages. A negative
// elapsed means the duration isn't known.
func (r *reporter) report(key testKey, result string, elapsed float64, failureMessage string) {
	output := ""
	if running, ok := r.running[key]; ok {
		output = strings.TrimRight(running.String(), "\n")
	}
	delete(r.running, key)
	r.ended[key] = true

	if pkg, ok := r.duplicateOf[key]; ok {
		result = "fail"
		failureMessage = fmt.Sprintf("Test name is already used by a test in %s - test names must be unique", pkg)
	}

	name := escape(key.name)
	fmt.Fprintf(r.out, "##teamcity[testStarted name='%s' captureStandardOutput='false']\n", name)
	if output != "" {
		fmt.Fprintf(r.out, "##teamcity[testStdOut name='%s' out='%s']\n", name, escape(output))
	}

	switch result {
	case "fail":
		r.failed = true
		r.failedTests[key.pkg]++
		fmt.Fprintf(r.out, "##teamcity[testFailed name='%s' message='%s']\n", name, escape(failureMessage))
	case "skip":
		fmt.Fprintf(r.out, "##teamcity[testIgnored name='%s' message='Test skipped']\n", name)
	}

	if elapsed < 0 {
		fmt.Fprintf(r.out, "##teamcity[testFinished name='%s']\n", name)
		return
	}
	fmt.Fprintf(r.out, "##teamcity[testFinished name='%s' duration='%d']\n", name, int64(elapsed*1000))
}

// buildProblem fails the build for a reason which isn't a failing test.
func (r *reporter) buildProblem(description string) {
	r.failed = true
	fmt.Fprintf(r.out, "##teamcity[buildProblem description='%s']\n", escape(description))
}

func (r *reporter) print(text string) {
	fmt.Fprintln(r.out, strings.TrimRight(text, "\n"))
}

var escaper = strings.NewReplacer(
	"|", "||",
	"'", "|'",
	"\n", "|n",
	"\r", "|r",
	"[", "|[",
	"]", "|]",
)

// escape makes a value safe to use within a TeamCity Service Message.
func escape(value string) string {
	return escaper.Replace(value)
}
