// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package engine

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/config"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/rules"
)

func testSchema() *providerschema.ProviderSchemaJSON {
	return &providerschema.ProviderSchemaJSON{
		ResourcesMap: map[string]providerschema.ResourceJSON{
			"azurerm_test": {
				Schema: map[string]providerschema.SchemaJSON{
					// Missing description -> SL001.
					"name": {Type: rules.TypeString, Required: true},
					// Nested block, missing description -> SL001.
					"block": {
						Type: rules.TypeList,
						Elem: &providerschema.ResourceJSON{
							Schema: map[string]providerschema.SchemaJSON{
								// Missing description -> SL001 at the dotted path block.inner.
								"inner": {Type: rules.TypeString, Optional: true},
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

	// The nested block property is reached with a dotted path.
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

func fixSchema() *providerschema.ProviderSchemaJSON {
	return &providerschema.ProviderSchemaJSON{
		ResourcesMap: map[string]providerschema.ResourceJSON{
			"azurerm_fix": {
				Timeouts: &providerschema.ResourceTimeoutJSON{Read: 5},
				Schema: map[string]providerschema.SchemaJSON{
					// A MaxItems:1 block with a single nested "enabled" bool ->
					// SL002 (single-property-block, fixable). Descriptions are set
					// so SL001 does not also fire.
					"sign_in": {
						Type:        rules.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "sign in settings",
						Elem: &providerschema.ResourceJSON{
							Schema: map[string]providerschema.SchemaJSON{
								"enabled": {Type: rules.TypeBool, Optional: true, Description: "whether sign in is enabled"},
							},
						},
					},
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
	f := hasFinding(l.Lint(fixSchema()), "SL002", "sign_in")
	if f == nil {
		t.Fatal("expected an SL002 finding")
	}
	if f.FixSuggestion == "" {
		t.Fatalf("expected a fix suggestion for SL002 when -fix is set, got none")
	}
}

func anyFinding(findings []rules.Finding, resourceType, path string) bool {
	for i := range findings {
		if findings[i].ResourceType == resourceType && findings[i].Path == path {
			return true
		}
	}
	return false
}

func TestLint_DiffMode(t *testing.T) {
	base := &providerschema.ProviderSchemaJSON{
		ResourcesMap: map[string]providerschema.ResourceJSON{
			"azurerm_test": {
				Schema: map[string]providerschema.SchemaJSON{
					"existing": {Type: rules.TypeString, Optional: true},
					"block": {
						Type: rules.TypeList,
						Elem: &providerschema.ResourceJSON{
							Schema: map[string]providerschema.SchemaJSON{
								"old": {Type: rules.TypeString, Optional: true},
							},
						},
					},
				},
			},
		},
	}

	current := &providerschema.ProviderSchemaJSON{
		ResourcesMap: map[string]providerschema.ResourceJSON{
			// Existing resource that gains a top-level property and a nested one.
			"azurerm_test": {
				Schema: map[string]providerschema.SchemaJSON{
					"existing":  {Type: rules.TypeString, Optional: true},
					"brand_new": {Type: rules.TypeString, Optional: true},
					"block": {
						Type: rules.TypeList,
						Elem: &providerschema.ResourceJSON{
							Schema: map[string]providerschema.SchemaJSON{
								"old":   {Type: rules.TypeString, Optional: true},
								"added": {Type: rules.TypeString, Optional: true},
							},
						},
					},
				},
			},
			// A brand new resource -> all of its properties are reported.
			"azurerm_new": {
				Schema: map[string]providerschema.SchemaJSON{
					"np": {Type: rules.TypeString, Optional: true},
				},
			},
		},
	}

	l := New(nil, Options{Config: &config.Config{}, BaseSchema: base})
	findings := l.Lint(current)

	// Pre-existing properties (and the pre-existing block itself) are not reported.
	for _, f := range findings {
		if f.ResourceType == "azurerm_test" && (f.Path == "existing" || f.Path == "block" || f.Path == "block.old") {
			t.Fatalf("did not expect a finding on pre-existing property, got %+v", f)
		}
	}

	// Newly-added top-level and nested properties are reported.
	if !anyFinding(findings, "azurerm_test", "brand_new") {
		t.Fatalf("expected a finding on the newly-added property brand_new, got %+v", findings)
	}
	if !anyFinding(findings, "azurerm_test", "block.added") {
		t.Fatalf("expected a finding on the newly-added nested property block.added, got %+v", findings)
	}

	// Every property of an entirely new resource is reported.
	if !anyFinding(findings, "azurerm_new", "np") {
		t.Fatalf("expected findings on the new resource azurerm_new, got %+v", findings)
	}
}
