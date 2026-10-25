package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VpnServerConfigurationId{})
}

var _ resourceids.ResourceId = &VpnServerConfigurationId{}

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
	parser := resourceids.NewParserFromResourceIdType(&VpnServerConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VpnServerConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVpnServerConfigurationIDInsensitively parses 'input' case-insensitively into a VpnServerConfigurationId
// note: this method should only be used for API response data and not user input
func ParseVpnServerConfigurationIDInsensitively(input string) (*VpnServerConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VpnServerConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VpnServerConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VpnServerConfigurationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VpnServerConfigurationName, ok = input.Parsed["vpnServerConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vpnServerConfigurationName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("vpnServerConfigurationName", "vpnServerConfigurationName"),
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
