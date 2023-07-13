package share

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProviderShareSubscriptionId{}

// ProviderShareSubscriptionId is a struct representing the Resource ID for a Provider Share Subscription
type ProviderShareSubscriptionId struct {
	SubscriptionId              string
	ResourceGroupName           string
	AccountName                 string
	ShareName                   string
	ProviderShareSubscriptionId string
}

// NewProviderShareSubscriptionID returns a new ProviderShareSubscriptionId struct
func NewProviderShareSubscriptionID(subscriptionId string, resourceGroupName string, accountName string, shareName string, providerShareSubscriptionId string) ProviderShareSubscriptionId {
	return ProviderShareSubscriptionId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		AccountName:                 accountName,
		ShareName:                   shareName,
		ProviderShareSubscriptionId: providerShareSubscriptionId,
	}
}

// ParseProviderShareSubscriptionID parses 'input' into a ProviderShareSubscriptionId
func ParseProviderShareSubscriptionID(input string) (*ProviderShareSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderShareSubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderShareSubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accountName", *parsed)
	}

	if id.ShareName, ok = parsed.Parsed["shareName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "shareName", *parsed)
	}

	if id.ProviderShareSubscriptionId, ok = parsed.Parsed["providerShareSubscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerShareSubscriptionId", *parsed)
	}

	return &id, nil
}

// ParseProviderShareSubscriptionIDInsensitively parses 'input' case-insensitively into a ProviderShareSubscriptionId
// note: this method should only be used for API response data and not user input
func ParseProviderShareSubscriptionIDInsensitively(input string) (*ProviderShareSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderShareSubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderShareSubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accountName", *parsed)
	}

	if id.ShareName, ok = parsed.Parsed["shareName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "shareName", *parsed)
	}

	if id.ProviderShareSubscriptionId, ok = parsed.Parsed["providerShareSubscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerShareSubscriptionId", *parsed)
	}

	return &id, nil
}

// ValidateProviderShareSubscriptionID checks that 'input' can be parsed as a Provider Share Subscription ID
func ValidateProviderShareSubscriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderShareSubscriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Share Subscription ID
func (id ProviderShareSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataShare/accounts/%s/shares/%s/providerShareSubscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName, id.ProviderShareSubscriptionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Share Subscription ID
func (id ProviderShareSubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataShare", "Microsoft.DataShare", "Microsoft.DataShare"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticShares", "shares", "shares"),
		resourceids.UserSpecifiedSegment("shareName", "shareValue"),
		resourceids.StaticSegment("staticProviderShareSubscriptions", "providerShareSubscriptions", "providerShareSubscriptions"),
		resourceids.UserSpecifiedSegment("providerShareSubscriptionId", "providerShareSubscriptionIdValue"),
	}
}

// String returns a human-readable description of this Provider Share Subscription ID
func (id ProviderShareSubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Share Name: %q", id.ShareName),
		fmt.Sprintf("Provider Share Subscription: %q", id.ProviderShareSubscriptionId),
	}
	return fmt.Sprintf("Provider Share Subscription (%s)", strings.Join(components, "\n"))
}
