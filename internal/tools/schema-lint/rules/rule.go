// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package rules defines the pluggable rule contracts for the schema linter and
// hosts the concrete rules that ship with it.
//
// A rule is a small, independent check (in the spirit of a markdown linter such
// as markdownlint) that inspects the provider schema and reports Findings. Rules
// come in two flavours:
//
//   - PropertyRule evaluates a single schema property node (including nested
//     block attributes), and
//   - ResourceRule evaluates an entire resource or data source schema.
//
// New rules are plugged in by implementing one of these interfaces and appending
// the rule to AllRules.
package rules

import (
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

// Severity indicates how a Finding should be treated by the reporter and the
// process exit code.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Kind identifies whether a schema belongs to a resource or a data source.
type Kind string

const (
	KindResource   Kind = "resource"
	KindDataSource Kind = "data source"
)

// Finding represents a single rule violation.
type Finding struct {
	RuleID       string   `json:"ruleId"`
	RuleName     string   `json:"ruleName"`
	Severity     Severity `json:"severity"`
	ResourceType string   `json:"resourceType"`
	Kind         Kind     `json:"kind"`
	// Path is the dotted property path from the resource root (empty for
	// resource-level findings).
	Path    string `json:"path,omitempty"`
	Message string `json:"message"`
	// FixSuggestion is a concrete suggested remediation. It is populated only
	// when fix suggestions are requested and the rule that produced the finding
	// implements Fixer.
	FixSuggestion string `json:"fixSuggestion,omitempty"`
}

// PropertyContext carries the information a PropertyRule needs to evaluate a
// single schema node.
type PropertyContext struct {
	ResourceType string
	Kind         Kind
	// Name is the leaf property name.
	Name string
	// Path is the full dotted path from the resource root (e.g. "network_profile.load_balancer_sku").
	Path   string
	Schema providerschema.SchemaJSON
	// Parent is the schema of the enclosing block, or nil at the top level.
	Parent *providerschema.SchemaJSON
	// Siblings is the set of properties at the same level as this one (including
	// this property), keyed by name. It is nil when unavailable.
	Siblings map[string]providerschema.SchemaJSON
}

// ResourceContext carries the information a ResourceRule needs to evaluate a
// whole resource or data source schema.
type ResourceContext struct {
	ResourceType string
	Kind         Kind
	Resource     providerschema.ResourceJSON
}

// Rule is the common interface implemented by every lint rule.
type Rule interface {
	// ID is the stable identifier for the rule (e.g. "SL001").
	ID() string
	// Name is a short human readable slug (e.g. "property-description-required").
	Name() string
	// Description explains what the rule checks.
	Description() string
	// DefaultSeverity is the severity used when the configuration does not
	// override it.
	DefaultSeverity() Severity
}

// PropertyRule is a Rule that evaluates a single schema property node.
type PropertyRule interface {
	Rule
	CheckProperty(ctx PropertyContext) []Finding
}

// ResourceRule is a Rule that evaluates an entire resource or data source schema.
type ResourceRule interface {
	Rule
	CheckResource(ctx ResourceContext) []Finding
}

// Fixer is an optional interface a Rule may implement to indicate it can suggest
// a remediation for its findings.
//
// Because the linter inspects the compiled provider schema rather than its Go
// source, a fix is surfaced as an actionable suggestion on Finding.FixSuggestion
// (computed when the finding is created) rather than an in-place source edit. The
// linter only surfaces those suggestions when they are requested.
type Fixer interface {
	Rule
	// FixHint is a short description of the remediation the rule suggests, shown
	// by the `list` command.
	FixHint() string
}

// Fixable reports whether a rule can suggest fixes for its findings.
func Fixable(r Rule) bool {
	_, ok := r.(Fixer)
	return ok
}

// AllRules is the registry of every rule known to the linter. Plug a new rule in
// by appending it here.
var AllRules = []Rule{
	propertyDescriptionRequired{},
	optionalAndRequired{},
	requiredAndComputed{},
	computedOnlyForceNew{},
	collectionElemRequired{},
	singlePropertyBlock{},
	limitsOnNonCollection{},
	avoidNoneValue{},
	validationRequired{},
	blockNeedsConstraint{},
	arrayLimits{},
	skuFieldNaming{},
	unitInNaming{},
	noAbbreviations{},
	redundantIsPrefix{},
	redundantSuffix{},
	idReferenceValidation{},
}

// ByID returns the registered rule with the given ID (case-insensitive) and
// whether it was found.
func ByID(id string) (Rule, bool) {
	for _, r := range AllRules {
		if strings.EqualFold(r.ID(), id) {
			return r, true
		}
	}
	return nil, false
}

// PropertyRules returns the subset of rules that evaluate property nodes.
func PropertyRules(in []Rule) []PropertyRule {
	out := make([]PropertyRule, 0, len(in))
	for _, r := range in {
		if pr, ok := r.(PropertyRule); ok {
			out = append(out, pr)
		}
	}
	return out
}

// ResourceRules returns the subset of rules that evaluate whole resources.
func ResourceRules(in []Rule) []ResourceRule {
	out := make([]ResourceRule, 0, len(in))
	for _, r := range in {
		if rr, ok := r.(ResourceRule); ok {
			out = append(out, rr)
		}
	}
	return out
}

// propertyFinding builds a Finding for a property-level rule using the rule's
// own ID, name and default severity.
func propertyFinding(r Rule, ctx PropertyContext, message string) Finding {
	return Finding{
		RuleID:       r.ID(),
		RuleName:     r.Name(),
		Severity:     r.DefaultSeverity(),
		ResourceType: ctx.ResourceType,
		Kind:         ctx.Kind,
		Path:         ctx.Path,
		Message:      message,
	}
}

// propertyFindingFix builds a property-level Finding that also carries a
// suggested remediation.
func propertyFindingFix(r Rule, ctx PropertyContext, message, suggestion string) Finding {
	f := propertyFinding(r, ctx, message)
	f.FixSuggestion = suggestion
	return f
}
