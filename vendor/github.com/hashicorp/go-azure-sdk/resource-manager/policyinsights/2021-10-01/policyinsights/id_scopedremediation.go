package policyinsights

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ScopedRemediationId{}

// ScopedRemediationId is a struct representing the Resource ID for a Scoped Remediation
type ScopedRemediationId struct {
	ResourceId      string
	RemediationName string
}

// NewScopedRemediationID returns a new ScopedRemediationId struct
func NewScopedRemediationID(resourceId string, remediationName string) ScopedRemediationId {
	return ScopedRemediationId{
		ResourceId:      resourceId,
		RemediationName: remediationName,
	}
}

// ParseScopedRemediationID parses 'input' into a ScopedRemediationId
func ParseScopedRemediationID(input string) (*ScopedRemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRemediationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRemediationId{}

	if id.ResourceId, ok = parsed.Parsed["resourceId"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceId' was not found in the resource id %q", input)
	}

	if id.RemediationName, ok = parsed.Parsed["remediationName"]; !ok {
		return nil, fmt.Errorf("the segment 'remediationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseScopedRemediationIDInsensitively parses 'input' case-insensitively into a ScopedRemediationId
// note: this method should only be used for API response data and not user input
func ParseScopedRemediationIDInsensitively(input string) (*ScopedRemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRemediationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRemediationId{}

	if id.ResourceId, ok = parsed.Parsed["resourceId"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceId' was not found in the resource id %q", input)
	}

	if id.RemediationName, ok = parsed.Parsed["remediationName"]; !ok {
		return nil, fmt.Errorf("the segment 'remediationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateScopedRemediationID checks that 'input' can be parsed as a Scoped Remediation ID
func ValidateScopedRemediationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRemediationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Remediation ID
func (id ScopedRemediationId) ID() string {
	fmtString := "/%s/providers/Microsoft.PolicyInsights/remediations/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceId, "/"), id.RemediationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Remediation ID
func (id ScopedRemediationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceId", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPolicyInsights", "Microsoft.PolicyInsights", "Microsoft.PolicyInsights"),
		resourceids.StaticSegment("staticRemediations", "remediations", "remediations"),
		resourceids.UserSpecifiedSegment("remediationName", "remediationValue"),
	}
}

// String returns a human-readable description of this Scoped Remediation ID
func (id ScopedRemediationId) String() string {
	components := []string{
		fmt.Sprintf("Resource: %q", id.ResourceId),
		fmt.Sprintf("Remediation Name: %q", id.RemediationName),
	}
	return fmt.Sprintf("Scoped Remediation (%s)", strings.Join(components, "\n"))
}
