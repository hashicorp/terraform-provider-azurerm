// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package engine

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/config"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/rules"
)

func testSchema() *providerjson.ProviderSchemaJSON {
	return &providerjson.ProviderSchemaJSON{
		ResourcesMap: map[string]providerjson.ResourceJSON{
			"azurerm_test": {
				Schema: map[string]providerjson.SchemaJSON{
					// Missing description -> SL001.
					"name": {Type: rules.TypeString, Required: true},
					// Nested block, missing description -> SL001.
					"block": {
						Type: rules.TypeList,
						Elem: &providerjson.ResourceJSON{
							Schema: map[string]providerjson.SchemaJSON{
								// Optional + Required -> SL002, and missing description -> SL001.
								"inner": {Type: rules.TypeString, Optional: true, Required: true},
							},
						},
					},
				},
			},
		},
	}
}

func hasFinding(findings []rules.Finding, ruleID, path string) *rules.Finding {
	for i := range findings {
		if findings[i].RuleID == ruleID && findings[i].Path == path {
			return &findings[i]
		}
	}
	return nil
}

func TestLint_NestedTraversal(t *testing.T) {
	l := New(nil, Options{Config: &config.Config{}})
	findings := l.Lint(testSchema())

	// Nested block property is reached with a dotted path and SL002 fires.
	if hasFinding(findings, "SL002", "block.inner") == nil {
		t.Fatalf("expected SL002 finding at path block.inner, got %+v", findings)
	}
	// The nested property is also checked for a missing description.
	if hasFinding(findings, "SL001", "block.inner") == nil {
		t.Fatalf("expected SL001 finding at path block.inner, got %+v", findings)
	}
}

func TestLint_DisableRule(t *testing.T) {
	l := New(nil, Options{Config: &config.Config{}, DisableRules: []string{"SL001"}})
	findings := l.Lint(testSchema())

	for _, f := range findings {
		if f.RuleID == "SL001" {
			t.Fatalf("SL001 should have been disabled, got %+v", f)
		}
	}
}

func TestLint_SeverityOverrideFromConfig(t *testing.T) {
	cfg := &config.Config{
		Rules: map[string]config.RuleConfig{
			"SL001": {Severity: "error"},
		},
	}
	l := New(nil, Options{Config: cfg})
	findings := l.Lint(testSchema())

	f := hasFinding(findings, "SL001", "name")
	if f == nil {
		t.Fatalf("expected SL001 finding, got %+v", findings)
	}
	if f.Severity != rules.SeverityError {
		t.Fatalf("expected SL001 severity to be overridden to error, got %q", f.Severity)
	}
}

func TestLint_OnlyRules(t *testing.T) {
	l := New(nil, Options{Config: &config.Config{}, OnlyRules: []string{"SL001"}})
	findings := l.Lint(testSchema())

	if len(findings) == 0 {
		t.Fatal("expected findings, got none")
	}
	for _, f := range findings {
		if f.RuleID != "SL001" {
			t.Fatalf("expected only SL001 findings, got %+v", f)
		}
	}
}

func TestLint_ResourceFilterSkip(t *testing.T) {
	cfg := &config.Config{SkipResources: []string{"azurerm_test"}}
	l := New(nil, Options{Config: cfg})
	findings := l.Lint(testSchema())

	if len(findings) != 0 {
		t.Fatalf("expected no findings when resource is skipped, got %+v", findings)
	}
}

func fixSchema() *providerjson.ProviderSchemaJSON {
	return &providerjson.ProviderSchemaJSON{
		ResourcesMap: map[string]providerjson.ResourceJSON{
			"azurerm_fix": {
				Timeouts: &providerjson.ResourceTimeoutJSON{Read: 5},
				Schema: map[string]providerjson.SchemaJSON{
					// Computed-only + ForceNew -> SL004 (fixable). Has a description
					// so SL001 does not also fire.
					"id_out": {Type: rules.TypeString, Computed: true, ForceNew: true, Description: "an output"},
				},
			},
		},
	}
}

func TestLint_FixSuggestionsDisabledByDefault(t *testing.T) {
	l := New(nil, Options{Config: &config.Config{}})
	for _, f := range l.Lint(fixSchema()) {
		if f.FixSuggestion != "" {
			t.Fatalf("did not expect a fix suggestion without -fix, got %+v", f)
		}
	}
}

func TestLint_FixSuggestionsWhenRequested(t *testing.T) {
	l := New(nil, Options{Config: &config.Config{}, SuggestFixes: true})
	f := hasFinding(l.Lint(fixSchema()), "SL004", "id_out")
	if f == nil {
		t.Fatal("expected an SL004 finding")
	}
	if f.FixSuggestion == "" {
		t.Fatalf("expected a fix suggestion for SL004 when -fix is set, got none")
	}
}
