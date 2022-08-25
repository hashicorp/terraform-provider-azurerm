package configurationassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ConfigurationAssignmentId{}

// ConfigurationAssignmentId is a struct representing the Resource ID for a Configuration Assignment
type ConfigurationAssignmentId struct {
	SubscriptionId              string
	ResourceGroupName           string
	ProviderName                string
	ResourceType                string
	ResourceName                string
	ConfigurationAssignmentName string
}

// NewConfigurationAssignmentID returns a new ConfigurationAssignmentId struct
func NewConfigurationAssignmentID(subscriptionId string, resourceGroupName string, providerName string, resourceType string, resourceName string, configurationAssignmentName string) ConfigurationAssignmentId {
	return ConfigurationAssignmentId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		ProviderName:                providerName,
		ResourceType:                resourceType,
		ResourceName:                resourceName,
		ConfigurationAssignmentName: configurationAssignmentName,
	}
}

// ParseConfigurationAssignmentID parses 'input' into a ConfigurationAssignmentId
func ParseConfigurationAssignmentID(input string) (*ConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationAssignmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, fmt.Errorf("the segment 'providerName' was not found in the resource id %q", input)
	}

	if id.ResourceType, ok = parsed.Parsed["resourceType"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceType' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.ConfigurationAssignmentName, ok = parsed.Parsed["configurationAssignmentName"]; !ok {
		return nil, fmt.Errorf("the segment 'configurationAssignmentName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseConfigurationAssignmentIDInsensitively parses 'input' case-insensitively into a ConfigurationAssignmentId
// note: this method should only be used for API response data and not user input
func ParseConfigurationAssignmentIDInsensitively(input string) (*ConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationAssignmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, fmt.Errorf("the segment 'providerName' was not found in the resource id %q", input)
	}

	if id.ResourceType, ok = parsed.Parsed["resourceType"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceType' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.ConfigurationAssignmentName, ok = parsed.Parsed["configurationAssignmentName"]; !ok {
		return nil, fmt.Errorf("the segment 'configurationAssignmentName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateConfigurationAssignmentID checks that 'input' can be parsed as a Configuration Assignment ID
func ValidateConfigurationAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration Assignment ID
func (id ConfigurationAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s/providers/Microsoft.Maintenance/configurationAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProviderName, id.ResourceType, id.ResourceName, id.ConfigurationAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration Assignment ID
func (id ConfigurationAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.UserSpecifiedSegment("providerName", "providerValue"),
		resourceids.UserSpecifiedSegment("resourceType", "resourceTypeValue"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMaintenance", "Microsoft.Maintenance", "Microsoft.Maintenance"),
		resourceids.StaticSegment("staticConfigurationAssignments", "configurationAssignments", "configurationAssignments"),
		resourceids.UserSpecifiedSegment("configurationAssignmentName", "configurationAssignmentValue"),
	}
}

// String returns a human-readable description of this Configuration Assignment ID
func (id ConfigurationAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
		fmt.Sprintf("Resource Type: %q", id.ResourceType),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Configuration Assignment Name: %q", id.ConfigurationAssignmentName),
	}
	return fmt.Sprintf("Configuration Assignment (%s)", strings.Join(components, "\n"))
}
