package alerts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ScopedAlertId{}

// ScopedAlertId is a struct representing the Resource ID for a Scoped Alert
type ScopedAlertId struct {
	Scope   string
	AlertId string
}

// NewScopedAlertID returns a new ScopedAlertId struct
func NewScopedAlertID(scope string, alertId string) ScopedAlertId {
	return ScopedAlertId{
		Scope:   scope,
		AlertId: alertId,
	}
}

// ParseScopedAlertID parses 'input' into a ScopedAlertId
func ParseScopedAlertID(input string) (*ScopedAlertId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedAlertId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedAlertId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, fmt.Errorf("the segment 'scope' was not found in the resource id %q", input)
	}

	if id.AlertId, ok = parsed.Parsed["alertId"]; !ok {
		return nil, fmt.Errorf("the segment 'alertId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseScopedAlertIDInsensitively parses 'input' case-insensitively into a ScopedAlertId
// note: this method should only be used for API response data and not user input
func ParseScopedAlertIDInsensitively(input string) (*ScopedAlertId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedAlertId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedAlertId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, fmt.Errorf("the segment 'scope' was not found in the resource id %q", input)
	}

	if id.AlertId, ok = parsed.Parsed["alertId"]; !ok {
		return nil, fmt.Errorf("the segment 'alertId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateScopedAlertID checks that 'input' can be parsed as a Scoped Alert ID
func ValidateScopedAlertID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedAlertID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Alert ID
func (id ScopedAlertId) ID() string {
	fmtString := "/%s/providers/Microsoft.CostManagement/alerts/%s"
	return fmt.Sprintf(fmtString, id.Scope, id.AlertId)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Alert ID
func (id ScopedAlertId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCostManagement", "Microsoft.CostManagement", "Microsoft.CostManagement"),
		resourceids.StaticSegment("staticAlerts", "alerts", "alerts"),
		resourceids.UserSpecifiedSegment("alertId", "alertIdValue"),
	}
}

// String returns a human-readable description of this Scoped Alert ID
func (id ScopedAlertId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Alert: %q", id.AlertId),
	}
	return fmt.Sprintf("Scoped Alert (%s)", strings.Join(components, "\n"))
}
