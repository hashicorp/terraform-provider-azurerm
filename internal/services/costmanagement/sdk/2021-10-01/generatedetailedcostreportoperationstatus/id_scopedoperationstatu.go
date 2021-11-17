package generatedetailedcostreportoperationstatus

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ScopedOperationStatuId{}

// ScopedOperationStatuId is a struct representing the Resource ID for a Scoped Operation Statu
type ScopedOperationStatuId struct {
	Scope       string
	OperationId string
}

// NewScopedOperationStatuID returns a new ScopedOperationStatuId struct
func NewScopedOperationStatuID(scope string, operationId string) ScopedOperationStatuId {
	return ScopedOperationStatuId{
		Scope:       scope,
		OperationId: operationId,
	}
}

// ParseScopedOperationStatuID parses 'input' into a ScopedOperationStatuId
func ParseScopedOperationStatuID(input string) (*ScopedOperationStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedOperationStatuId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedOperationStatuId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, fmt.Errorf("the segment 'scope' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseScopedOperationStatuIDInsensitively parses 'input' case-insensitively into a ScopedOperationStatuId
// note: this method should only be used for API response data and not user input
func ParseScopedOperationStatuIDInsensitively(input string) (*ScopedOperationStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedOperationStatuId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedOperationStatuId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, fmt.Errorf("the segment 'scope' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateScopedOperationStatuID checks that 'input' can be parsed as a Scoped Operation Statu ID
func ValidateScopedOperationStatuID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedOperationStatuID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Operation Statu ID
func (id ScopedOperationStatuId) ID() string {
	fmtString := "/%s/providers/Microsoft.CostManagement/operationStatus/%s"
	return fmt.Sprintf(fmtString, id.Scope, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Operation Statu ID
func (id ScopedOperationStatuId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCostManagement", "Microsoft.CostManagement", "Microsoft.CostManagement"),
		resourceids.StaticSegment("staticOperationStatus", "operationStatus", "operationStatus"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Scoped Operation Statu ID
func (id ScopedOperationStatuId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Scoped Operation Statu (%s)", strings.Join(components, "\n"))
}
