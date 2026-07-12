// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package engine loads the provider schema and applies the configured lint
// rules to every resource, data source and nested block property.
package engine

import (
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/config"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/rules"
)

// Options controls a lint run. CLI flags populate these and take precedence over
// the config file.
type Options struct {
	// Config is the loaded configuration file (never nil once passed to New).
	Config *config.Config
	// OnlyRules, when non-empty, restricts the run to these rule IDs (from the
	// CLI, excluding the "all" sentinel).
	OnlyRules []string
	// DisableRules lists rule IDs to disable; this takes precedence over
	// everything else.
	DisableRules []string
	// SuggestFixes, when true, keeps the fix suggestions that fixable rules
	// attach to their findings. When false the suggestions are discarded.
	SuggestFixes bool
}

// Linter applies a set of rules to the provider schema.
type Linter struct {
	rules   []rules.Rule
	options Options
}

// New returns a Linter over the given rule set (defaults to rules.AllRules when
// nil).
func New(ruleSet []rules.Rule, opts Options) *Linter {
	if ruleSet == nil {
		ruleSet = rules.AllRules
	}
	if opts.Config == nil {
		opts.Config = &config.Config{}
	}

	return &Linter{rules: ruleSet, options: opts}
}

// Run loads the live provider schema and lints it.
func (l *Linter) Run() ([]rules.Finding, error) {
	schema, err := providerjson.ProviderFromRaw(providerjson.LoadData())
	if err != nil {
		return nil, err
	}

	return l.Lint(schema), nil
}

// Lint applies the active rules to the supplied provider schema. It is separated
// from Run so it can be exercised with a synthetic schema in tests.
func (l *Linter) Lint(schema *providerjson.ProviderSchemaJSON) []rules.Finding {
	propertyRules, resourceRules, severity := l.activeRules()

	findings := make([]rules.Finding, 0)

	lintOne := func(name string, res providerjson.ResourceJSON, kind rules.Kind) {
		if !l.options.Config.ResourceAllowed(name) {
			return
		}

		resourceCtx := rules.ResourceContext{ResourceType: name, Kind: kind, Resource: res}
		for _, rr := range resourceRules {
			findings = append(findings, rr.CheckResource(resourceCtx)...)
		}

		for _, propName := range sortedKeys(res.Schema) {
			walkProperty(name, kind, propName, propName, res.Schema[propName], nil, res.Schema, propertyRules, &findings)
		}
	}

	if schema != nil {
		for _, name := range sortedResourceKeys(schema.ResourcesMap) {
			lintOne(name, schema.ResourcesMap[name], rules.KindResource)
		}
		for _, name := range sortedResourceKeys(schema.DataSourcesMap) {
			lintOne(name, schema.DataSourcesMap[name], rules.KindDataSource)
		}
	}

	// Apply configured severity overrides.
	for i := range findings {
		if s, ok := severity[strings.ToUpper(findings[i].RuleID)]; ok {
			findings[i].Severity = s
		}
	}

	// Fix suggestions are computed by rules during checking; drop them unless
	// they were explicitly requested.
	if !l.options.SuggestFixes {
		for i := range findings {
			findings[i].FixSuggestion = ""
		}
	}

	sortFindings(findings)

	return findings
}

// activeRules resolves which rules run and their effective severity by merging
// rule defaults, the config file and the CLI overrides (CLI wins).
func (l *Linter) activeRules() (propertyRules []rules.PropertyRule, resourceRules []rules.ResourceRule, severity map[string]rules.Severity) {
	severity = make(map[string]rules.Severity)
	only := toSet(l.options.OnlyRules)
	disabled := toSet(l.options.DisableRules)
	cfg := l.options.Config

	for _, r := range l.rules {
		id := strings.ToUpper(r.ID())

		// CLI -disable takes precedence over everything.
		if disabled[id] {
			continue
		}
		// CLI -rules restricts to an explicit set.
		if len(only) > 0 && !only[id] {
			continue
		}

		enabled := cfg.DefaultEnabled()
		sev := r.DefaultSeverity()
		if rc, ok := cfg.Rule(id); ok {
			if rc.Enabled != nil {
				enabled = *rc.Enabled
			}
			if rc.Severity != "" {
				sev = rules.Severity(strings.ToLower(rc.Severity))
			}
		}

		if !enabled {
			continue
		}

		severity[id] = sev
		if pr, ok := r.(rules.PropertyRule); ok {
			propertyRules = append(propertyRules, pr)
		}
		if rr, ok := r.(rules.ResourceRule); ok {
			resourceRules = append(resourceRules, rr)
		}
	}

	return propertyRules, resourceRules, severity
}

// walkProperty applies the property rules to a node and recurses into nested
// blocks, building a dotted path as it descends.
func walkProperty(resourceType string, kind rules.Kind, name, path string, s providerjson.SchemaJSON, parent *providerjson.SchemaJSON, siblings map[string]providerjson.SchemaJSON, propertyRules []rules.PropertyRule, out *[]rules.Finding) {
	ctx := rules.PropertyContext{
		ResourceType: resourceType,
		Kind:         kind,
		Name:         name,
		Path:         path,
		Schema:       s,
		Parent:       parent,
		Siblings:     siblings,
	}

	for _, pr := range propertyRules {
		*out = append(*out, pr.CheckProperty(ctx)...)
	}

	if block, ok := rules.BlockElem(s); ok {
		for _, childName := range sortedKeys(block.Schema) {
			walkProperty(resourceType, kind, childName, path+"."+childName, block.Schema[childName], &s, block.Schema, propertyRules, out)
		}
	}
}

func sortedKeys(m map[string]providerjson.SchemaJSON) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func sortedResourceKeys(m map[string]providerjson.ResourceJSON) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func sortFindings(findings []rules.Finding) {
	sort.SliceStable(findings, func(i, j int) bool {
		a, b := findings[i], findings[j]
		if a.ResourceType != b.ResourceType {
			return a.ResourceType < b.ResourceType
		}
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.RuleID < b.RuleID
	})
}

func toSet(in []string) map[string]bool {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]bool, len(in))
	for _, v := range in {
		if v = strings.TrimSpace(v); v != "" {
			out[strings.ToUpper(v)] = true
		}
	}

	return out
}
