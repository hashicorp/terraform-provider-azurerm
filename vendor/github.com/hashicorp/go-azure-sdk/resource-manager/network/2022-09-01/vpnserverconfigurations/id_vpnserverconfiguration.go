package vpnserverconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VpnServerConfigurationId{}

// VpnServerConfigurationId is a struct representing the Resource ID for a Vpn Server Configuration
type VpnServerConfigurationId struct {
	SubscriptionId             string
	ResourceGroupName          string
	VpnServerConfigurationName string
}

// NewVpnServerConfigurationID returns a new VpnServerConfigurationId struct
func NewVpnServerConfigurationID(subscriptionId string, resourceGroupName string, vpnServerConfigurationName string) VpnServerConfigurationId {
	return VpnServerConfigurationId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		VpnServerConfigurationName: vpnServerConfigurationName,
	}
}

// ParseVpnServerConfigurationID parses 'input' into a VpnServerConfigurationId
func ParseVpnServerConfigurationID(input string) (*VpnServerConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnServerConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnServerConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnServerConfigurationName, ok = parsed.Parsed["vpnServerConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnServerConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseVpnServerConfigurationIDInsensitively parses 'input' case-insensitively into a VpnServerConfigurationId
// note: this method should only be used for API response data and not user input
func ParseVpnServerConfigurationIDInsensitively(input string) (*VpnServerConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnServerConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnServerConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnServerConfigurationName, ok = parsed.Parsed["vpnServerConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnServerConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateVpnServerConfigurationID checks that 'input' can be parsed as a Vpn Server Configuration ID
func ValidateVpnServerConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVpnServerConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Vpn Server Configuration ID
func (id VpnServerConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnServerConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VpnServerConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Vpn Server Configuration ID
func (id VpnServerConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVpnServerConfigurations", "vpnServerConfigurations", "vpnServerConfigurations"),
		resourceids.UserSpecifiedSegment("vpnServerConfigurationName", "vpnServerConfigurationValue"),
	}
}

// String returns a human-readable description of this Vpn Server Configuration ID
func (id VpnServerConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vpn Server Configuration Name: %q", id.VpnServerConfigurationName),
	}
	return fmt.Sprintf("Vpn Server Configuration (%s)", strings.Join(components, "\n"))
}
