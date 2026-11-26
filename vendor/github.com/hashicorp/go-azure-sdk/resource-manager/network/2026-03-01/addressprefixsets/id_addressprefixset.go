package addressprefixsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AddressPrefixSetId{})
}

var _ resourceids.ResourceId = &AddressPrefixSetId{}

// AddressPrefixSetId is a struct representing the Resource ID for a Address Prefix Set
type AddressPrefixSetId struct {
	SubscriptionId               string
	ResourceGroupName            string
	ApplicationSecurityGroupName string
	AddressPrefixSetName         string
}

// NewAddressPrefixSetID returns a new AddressPrefixSetId struct
func NewAddressPrefixSetID(subscriptionId string, resourceGroupName string, applicationSecurityGroupName string, addressPrefixSetName string) AddressPrefixSetId {
	return AddressPrefixSetId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		ApplicationSecurityGroupName: applicationSecurityGroupName,
		AddressPrefixSetName:         addressPrefixSetName,
	}
}

// ParseAddressPrefixSetID parses 'input' into a AddressPrefixSetId
func ParseAddressPrefixSetID(input string) (*AddressPrefixSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AddressPrefixSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AddressPrefixSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAddressPrefixSetIDInsensitively parses 'input' case-insensitively into a AddressPrefixSetId
// note: this method should only be used for API response data and not user input
func ParseAddressPrefixSetIDInsensitively(input string) (*AddressPrefixSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AddressPrefixSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AddressPrefixSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AddressPrefixSetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ApplicationSecurityGroupName, ok = input.Parsed["applicationSecurityGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationSecurityGroupName", input)
	}

	if id.AddressPrefixSetName, ok = input.Parsed["addressPrefixSetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "addressPrefixSetName", input)
	}

	return nil
}

// ValidateAddressPrefixSetID checks that 'input' can be parsed as a Address Prefix Set ID
func ValidateAddressPrefixSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAddressPrefixSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Address Prefix Set ID
func (id AddressPrefixSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationSecurityGroups/%s/addressPrefixSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationSecurityGroupName, id.AddressPrefixSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Address Prefix Set ID
func (id AddressPrefixSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticApplicationSecurityGroups", "applicationSecurityGroups", "applicationSecurityGroups"),
		resourceids.UserSpecifiedSegment("applicationSecurityGroupName", "applicationSecurityGroupName"),
		resourceids.StaticSegment("staticAddressPrefixSets", "addressPrefixSets", "addressPrefixSets"),
		resourceids.UserSpecifiedSegment("addressPrefixSetName", "addressPrefixSetName"),
	}
}

// String returns a human-readable description of this Address Prefix Set ID
func (id AddressPrefixSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Security Group Name: %q", id.ApplicationSecurityGroupName),
		fmt.Sprintf("Address Prefix Set Name: %q", id.AddressPrefixSetName),
	}
	return fmt.Sprintf("Address Prefix Set (%s)", strings.Join(components, "\n"))
}
