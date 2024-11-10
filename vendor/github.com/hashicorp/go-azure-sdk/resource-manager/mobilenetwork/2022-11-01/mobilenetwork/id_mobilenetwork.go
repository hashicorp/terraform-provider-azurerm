package mobilenetwork

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MobileNetworkId{})
}

var _ resourceids.ResourceId = &MobileNetworkId{}

// MobileNetworkId is a struct representing the Resource ID for a Mobile Network
type MobileNetworkId struct {
	SubscriptionId    string
	ResourceGroupName string
	MobileNetworkName string
}

// NewMobileNetworkID returns a new MobileNetworkId struct
func NewMobileNetworkID(subscriptionId string, resourceGroupName string, mobileNetworkName string) MobileNetworkId {
	return MobileNetworkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MobileNetworkName: mobileNetworkName,
	}
}

// ParseMobileNetworkID parses 'input' into a MobileNetworkId
func ParseMobileNetworkID(input string) (*MobileNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MobileNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MobileNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMobileNetworkIDInsensitively parses 'input' case-insensitively into a MobileNetworkId
// note: this method should only be used for API response data and not user input
func ParseMobileNetworkIDInsensitively(input string) (*MobileNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MobileNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MobileNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MobileNetworkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MobileNetworkName, ok = input.Parsed["mobileNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "mobileNetworkName", input)
	}

	return nil
}

// ValidateMobileNetworkID checks that 'input' can be parsed as a Mobile Network ID
func ValidateMobileNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMobileNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mobile Network ID
func (id MobileNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MobileNetwork/mobileNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Mobile Network ID
func (id MobileNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMobileNetwork", "Microsoft.MobileNetwork", "Microsoft.MobileNetwork"),
		resourceids.StaticSegment("staticMobileNetworks", "mobileNetworks", "mobileNetworks"),
		resourceids.UserSpecifiedSegment("mobileNetworkName", "mobileNetworkName"),
	}
}

// String returns a human-readable description of this Mobile Network ID
func (id MobileNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Mobile Network Name: %q", id.MobileNetworkName),
	}
	return fmt.Sprintf("Mobile Network (%s)", strings.Join(components, "\n"))
}
