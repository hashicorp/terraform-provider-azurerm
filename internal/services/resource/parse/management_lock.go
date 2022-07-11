package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ManagementLockId{}

// ManagementLockId is a struct representing the Resource ID for a Management Lock
type ManagementLockId struct {
	Scope string
	Name  string
}

// NewManagementLockID returns a new ManagementLockId struct
func NewManagementLockID(scope string, name string) ManagementLockId {
	return ManagementLockId{
		Scope: scope,
		Name:  name,
	}
}

// ParseManagementLockID parses 'input' into a ManagementLockId
func ParseManagementLockID(input string) (*ManagementLockId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagementLockId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagementLockId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, fmt.Errorf("the segment 'scope' was not found in the resource id %q", input)
	}

	if id.Name, ok = parsed.Parsed["lockName"]; !ok {
		return nil, fmt.Errorf("the segment 'lockName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateManagementLockID checks that 'input' can be parsed as a Management Lock ID
func ValidateManagementLockID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagementLockID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Management Lock ID
func (id ManagementLockId) ID() string {
	fmtString := "%s/providers/Microsoft.Authorization/locks/%s"
	return fmt.Sprintf(fmtString, id.Scope, id.Name)
}

// Segments returns a slice of Resource ID Segments which comprise this Management Lock ID
func (id ManagementLockId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("locks", "locks", "locks"),
		resourceids.UserSpecifiedSegment("lockName", "lockValue"),
	}
}

// String returns a human-readable description of this Management Lock ID
func (id ManagementLockId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Name: %q", id.Name),
	}
	return fmt.Sprintf("Management Lock (%s)", strings.Join(components, "\n"))
}
