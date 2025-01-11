package budgets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedBudgetId{})
}

var _ resourceids.ResourceId = &ScopedBudgetId{}

// ScopedBudgetId is a struct representing the Resource ID for a Scoped Budget
type ScopedBudgetId struct {
	Scope      string
	BudgetName string
}

// NewScopedBudgetID returns a new ScopedBudgetId struct
func NewScopedBudgetID(scope string, budgetName string) ScopedBudgetId {
	return ScopedBudgetId{
		Scope:      scope,
		BudgetName: budgetName,
	}
}

// ParseScopedBudgetID parses 'input' into a ScopedBudgetId
func ParseScopedBudgetID(input string) (*ScopedBudgetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedBudgetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedBudgetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedBudgetIDInsensitively parses 'input' case-insensitively into a ScopedBudgetId
// note: this method should only be used for API response data and not user input
func ParseScopedBudgetIDInsensitively(input string) (*ScopedBudgetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedBudgetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedBudgetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedBudgetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.BudgetName, ok = input.Parsed["budgetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "budgetName", input)
	}

	return nil
}

// ValidateScopedBudgetID checks that 'input' can be parsed as a Scoped Budget ID
func ValidateScopedBudgetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedBudgetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Budget ID
func (id ScopedBudgetId) ID() string {
	fmtString := "/%s/providers/Microsoft.Consumption/budgets/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.BudgetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Budget ID
func (id ScopedBudgetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftConsumption", "Microsoft.Consumption", "Microsoft.Consumption"),
		resourceids.StaticSegment("staticBudgets", "budgets", "budgets"),
		resourceids.UserSpecifiedSegment("budgetName", "budgetName"),
	}
}

// String returns a human-readable description of this Scoped Budget ID
func (id ScopedBudgetId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Budget Name: %q", id.BudgetName),
	}
	return fmt.Sprintf("Scoped Budget (%s)", strings.Join(components, "\n"))
}
