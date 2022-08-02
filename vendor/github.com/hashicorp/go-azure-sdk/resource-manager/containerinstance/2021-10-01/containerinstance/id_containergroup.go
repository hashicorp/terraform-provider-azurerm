package containerinstance

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ContainerGroupId{}

// ContainerGroupId is a struct representing the Resource ID for a Container Group
type ContainerGroupId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ContainerGroupName string
}

// NewContainerGroupID returns a new ContainerGroupId struct
func NewContainerGroupID(subscriptionId string, resourceGroupName string, containerGroupName string) ContainerGroupId {
	return ContainerGroupId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ContainerGroupName: containerGroupName,
	}
}

// ParseContainerGroupID parses 'input' into a ContainerGroupId
func ParseContainerGroupID(input string) (*ContainerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ContainerGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ContainerGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ContainerGroupName, ok = parsed.Parsed["containerGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'containerGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseContainerGroupIDInsensitively parses 'input' case-insensitively into a ContainerGroupId
// note: this method should only be used for API response data and not user input
func ParseContainerGroupIDInsensitively(input string) (*ContainerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ContainerGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ContainerGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ContainerGroupName, ok = parsed.Parsed["containerGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'containerGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateContainerGroupID checks that 'input' can be parsed as a Container Group ID
func ValidateContainerGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContainerGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Container Group ID
func (id ContainerGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerInstance/containerGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Container Group ID
func (id ContainerGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerInstance", "Microsoft.ContainerInstance", "Microsoft.ContainerInstance"),
		resourceids.StaticSegment("staticContainerGroups", "containerGroups", "containerGroups"),
		resourceids.UserSpecifiedSegment("containerGroupName", "containerGroupValue"),
	}
}

// String returns a human-readable description of this Container Group ID
func (id ContainerGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container Group Name: %q", id.ContainerGroupName),
	}
	return fmt.Sprintf("Container Group (%s)", strings.Join(components, "\n"))
}
