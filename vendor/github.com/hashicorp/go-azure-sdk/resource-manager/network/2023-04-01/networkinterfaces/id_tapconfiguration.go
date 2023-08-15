package networkinterfaces

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TapConfigurationId{}

// TapConfigurationId is a struct representing the Resource ID for a Tap Configuration
type TapConfigurationId struct {
	SubscriptionId       string
	ResourceGroupName    string
	NetworkInterfaceName string
	TapConfigurationName string
}

// NewTapConfigurationID returns a new TapConfigurationId struct
func NewTapConfigurationID(subscriptionId string, resourceGroupName string, networkInterfaceName string, tapConfigurationName string) TapConfigurationId {
	return TapConfigurationId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		NetworkInterfaceName: networkInterfaceName,
		TapConfigurationName: tapConfigurationName,
	}
}

// ParseTapConfigurationID parses 'input' into a TapConfigurationId
func ParseTapConfigurationID(input string) (*TapConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(TapConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TapConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkInterfaceName, ok = parsed.Parsed["networkInterfaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkInterfaceName", *parsed)
	}

	if id.TapConfigurationName, ok = parsed.Parsed["tapConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "tapConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseTapConfigurationIDInsensitively parses 'input' case-insensitively into a TapConfigurationId
// note: this method should only be used for API response data and not user input
func ParseTapConfigurationIDInsensitively(input string) (*TapConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(TapConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TapConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkInterfaceName, ok = parsed.Parsed["networkInterfaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkInterfaceName", *parsed)
	}

	if id.TapConfigurationName, ok = parsed.Parsed["tapConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "tapConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateTapConfigurationID checks that 'input' can be parsed as a Tap Configuration ID
func ValidateTapConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTapConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Tap Configuration ID
func (id TapConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkInterfaces/%s/tapConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkInterfaceName, id.TapConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Tap Configuration ID
func (id TapConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkInterfaces", "networkInterfaces", "networkInterfaces"),
		resourceids.UserSpecifiedSegment("networkInterfaceName", "networkInterfaceValue"),
		resourceids.StaticSegment("staticTapConfigurations", "tapConfigurations", "tapConfigurations"),
		resourceids.UserSpecifiedSegment("tapConfigurationName", "tapConfigurationValue"),
	}
}

// String returns a human-readable description of this Tap Configuration ID
func (id TapConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Interface Name: %q", id.NetworkInterfaceName),
		fmt.Sprintf("Tap Configuration Name: %q", id.TapConfigurationName),
	}
	return fmt.Sprintf("Tap Configuration (%s)", strings.Join(components, "\n"))
}
