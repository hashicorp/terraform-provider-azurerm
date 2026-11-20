// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"reflect"
	"testing"
)

func TestSentences(t *testing.T) {
	cases := []struct {
		name string
		text string
		want []sentence
	}{
		{"one sentence", "The name of the Example.", []sentence{{"The name of the Example.", 0}}},
		{"two sentences", "The name. Changing this forces a new resource.", []sentence{{"The name.", 0}, {"Changing this forces a new resource.", 10}}},
		{"question", "Should this be enabled? Defaults to `true`.", []sentence{{"Should this be enabled?", 0}, {"Defaults to `true`.", 24}}},
		{"no full stop", "Manages an Example", []sentence{{"Manages an Example", 0}}},
		{"full stop inside a value", "The version, such as `1.2`.", []sentence{{"The version, such as `1.2`.", 0}}},
		{"empty", "", []sentence{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := sentences(c.text); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("sentences(%q) = %+v, want %+v", c.text, got, c.want)
			}
		})
	}
}

func TestSentenceAt(t *testing.T) {
	lines := []string{
		"* `ttl` - (Optional) A optional time to live value. If omitted, the suppression will not expire.",
		"Manages an Certificate within an API Management Service.",
		"* `id` - The ID of the the Example.",
	}

	cases := []struct {
		name       string
		line       int
		column     int
		want       string
		wantOffset int
		wantFound  bool
	}{
		{"start of a property's description", 1, 22, "A optional time to live value.", 0, true},
		{"second sentence of a property", 1, 56, "If omitted, the suppression will not expire.", 3, true},
		{"line with no property", 2, 9, "Manages an Certificate within an API Management Service.", 8, true},
		{"property with no Required or Optional", 3, 21, "The ID of the the Example.", 11, true},
		{"inside the property name", 1, 4, "", 0, false},
		{"line that does not exist", 9, 1, "", 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, offset, found := sentenceAt(lines, c.line, c.column)
			if got != c.want || offset != c.wantOffset || found != c.wantFound {
				t.Fatalf("sentenceAt(line %d, column %d) = %q, %d, %t, want %q, %d, %t", c.line, c.column, got, offset, found, c.want, c.wantOffset, c.wantFound)
			}
		})
	}
}
