package deviceupdates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PrivateEndpointConnectionProxyId{}

// PrivateEndpointConnectionProxyId is a struct representing the Resource ID for a Private Endpoint Connection Proxy
type PrivateEndpointConnectionProxyId struct {
	SubscriptionId                   string
	ResourceGroupName                string
	AccountName                      string
	PrivateEndpointConnectionProxyId string
}

// NewPrivateEndpointConnectionProxyID returns a new PrivateEndpointConnectionProxyId struct
func NewPrivateEndpointConnectionProxyID(subscriptionId string, resourceGroupName string, accountName string, privateEndpointConnectionProxyId string) PrivateEndpointConnectionProxyId {
	return PrivateEndpointConnectionProxyId{
		SubscriptionId:                   subscriptionId,
		ResourceGroupName:                resourceGroupName,
		AccountName:                      accountName,
		PrivateEndpointConnectionProxyId: privateEndpointConnectionProxyId,
	}
}

// ParsePrivateEndpointConnectionProxyID parses 'input' into a PrivateEndpointConnectionProxyId
func ParsePrivateEndpointConnectionProxyID(input string) (*PrivateEndpointConnectionProxyId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateEndpointConnectionProxyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateEndpointConnectionProxyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accountName", *parsed)
	}

	if id.PrivateEndpointConnectionProxyId, ok = parsed.Parsed["privateEndpointConnectionProxyId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointConnectionProxyId", *parsed)
	}

	return &id, nil
}

// ParsePrivateEndpointConnectionProxyIDInsensitively parses 'input' case-insensitively into a PrivateEndpointConnectionProxyId
// note: this method should only be used for API response data and not user input
func ParsePrivateEndpointConnectionProxyIDInsensitively(input string) (*PrivateEndpointConnectionProxyId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateEndpointConnectionProxyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateEndpointConnectionProxyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accountName", *parsed)
	}

	if id.PrivateEndpointConnectionProxyId, ok = parsed.Parsed["privateEndpointConnectionProxyId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointConnectionProxyId", *parsed)
	}

	return &id, nil
}

// ValidatePrivateEndpointConnectionProxyID checks that 'input' can be parsed as a Private Endpoint Connection Proxy ID
func ValidatePrivateEndpointConnectionProxyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateEndpointConnectionProxyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Endpoint Connection Proxy ID
func (id PrivateEndpointConnectionProxyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DeviceUpdate/accounts/%s/privateEndpointConnectionProxies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.PrivateEndpointConnectionProxyId)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Endpoint Connection Proxy ID
func (id PrivateEndpointConnectionProxyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDeviceUpdate", "Microsoft.DeviceUpdate", "Microsoft.DeviceUpdate"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticPrivateEndpointConnectionProxies", "privateEndpointConnectionProxies", "privateEndpointConnectionProxies"),
		resourceids.UserSpecifiedSegment("privateEndpointConnectionProxyId", "privateEndpointConnectionProxyIdValue"),
	}
}

// String returns a human-readable description of this Private Endpoint Connection Proxy ID
func (id PrivateEndpointConnectionProxyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Private Endpoint Connection Proxy: %q", id.PrivateEndpointConnectionProxyId),
	}
	return fmt.Sprintf("Private Endpoint Connection Proxy (%s)", strings.Join(components, "\n"))
}
