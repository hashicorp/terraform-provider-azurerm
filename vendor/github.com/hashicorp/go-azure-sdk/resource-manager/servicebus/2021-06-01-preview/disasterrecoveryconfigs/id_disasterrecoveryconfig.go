package disasterrecoveryconfigs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DisasterRecoveryConfigId{}

// DisasterRecoveryConfigId is a struct representing the Resource ID for a Disaster Recovery Config
type DisasterRecoveryConfigId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	Alias             string
}

// NewDisasterRecoveryConfigID returns a new DisasterRecoveryConfigId struct
func NewDisasterRecoveryConfigID(subscriptionId string, resourceGroupName string, namespaceName string, alias string) DisasterRecoveryConfigId {
	return DisasterRecoveryConfigId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		Alias:             alias,
	}
}

// ParseDisasterRecoveryConfigID parses 'input' into a DisasterRecoveryConfigId
func ParseDisasterRecoveryConfigID(input string) (*DisasterRecoveryConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(DisasterRecoveryConfigId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DisasterRecoveryConfigId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.Alias, ok = parsed.Parsed["alias"]; !ok {
		return nil, fmt.Errorf("the segment 'alias' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDisasterRecoveryConfigIDInsensitively parses 'input' case-insensitively into a DisasterRecoveryConfigId
// note: this method should only be used for API response data and not user input
func ParseDisasterRecoveryConfigIDInsensitively(input string) (*DisasterRecoveryConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(DisasterRecoveryConfigId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DisasterRecoveryConfigId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.Alias, ok = parsed.Parsed["alias"]; !ok {
		return nil, fmt.Errorf("the segment 'alias' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDisasterRecoveryConfigID checks that 'input' can be parsed as a Disaster Recovery Config ID
func ValidateDisasterRecoveryConfigID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDisasterRecoveryConfigID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Disaster Recovery Config ID
func (id DisasterRecoveryConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/disasterRecoveryConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.Alias)
}

// Segments returns a slice of Resource ID Segments which comprise this Disaster Recovery Config ID
func (id DisasterRecoveryConfigId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceBus", "Microsoft.ServiceBus", "Microsoft.ServiceBus"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticDisasterRecoveryConfigs", "disasterRecoveryConfigs", "disasterRecoveryConfigs"),
		resourceids.UserSpecifiedSegment("alias", "aliasValue"),
	}
}

// String returns a human-readable description of this Disaster Recovery Config ID
func (id DisasterRecoveryConfigId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Alias: %q", id.Alias),
	}
	return fmt.Sprintf("Disaster Recovery Config (%s)", strings.Join(components, "\n"))
}
