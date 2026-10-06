// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"strings"
	"testing"
)

func TestReporter(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
		passed   bool
	}{
		{
			name: "passing test is reported under its bare name as one block",
			input: []string{
				`{"Action":"start","Package":"example.com/a"}`,
				`{"Action":"run","Package":"example.com/a","Test":"TestAccThing_basic"}`,
				`{"Action":"output","Package":"example.com/a","Test":"TestAccThing_basic","Output":"=== RUN   TestAccThing_basic\n"}`,
				`{"Action":"output","Package":"example.com/a","Test":"TestAccThing_basic","Output":"--- PASS: TestAccThing_basic (1.50s)\n"}`,
				`{"Action":"pass","Package":"example.com/a","Test":"TestAccThing_basic","Elapsed":1.5}`,
				`{"Action":"output","Package":"example.com/a","Output":"PASS\n"}`,
				`{"Action":"pass","Package":"example.com/a","Elapsed":1.6}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccThing_basic' captureStandardOutput='false']`,
				`##teamcity[testStdOut name='TestAccThing_basic' out='=== RUN   TestAccThing_basic|n--- PASS: TestAccThing_basic (1.50s)']`,
				`##teamcity[testFinished name='TestAccThing_basic' duration='1500']`,
				`PASS`,
			},
			passed: true,
		},
		{
			name: "output of tests running in parallel is not interleaved",
			input: []string{
				`{"Action":"output","Package":"example.com/a","Test":"TestAccThing_one","Output":"one: first\n"}`,
				`{"Action":"output","Package":"example.com/a","Test":"TestAccThing_two","Output":"two: first\n"}`,
				`{"Action":"output","Package":"example.com/a","Test":"TestAccThing_one","Output":"one: second\n"}`,
				`{"Action":"pass","Package":"example.com/a","Test":"TestAccThing_two","Elapsed":0}`,
				`{"Action":"pass","Package":"example.com/a","Test":"TestAccThing_one","Elapsed":0}`,
				`{"Action":"pass","Package":"example.com/a","Elapsed":0}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccThing_two' captureStandardOutput='false']`,
				`##teamcity[testStdOut name='TestAccThing_two' out='two: first']`,
				`##teamcity[testFinished name='TestAccThing_two' duration='0']`,
				`##teamcity[testStarted name='TestAccThing_one' captureStandardOutput='false']`,
				`##teamcity[testStdOut name='TestAccThing_one' out='one: first|none: second']`,
				`##teamcity[testFinished name='TestAccThing_one' duration='0']`,
			},
			passed: true,
		},
		{
			name: "failing test",
			input: []string{
				`{"Action":"output","Package":"example.com/a","Test":"TestAccThing_fails","Output":"    a_test.go:17: it broke ['here'] | there\n"}`,
				`{"Action":"fail","Package":"example.com/a","Test":"TestAccThing_fails","Elapsed":0.25}`,
				`{"Action":"fail","Package":"example.com/a","Elapsed":0.3}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccThing_fails' captureStandardOutput='false']`,
				`##teamcity[testStdOut name='TestAccThing_fails' out='    a_test.go:17: it broke |[|'here|'|] || there']`,
				`##teamcity[testFailed name='TestAccThing_fails' message='Test failed']`,
				`##teamcity[testFinished name='TestAccThing_fails' duration='250']`,
			},
			passed: false,
		},
		{
			name: "skipped test",
			input: []string{
				`{"Action":"output","Package":"example.com/a","Test":"TestAccThing_skipped","Output":"    a_test.go:21: not today\n"}`,
				`{"Action":"skip","Package":"example.com/a","Test":"TestAccThing_skipped","Elapsed":0}`,
				`{"Action":"pass","Package":"example.com/a","Elapsed":0}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccThing_skipped' captureStandardOutput='false']`,
				`##teamcity[testStdOut name='TestAccThing_skipped' out='    a_test.go:21: not today']`,
				`##teamcity[testIgnored name='TestAccThing_skipped' message='Test skipped']`,
				`##teamcity[testFinished name='TestAccThing_skipped' duration='0']`,
			},
			passed: true,
		},
		{
			name: "subtests are folded into their top-level test",
			input: []string{
				`{"Action":"run","Package":"example.com/a","Test":"TestAccThing_sequence"}`,
				`{"Action":"run","Package":"example.com/a","Test":"TestAccThing_sequence/group/one"}`,
				`{"Action":"output","Package":"example.com/a","Test":"TestAccThing_sequence/group/one","Output":"    a_test.go:26: sub one\n"}`,
				`{"Action":"pass","Package":"example.com/a","Test":"TestAccThing_sequence/group/one","Elapsed":0}`,
				`{"Action":"output","Package":"example.com/a","Test":"TestAccThing_sequence/group/two","Output":"    a_test.go:27: sub two failed\n"}`,
				`{"Action":"fail","Package":"example.com/a","Test":"TestAccThing_sequence/group/two","Elapsed":0}`,
				`{"Action":"fail","Package":"example.com/a","Test":"TestAccThing_sequence/group","Elapsed":0}`,
				`{"Action":"fail","Package":"example.com/a","Test":"TestAccThing_sequence","Elapsed":2}`,
				`{"Action":"fail","Package":"example.com/a","Elapsed":2}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccThing_sequence' captureStandardOutput='false']`,
				`##teamcity[testStdOut name='TestAccThing_sequence' out='    a_test.go:26: sub one|n    a_test.go:27: sub two failed']`,
				`##teamcity[testFailed name='TestAccThing_sequence' message='Test failed']`,
				`##teamcity[testFinished name='TestAccThing_sequence' duration='2000']`,
			},
			passed: false,
		},
		{
			name: "tests cut short by a panic in another test are reported as failed",
			input: []string{
				`{"Action":"output","Package":"example.com/b","Test":"TestAccOther_slow","Output":"=== RUN   TestAccOther_slow\n"}`,
				`{"Action":"output","Package":"example.com/b","Test":"TestAccOther_panics","Output":"panic: assignment to entry in nil map\n"}`,
				`{"Action":"fail","Package":"example.com/b","Test":"TestAccOther_panics","Elapsed":0.02}`,
				`{"Action":"output","Package":"example.com/b","Output":"FAIL\texample.com/b\t0.576s\n"}`,
				`{"Action":"fail","Package":"example.com/b","Elapsed":0.576}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccOther_panics' captureStandardOutput='false']`,
				`##teamcity[testStdOut name='TestAccOther_panics' out='panic: assignment to entry in nil map']`,
				`##teamcity[testFailed name='TestAccOther_panics' message='Test failed']`,
				`##teamcity[testFinished name='TestAccOther_panics' duration='20']`,
				"FAIL\texample.com/b\t0.576s",
				`##teamcity[testStarted name='TestAccOther_slow' captureStandardOutput='false']`,
				`##teamcity[testStdOut name='TestAccOther_slow' out='=== RUN   TestAccOther_slow']`,
				`##teamcity[testFailed name='TestAccOther_slow' message='Test did not finish: the test binary exited before it reported a result (another test panicked, or the run timed out)']`,
				`##teamcity[testFinished name='TestAccOther_slow']`,
			},
			passed: false,
		},
		{
			name: "package which doesn't compile is a build problem",
			input: []string{
				`{"ImportPath":"example.com/c [example.com/c.test]","Action":"build-output","Output":"c/c_test.go:6:2: undefined: undefinedFunction\n"}`,
				`{"ImportPath":"example.com/c [example.com/c.test]","Action":"build-fail"}`,
				`{"Action":"start","Package":"example.com/c"}`,
				`{"Action":"fail","Package":"example.com/c","Elapsed":0,"FailedBuild":"example.com/c [example.com/c.test]"}`,
			},
			expected: []string{
				`c/c_test.go:6:2: undefined: undefinedFunction`,
				`##teamcity[buildProblem description='example.com/c failed to build']`,
			},
			passed: false,
		},
		{
			name: "package which fails outside of any test is a build problem",
			input: []string{
				`{"Action":"pass","Package":"example.com/a","Test":"TestAccThing_basic","Elapsed":0}`,
				`{"Action":"output","Package":"example.com/a","Output":"exit status 1\n"}`,
				`{"Action":"fail","Package":"example.com/a","Elapsed":0}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccThing_basic' captureStandardOutput='false']`,
				`##teamcity[testFinished name='TestAccThing_basic' duration='0']`,
				`exit status 1`,
				`##teamcity[buildProblem description='example.com/a failed without a failing test']`,
			},
			passed: false,
		},
		{
			name: "the same test name in two packages fails the build",
			input: []string{
				`{"Action":"run","Package":"example.com/a","Test":"TestAccThing_basic"}`,
				`{"Action":"pass","Package":"example.com/a","Test":"TestAccThing_basic","Elapsed":0}`,
				`{"Action":"pass","Package":"example.com/a","Elapsed":0}`,
				`{"Action":"run","Package":"example.com/d","Test":"TestAccThing_basic"}`,
				`{"Action":"pass","Package":"example.com/d","Test":"TestAccThing_basic","Elapsed":0}`,
				`{"Action":"pass","Package":"example.com/d","Elapsed":0}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccThing_basic' captureStandardOutput='false']`,
				`##teamcity[testFinished name='TestAccThing_basic' duration='0']`,
				`##teamcity[buildProblem description='TestAccThing_basic is defined in both example.com/a and example.com/d - test names must be unique']`,
				`##teamcity[testStarted name='TestAccThing_basic' captureStandardOutput='false']`,
				`##teamcity[testFailed name='TestAccThing_basic' message='Test name is already used by a test in example.com/a - test names must be unique']`,
				`##teamcity[testFinished name='TestAccThing_basic' duration='0']`,
			},
			passed: false,
		},
		{
			name: "a test run more than once in the same package is not a duplicate",
			input: []string{
				`{"Action":"run","Package":"example.com/a","Test":"TestAccThing_basic"}`,
				`{"Action":"pass","Package":"example.com/a","Test":"TestAccThing_basic","Elapsed":0}`,
				`{"Action":"run","Package":"example.com/a","Test":"TestAccThing_basic"}`,
				`{"Action":"pass","Package":"example.com/a","Test":"TestAccThing_basic","Elapsed":0}`,
				`{"Action":"pass","Package":"example.com/a","Elapsed":0}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccThing_basic' captureStandardOutput='false']`,
				`##teamcity[testFinished name='TestAccThing_basic' duration='0']`,
				`##teamcity[testStarted name='TestAccThing_basic' captureStandardOutput='false']`,
				`##teamcity[testFinished name='TestAccThing_basic' duration='0']`,
			},
			passed: true,
		},
		{
			name: "lines which aren't test events are passed through",
			input: []string{
				`go: downloading example.com/dependency v1.0.0`,
				`{"Action":"pass","Package":"example.com/a","Elapsed":0}`,
			},
			expected: []string{
				`go: downloading example.com/dependency v1.0.0`,
			},
			passed: true,
		},
		{
			name: "tests still running when the input ends are reported as failed",
			input: []string{
				`{"Action":"run","Package":"example.com/a","Test":"TestAccThing_basic"}`,
			},
			expected: []string{
				`##teamcity[testStarted name='TestAccThing_basic' captureStandardOutput='false']`,
				"##teamcity[testFailed name='TestAccThing_basic' message='Test did not finish: `go test` exited before it reported a result']",
				`##teamcity[testFinished name='TestAccThing_basic']`,
			},
			passed: false,
		},
		{
			name:  "no test results at all is a build problem",
			input: []string{},
			expected: []string{
				"##teamcity[buildProblem description='no test results were read from `go test`']",
			},
			passed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out strings.Builder
			passed := newReporter(&out).run(strings.NewReader(strings.Join(tc.input, "\n")))

			if passed != tc.passed {
				t.Errorf("expected passed to be %t but got %t", tc.passed, passed)
			}

			expected := ""
			if len(tc.expected) > 0 {
				expected = strings.Join(tc.expected, "\n") + "\n"
			}
			if actual := out.String(); actual != expected {
				t.Errorf("expected:\n%s\nbut got:\n%s", expected, actual)
			}
		})
	}
}
