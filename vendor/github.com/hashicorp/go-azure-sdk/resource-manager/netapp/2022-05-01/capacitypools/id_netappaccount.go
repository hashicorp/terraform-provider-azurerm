package capacitypools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NetAppAccountId{}

// NetAppAccountId is a struct representing the Resource ID for a Net App Account
type NetAppAccountId struct {
	SubscriptionId    string
	ResourceGroupName string
	NetAppAccountName string
}

// NewNetAppAccountID returns a new NetAppAccountId struct
func NewNetAppAccountID(subscriptionId string, resourceGroupName string, netAppAccountName string) NetAppAccountId {
	return NetAppAccountId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NetAppAccountName: netAppAccountName,
	}
}

// ParseNetAppAccountID parses 'input' into a NetAppAccountId
func ParseNetAppAccountID(input string) (*NetAppAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetAppAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetAppAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetAppAccountName, ok = parsed.Parsed["netAppAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "netAppAccountName", *parsed)
	}

	return &id, nil
}

// ParseNetAppAccountIDInsensitively parses 'input' case-insensitively into a NetAppAccountId
// note: this method should only be used for API response data and not user input
func ParseNetAppAccountIDInsensitively(input string) (*NetAppAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetAppAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetAppAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetAppAccountName, ok = parsed.Parsed["netAppAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "netAppAccountName", *parsed)
	}

	return &id, nil
}

// ValidateNetAppAccountID checks that 'input' can be parsed as a Net App Account ID
func ValidateNetAppAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetAppAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Net App Account ID
func (id NetAppAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Net App Account ID
func (id NetAppAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetApp", "Microsoft.NetApp", "Microsoft.NetApp"),
		resourceids.StaticSegment("staticNetAppAccounts", "netAppAccounts", "netAppAccounts"),
		resourceids.UserSpecifiedSegment("netAppAccountName", "netAppAccountValue"),
	}
}

// String returns a human-readable description of this Net App Account ID
func (id NetAppAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Net App Account Name: %q", id.NetAppAccountName),
	}
	return fmt.Sprintf("Net App Account (%s)", strings.Join(components, "\n"))
}
