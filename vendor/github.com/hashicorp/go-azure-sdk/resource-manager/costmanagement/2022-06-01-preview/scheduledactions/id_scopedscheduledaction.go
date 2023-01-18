package scheduledactions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ScopedScheduledActionId{}

// ScopedScheduledActionId is a struct representing the Resource ID for a Scoped Scheduled Action
type ScopedScheduledActionId struct {
	Scope string
	Name  string
}

// NewScopedScheduledActionID returns a new ScopedScheduledActionId struct
func NewScopedScheduledActionID(scope string, name string) ScopedScheduledActionId {
	return ScopedScheduledActionId{
		Scope: scope,
		Name:  name,
	}
}

// ParseScopedScheduledActionID parses 'input' into a ScopedScheduledActionId
func ParseScopedScheduledActionID(input string) (*ScopedScheduledActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedScheduledActionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedScheduledActionId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, fmt.Errorf("the segment 'scope' was not found in the resource id %q", input)
	}

	if id.Name, ok = parsed.Parsed["name"]; !ok {
		return nil, fmt.Errorf("the segment 'name' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseScopedScheduledActionIDInsensitively parses 'input' case-insensitively into a ScopedScheduledActionId
// note: this method should only be used for API response data and not user input
func ParseScopedScheduledActionIDInsensitively(input string) (*ScopedScheduledActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedScheduledActionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedScheduledActionId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, fmt.Errorf("the segment 'scope' was not found in the resource id %q", input)
	}

	if id.Name, ok = parsed.Parsed["name"]; !ok {
		return nil, fmt.Errorf("the segment 'name' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateScopedScheduledActionID checks that 'input' can be parsed as a Scoped Scheduled Action ID
func ValidateScopedScheduledActionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedScheduledActionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Scheduled Action ID
func (id ScopedScheduledActionId) ID() string {
	fmtString := "/%s/providers/Microsoft.CostManagement/scheduledActions/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.Name)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Scheduled Action ID
func (id ScopedScheduledActionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCostManagement", "Microsoft.CostManagement", "Microsoft.CostManagement"),
		resourceids.StaticSegment("staticScheduledActions", "scheduledActions", "scheduledActions"),
		resourceids.UserSpecifiedSegment("name", "nameValue"),
	}
}

// String returns a human-readable description of this Scoped Scheduled Action ID
func (id ScopedScheduledActionId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Name: %q", id.Name),
	}
	return fmt.Sprintf("Scoped Scheduled Action (%s)", strings.Join(components, "\n"))
}
