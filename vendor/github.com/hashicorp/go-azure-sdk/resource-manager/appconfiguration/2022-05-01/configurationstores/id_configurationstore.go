package configurationstores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ConfigurationStoreId{}

// ConfigurationStoreId is a struct representing the Resource ID for a Configuration Store
type ConfigurationStoreId struct {
	SubscriptionId    string
	ResourceGroupName string
	ConfigStoreName   string
}

// NewConfigurationStoreID returns a new ConfigurationStoreId struct
func NewConfigurationStoreID(subscriptionId string, resourceGroupName string, configStoreName string) ConfigurationStoreId {
	return ConfigurationStoreId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ConfigStoreName:   configStoreName,
	}
}

// ParseConfigurationStoreID parses 'input' into a ConfigurationStoreId
func ParseConfigurationStoreID(input string) (*ConfigurationStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationStoreId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationStoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ConfigStoreName, ok = parsed.Parsed["configStoreName"]; !ok {
		return nil, fmt.Errorf("the segment 'configStoreName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseConfigurationStoreIDInsensitively parses 'input' case-insensitively into a ConfigurationStoreId
// note: this method should only be used for API response data and not user input
func ParseConfigurationStoreIDInsensitively(input string) (*ConfigurationStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationStoreId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationStoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ConfigStoreName, ok = parsed.Parsed["configStoreName"]; !ok {
		return nil, fmt.Errorf("the segment 'configStoreName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateConfigurationStoreID checks that 'input' can be parsed as a Configuration Store ID
func ValidateConfigurationStoreID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationStoreID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration Store ID
func (id ConfigurationStoreId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppConfiguration/configurationStores/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ConfigStoreName)
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration Store ID
func (id ConfigurationStoreId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppConfiguration", "Microsoft.AppConfiguration", "Microsoft.AppConfiguration"),
		resourceids.StaticSegment("staticConfigurationStores", "configurationStores", "configurationStores"),
		resourceids.UserSpecifiedSegment("configStoreName", "configStoreValue"),
	}
}

// String returns a human-readable description of this Configuration Store ID
func (id ConfigurationStoreId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Config Store Name: %q", id.ConfigStoreName),
	}
	return fmt.Sprintf("Configuration Store (%s)", strings.Join(components, "\n"))
}
