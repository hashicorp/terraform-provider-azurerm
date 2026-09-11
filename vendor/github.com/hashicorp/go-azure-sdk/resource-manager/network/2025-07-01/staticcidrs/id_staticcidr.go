package staticcidrs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StaticCidrId{})
}

var _ resourceids.ResourceId = &StaticCidrId{}

// StaticCidrId is a struct representing the Resource ID for a Static Cidr
type StaticCidrId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetworkManagerName string
	IpamPoolName       string
	StaticCidrName     string
}

// NewStaticCidrID returns a new StaticCidrId struct
func NewStaticCidrID(subscriptionId string, resourceGroupName string, networkManagerName string, ipamPoolName string, staticCidrName string) StaticCidrId {
	return StaticCidrId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetworkManagerName: networkManagerName,
		IpamPoolName:       ipamPoolName,
		StaticCidrName:     staticCidrName,
	}
}

// ParseStaticCidrID parses 'input' into a StaticCidrId
func ParseStaticCidrID(input string) (*StaticCidrId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StaticCidrId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StaticCidrId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStaticCidrIDInsensitively parses 'input' case-insensitively into a StaticCidrId
// note: this method should only be used for API response data and not user input
func ParseStaticCidrIDInsensitively(input string) (*StaticCidrId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StaticCidrId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StaticCidrId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StaticCidrId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.IpamPoolName, ok = input.Parsed["ipamPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "ipamPoolName", input)
	}

	if id.StaticCidrName, ok = input.Parsed["staticCidrName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "staticCidrName", input)
	}

	return nil
}

// ValidateStaticCidrID checks that 'input' can be parsed as a Static Cidr ID
func ValidateStaticCidrID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStaticCidrID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Static Cidr ID
func (id StaticCidrId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/ipamPools/%s/staticCidrs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.IpamPoolName, id.StaticCidrName)
}

// Segments returns a slice of Resource ID Segments which comprise this Static Cidr ID
func (id StaticCidrId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerName"),
		resourceids.StaticSegment("staticIpamPools", "ipamPools", "ipamPools"),
		resourceids.UserSpecifiedSegment("ipamPoolName", "ipamPoolName"),
		resourceids.StaticSegment("staticStaticCidrs", "staticCidrs", "staticCidrs"),
		resourceids.UserSpecifiedSegment("staticCidrName", "staticCidrName"),
	}
}

// String returns a human-readable description of this Static Cidr ID
func (id StaticCidrId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Ipam Pool Name: %q", id.IpamPoolName),
		fmt.Sprintf("Static Cidr Name: %q", id.StaticCidrName),
	}
	return fmt.Sprintf("Static Cidr (%s)", strings.Join(components, "\n"))
}
