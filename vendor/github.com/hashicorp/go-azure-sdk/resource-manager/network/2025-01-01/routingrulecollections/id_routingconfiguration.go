package routingrulecollections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RoutingConfigurationId{})
}

var _ resourceids.ResourceId = &RoutingConfigurationId{}

// RoutingConfigurationId is a struct representing the Resource ID for a Routing Configuration
type RoutingConfigurationId struct {
	SubscriptionId           string
	ResourceGroupName        string
	NetworkManagerName       string
	RoutingConfigurationName string
}

// NewRoutingConfigurationID returns a new RoutingConfigurationId struct
func NewRoutingConfigurationID(subscriptionId string, resourceGroupName string, networkManagerName string, routingConfigurationName string) RoutingConfigurationId {
	return RoutingConfigurationId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		NetworkManagerName:       networkManagerName,
		RoutingConfigurationName: routingConfigurationName,
	}
}

// ParseRoutingConfigurationID parses 'input' into a RoutingConfigurationId
func ParseRoutingConfigurationID(input string) (*RoutingConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoutingConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoutingConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRoutingConfigurationIDInsensitively parses 'input' case-insensitively into a RoutingConfigurationId
// note: this method should only be used for API response data and not user input
func ParseRoutingConfigurationIDInsensitively(input string) (*RoutingConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoutingConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoutingConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RoutingConfigurationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkManagerName, ok = input.Parsed["networkManagerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", input)
	}

	if id.RoutingConfigurationName, ok = input.Parsed["routingConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routingConfigurationName", input)
	}

	return nil
}

// ValidateRoutingConfigurationID checks that 'input' can be parsed as a Routing Configuration ID
func ValidateRoutingConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRoutingConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Routing Configuration ID
func (id RoutingConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/routingConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.RoutingConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Routing Configuration ID
func (id RoutingConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerName"),
		resourceids.StaticSegment("staticRoutingConfigurations", "routingConfigurations", "routingConfigurations"),
		resourceids.UserSpecifiedSegment("routingConfigurationName", "routingConfigurationName"),
	}
}

// String returns a human-readable description of this Routing Configuration ID
func (id RoutingConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Routing Configuration Name: %q", id.RoutingConfigurationName),
	}
	return fmt.Sprintf("Routing Configuration (%s)", strings.Join(components, "\n"))
}
