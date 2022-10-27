package managedclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CommandResultId{}

// CommandResultId is a struct representing the Resource ID for a Command Result
type CommandResultId struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceName      string
	CommandId         string
}

// NewCommandResultID returns a new CommandResultId struct
func NewCommandResultID(subscriptionId string, resourceGroupName string, resourceName string, commandId string) CommandResultId {
	return CommandResultId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ResourceName:      resourceName,
		CommandId:         commandId,
	}
}

// ParseCommandResultID parses 'input' into a CommandResultId
func ParseCommandResultID(input string) (*CommandResultId, error) {
	parser := resourceids.NewParserFromResourceIdType(CommandResultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CommandResultId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.CommandId, ok = parsed.Parsed["commandId"]; !ok {
		return nil, fmt.Errorf("the segment 'commandId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCommandResultIDInsensitively parses 'input' case-insensitively into a CommandResultId
// note: this method should only be used for API response data and not user input
func ParseCommandResultIDInsensitively(input string) (*CommandResultId, error) {
	parser := resourceids.NewParserFromResourceIdType(CommandResultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CommandResultId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.CommandId, ok = parsed.Parsed["commandId"]; !ok {
		return nil, fmt.Errorf("the segment 'commandId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCommandResultID checks that 'input' can be parsed as a Command Result ID
func ValidateCommandResultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCommandResultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Command Result ID
func (id CommandResultId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s/commandResults/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceName, id.CommandId)
}

// Segments returns a slice of Resource ID Segments which comprise this Command Result ID
func (id CommandResultId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticManagedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
		resourceids.StaticSegment("staticCommandResults", "commandResults", "commandResults"),
		resourceids.UserSpecifiedSegment("commandId", "commandIdValue"),
	}
}

// String returns a human-readable description of this Command Result ID
func (id CommandResultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Command: %q", id.CommandId),
	}
	return fmt.Sprintf("Command Result (%s)", strings.Join(components, "\n"))
}
