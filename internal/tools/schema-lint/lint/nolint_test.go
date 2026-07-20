// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lint

import "testing"

// TestNolintBaseline confirms the property used by the suppression tests reports
// the findings we then expect //nolint to silence.
func TestNolintBaseline(t *testing.T) {
	got := rulesAtPath(lintSrc(t, wrap(`"s": {Type: pluginsdk.TypeString, Optional: true},`)), "s")
	if !got["SL001"] || !got["SL005"] {
		t.Fatalf("expected SL001 and SL005 on s, got %v", got)
	}
}

func TestNolintSuppresses(t *testing.T) {
	cases := []struct {
		name string
		body string
		// suppressed maps a rule to whether it should be silenced for "s".
		suppressed map[string]bool
	}{
		{
			name:       "bare trailing suppresses all",
			body:       `"s": {Type: pluginsdk.TypeString, Optional: true}, //nolint`,
			suppressed: map[string]bool{"SL001": true, "SL005": true},
		},
		{
			name:       "bare trailing with a space suppresses all",
			body:       `"s": {Type: pluginsdk.TypeString, Optional: true}, // nolint`,
			suppressed: map[string]bool{"SL001": true, "SL005": true},
		},
		{
			name:       "specific trailing suppresses one",
			body:       `"s": {Type: pluginsdk.TypeString, Optional: true}, //nolint:SL001`,
			suppressed: map[string]bool{"SL001": true, "SL005": false},
		},
		{
			name:       "specific trailing with a space suppresses one",
			body:       `"s": {Type: pluginsdk.TypeString, Optional: true}, // nolint:SL001`,
			suppressed: map[string]bool{"SL001": true, "SL005": false},
		},
		{
			name:       "comma separated rules",
			body:       `"s": {Type: pluginsdk.TypeString, Optional: true}, //nolint:SL001,SL005`,
			suppressed: map[string]bool{"SL001": true, "SL005": true},
		},
		{
			name:       "trailing explanation is ignored",
			body:       `"s": {Type: pluginsdk.TypeString, Optional: true}, //nolint:SL001 // intentional`,
			suppressed: map[string]bool{"SL001": true, "SL005": false},
		},
		{
			name:       "rule id is case insensitive",
			body:       `"s": {Type: pluginsdk.TypeString, Optional: true}, //nolint:sl001`,
			suppressed: map[string]bool{"SL001": true, "SL005": false},
		},
		{
			name:       "leading comment on its own line",
			body:       "//nolint:SL001\n\t\t\t\"s\": {Type: pluginsdk.TypeString, Optional: true},",
			suppressed: map[string]bool{"SL001": true, "SL005": false},
		},
		{
			name:       "leading comment on its own line with a space",
			body:       "// nolint:SL001\n\t\t\t\"s\": {Type: pluginsdk.TypeString, Optional: true},",
			suppressed: map[string]bool{"SL001": true, "SL005": false},
		},
		{
			name:       "trailing on a multi-line property key line",
			body:       "\"s\": { //nolint:SL001\n\t\t\t\tType:     pluginsdk.TypeString,\n\t\t\t\tOptional: true,\n\t\t\t},",
			suppressed: map[string]bool{"SL001": true, "SL005": false},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rulesAtPath(lintSrc(t, wrap(tc.body)), "s")
			for rule, want := range tc.suppressed {
				if want && got[rule] {
					t.Errorf("expected %s suppressed, still reported", rule)
				}
				if !want && !got[rule] {
					t.Errorf("expected %s reported, was suppressed", rule)
				}
			}
		})
	}
}

// TestParseNolintForms documents that the directive is accepted both as the
// gofmt-preserving machine form (//nolint) and with a leading space (// nolint),
// with or without a rule list.
func TestParseNolintForms(t *testing.T) {
	all := []struct {
		text string
	}{
		{"//nolint"},
		{"// nolint"},
		{"//  nolint"},
	}
	for _, tc := range all {
		d, ok := parseNolint(tc.text)
		if !ok || !d.all {
			t.Errorf("parseNolint(%q) = (%+v, %v); want all-rules directive", tc.text, d, ok)
		}
	}

	specific := []string{"//nolint:SL001", "// nolint:SL001", "// nolint:SL001 // why"}
	for _, text := range specific {
		d, ok := parseNolint(text)
		if !ok || d.all || !d.rules["SL001"] {
			t.Errorf("parseNolint(%q) = (%+v, %v); want SL001 directive", text, d, ok)
		}
	}

	if _, ok := parseNolint("// not a directive"); ok {
		t.Error("parseNolint should reject a non-directive comment")
	}
	if _, ok := parseNolint("//nolintfoo"); ok {
		t.Error("parseNolint should reject //nolintfoo")
	}
}

// TestNolintScopedToProperty confirms a directive only affects its own property,
// not a sibling on an adjacent line.
func TestNolintScopedToProperty(t *testing.T) {
	body := "\t\t\t\"a\": {Type: pluginsdk.TypeString, Optional: true}, //nolint\n" +
		"\t\t\t\"b\": {Type: pluginsdk.TypeString, Optional: true},"
	f := lintSrc(t, wrap(body))
	if got := rulesAtPath(f, "a"); len(got) != 0 {
		t.Errorf("expected a fully suppressed, got %v", got)
	}
	if !rulesAtPath(f, "b")["SL001"] {
		t.Errorf("expected b still reported, got %v", rulesAtPath(f, "b"))
	}
}
