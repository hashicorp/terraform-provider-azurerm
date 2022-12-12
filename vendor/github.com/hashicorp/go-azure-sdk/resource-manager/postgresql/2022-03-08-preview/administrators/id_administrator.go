package administrators

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AdministratorId{}

// AdministratorId is a struct representing the Resource ID for a Administrator
type AdministratorId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
	ObjectId          string
}

// NewAdministratorID returns a new AdministratorId struct
func NewAdministratorID(subscriptionId string, resourceGroupName string, serverName string, objectId string) AdministratorId {
	return AdministratorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
		ObjectId:          objectId,
	}
}

// ParseAdministratorID parses 'input' into a AdministratorId
func ParseAdministratorID(input string) (*AdministratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(AdministratorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AdministratorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, fmt.Errorf("the segment 'serverName' was not found in the resource id %q", input)
	}

	if id.ObjectId, ok = parsed.Parsed["objectId"]; !ok {
		return nil, fmt.Errorf("the segment 'objectId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseAdministratorIDInsensitively parses 'input' case-insensitively into a AdministratorId
// note: this method should only be used for API response data and not user input
func ParseAdministratorIDInsensitively(input string) (*AdministratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(AdministratorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AdministratorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, fmt.Errorf("the segment 'serverName' was not found in the resource id %q", input)
	}

	if id.ObjectId, ok = parsed.Parsed["objectId"]; !ok {
		return nil, fmt.Errorf("the segment 'objectId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateAdministratorID checks that 'input' can be parsed as a Administrator ID
func ValidateAdministratorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAdministratorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Administrator ID
func (id AdministratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/flexibleServers/%s/administrators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.ObjectId)
}

// Segments returns a slice of Resource ID Segments which comprise this Administrator ID
func (id AdministratorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticFlexibleServers", "flexibleServers", "flexibleServers"),
		resourceids.UserSpecifiedSegment("serverName", "serverValue"),
		resourceids.StaticSegment("staticAdministrators", "administrators", "administrators"),
		resourceids.UserSpecifiedSegment("objectId", "objectIdValue"),
	}
}

// String returns a human-readable description of this Administrator ID
func (id AdministratorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Object: %q", id.ObjectId),
	}
	return fmt.Sprintf("Administrator (%s)", strings.Join(components, "\n"))
}
