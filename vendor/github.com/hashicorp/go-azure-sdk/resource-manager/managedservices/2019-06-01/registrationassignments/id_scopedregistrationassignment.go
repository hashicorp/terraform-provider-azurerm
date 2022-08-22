package registrationassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ScopedRegistrationAssignmentId{}

// ScopedRegistrationAssignmentId is a struct representing the Resource ID for a Scoped Registration Assignment
type ScopedRegistrationAssignmentId struct {
	Scope                    string
	RegistrationAssignmentId string
}

// NewScopedRegistrationAssignmentID returns a new ScopedRegistrationAssignmentId struct
func NewScopedRegistrationAssignmentID(scope string, registrationAssignmentId string) ScopedRegistrationAssignmentId {
	return ScopedRegistrationAssignmentId{
		Scope:                    scope,
		RegistrationAssignmentId: registrationAssignmentId,
	}
}

// ParseScopedRegistrationAssignmentID parses 'input' into a ScopedRegistrationAssignmentId
func ParseScopedRegistrationAssignmentID(input string) (*ScopedRegistrationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRegistrationAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRegistrationAssignmentId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, fmt.Errorf("the segment 'scope' was not found in the resource id %q", input)
	}

	if id.RegistrationAssignmentId, ok = parsed.Parsed["registrationAssignmentId"]; !ok {
		return nil, fmt.Errorf("the segment 'registrationAssignmentId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseScopedRegistrationAssignmentIDInsensitively parses 'input' case-insensitively into a ScopedRegistrationAssignmentId
// note: this method should only be used for API response data and not user input
func ParseScopedRegistrationAssignmentIDInsensitively(input string) (*ScopedRegistrationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRegistrationAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRegistrationAssignmentId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, fmt.Errorf("the segment 'scope' was not found in the resource id %q", input)
	}

	if id.RegistrationAssignmentId, ok = parsed.Parsed["registrationAssignmentId"]; !ok {
		return nil, fmt.Errorf("the segment 'registrationAssignmentId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateScopedRegistrationAssignmentID checks that 'input' can be parsed as a Scoped Registration Assignment ID
func ValidateScopedRegistrationAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRegistrationAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Registration Assignment ID
func (id ScopedRegistrationAssignmentId) ID() string {
	fmtString := "/%s/providers/Microsoft.ManagedServices/registrationAssignments/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RegistrationAssignmentId)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Registration Assignment ID
func (id ScopedRegistrationAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagedServices", "Microsoft.ManagedServices", "Microsoft.ManagedServices"),
		resourceids.StaticSegment("staticRegistrationAssignments", "registrationAssignments", "registrationAssignments"),
		resourceids.UserSpecifiedSegment("registrationAssignmentId", "registrationAssignmentIdValue"),
	}
}

// String returns a human-readable description of this Scoped Registration Assignment ID
func (id ScopedRegistrationAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Registration Assignment: %q", id.RegistrationAssignmentId),
	}
	return fmt.Sprintf("Scoped Registration Assignment (%s)", strings.Join(components, "\n"))
}
