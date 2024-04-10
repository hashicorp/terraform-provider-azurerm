package connectivityconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ConnectivityConfigurationId{}

// ConnectivityConfigurationId is a struct representing the Resource ID for a Connectivity Configuration
type ConnectivityConfigurationId struct {
	SubscriptionId                string
	ResourceGroupName             string
	NetworkManagerName            string
	ConnectivityConfigurationName string
}

// NewConnectivityConfigurationID returns a new ConnectivityConfigurationId struct
func NewConnectivityConfigurationID(subscriptionId string, resourceGroupName string, networkManagerName string, connectivityConfigurationName string) ConnectivityConfigurationId {
	return ConnectivityConfigurationId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		NetworkManagerName:            networkManagerName,
		ConnectivityConfigurationName: connectivityConfigurationName,
	}
}

// ParseConnectivityConfigurationID parses 'input' into a ConnectivityConfigurationId
func ParseConnectivityConfigurationID(input string) (*ConnectivityConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectivityConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectivityConfigurationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectivityConfigurationIDInsensitively parses 'input' case-insensitively into a ConnectivityConfigurationId
// note: this method should only be used for API response data and not user input
func ParseConnectivityConfigurationIDInsensitively(input string) (*ConnectivityConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectivityConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectivityConfigurationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectivityConfigurationId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ConnectivityConfigurationName, ok = input.Parsed["connectivityConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectivityConfigurationName", input)
	}

	return nil
}

// ValidateConnectivityConfigurationID checks that 'input' can be parsed as a Connectivity Configuration ID
func ValidateConnectivityConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectivityConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connectivity Configuration ID
func (id ConnectivityConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/connectivityConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.ConnectivityConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connectivity Configuration ID
func (id ConnectivityConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerValue"),
		resourceids.StaticSegment("staticConnectivityConfigurations", "connectivityConfigurations", "connectivityConfigurations"),
		resourceids.UserSpecifiedSegment("connectivityConfigurationName", "connectivityConfigurationValue"),
	}
}

// String returns a human-readable description of this Connectivity Configuration ID
func (id ConnectivityConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Connectivity Configuration Name: %q", id.ConnectivityConfigurationName),
	}
	return fmt.Sprintf("Connectivity Configuration (%s)", strings.Join(components, "\n"))
}
