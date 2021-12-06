package managedidentity

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = UserAssignedIdentitiesId{}

// UserAssignedIdentitiesId is a struct representing the Resource ID for a User Assigned Identities
type UserAssignedIdentitiesId struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceName      string
}

// NewUserAssignedIdentitiesID returns a new UserAssignedIdentitiesId struct
func NewUserAssignedIdentitiesID(subscriptionId string, resourceGroupName string, resourceName string) UserAssignedIdentitiesId {
	return UserAssignedIdentitiesId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ResourceName:      resourceName,
	}
}

// ParseUserAssignedIdentitiesID parses 'input' into a UserAssignedIdentitiesId
func ParseUserAssignedIdentitiesID(input string) (*UserAssignedIdentitiesId, error) {
	parser := resourceids.NewParserFromResourceIdType(UserAssignedIdentitiesId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UserAssignedIdentitiesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseUserAssignedIdentitiesIDInsensitively parses 'input' case-insensitively into a UserAssignedIdentitiesId
// note: this method should only be used for API response data and not user input
func ParseUserAssignedIdentitiesIDInsensitively(input string) (*UserAssignedIdentitiesId, error) {
	parser := resourceids.NewParserFromResourceIdType(UserAssignedIdentitiesId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UserAssignedIdentitiesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateUserAssignedIdentitiesID checks that 'input' can be parsed as a User Assigned Identities ID
func ValidateUserAssignedIdentitiesID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUserAssignedIdentitiesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted User Assigned Identities ID
func (id UserAssignedIdentitiesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this User Assigned Identities ID
func (id UserAssignedIdentitiesId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagedIdentity", "Microsoft.ManagedIdentity", "Microsoft.ManagedIdentity"),
		resourceids.StaticSegment("staticUserAssignedIdentities", "userAssignedIdentities", "userAssignedIdentities"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
	}
}

// String returns a human-readable description of this User Assigned Identities ID
func (id UserAssignedIdentitiesId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
	}
	return fmt.Sprintf("User Assigned Identities (%s)", strings.Join(components, "\n"))
}
