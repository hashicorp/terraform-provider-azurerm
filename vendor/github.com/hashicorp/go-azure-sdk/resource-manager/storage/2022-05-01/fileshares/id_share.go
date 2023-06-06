package fileshares

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ShareId{}

// ShareId is a struct representing the Resource ID for a Share
type ShareId struct {
	SubscriptionId     string
	ResourceGroupName  string
	StorageAccountName string
	ShareName          string
}

// NewShareID returns a new ShareId struct
func NewShareID(subscriptionId string, resourceGroupName string, storageAccountName string, shareName string) ShareId {
	return ShareId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		StorageAccountName: storageAccountName,
		ShareName:          shareName,
	}
}

// ParseShareID parses 'input' into a ShareId
func ParseShareID(input string) (*ShareId, error) {
	parser := resourceids.NewParserFromResourceIdType(ShareId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ShareId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageAccountName, ok = parsed.Parsed["storageAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageAccountName", *parsed)
	}

	if id.ShareName, ok = parsed.Parsed["shareName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "shareName", *parsed)
	}

	return &id, nil
}

// ParseShareIDInsensitively parses 'input' case-insensitively into a ShareId
// note: this method should only be used for API response data and not user input
func ParseShareIDInsensitively(input string) (*ShareId, error) {
	parser := resourceids.NewParserFromResourceIdType(ShareId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ShareId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageAccountName, ok = parsed.Parsed["storageAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageAccountName", *parsed)
	}

	if id.ShareName, ok = parsed.Parsed["shareName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "shareName", *parsed)
	}

	return &id, nil
}

// ValidateShareID checks that 'input' can be parsed as a Share ID
func ValidateShareID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseShareID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Share ID
func (id ShareId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/default/shares/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, id.ShareName)
}

// Segments returns a slice of Resource ID Segments which comprise this Share ID
func (id ShareId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("storageAccountName", "storageAccountValue"),
		resourceids.StaticSegment("staticFileServices", "fileServices", "fileServices"),
		resourceids.StaticSegment("staticDefault", "default", "default"),
		resourceids.StaticSegment("staticShares", "shares", "shares"),
		resourceids.UserSpecifiedSegment("shareName", "shareValue"),
	}
}

// String returns a human-readable description of this Share ID
func (id ShareId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Account Name: %q", id.StorageAccountName),
		fmt.Sprintf("Share Name: %q", id.ShareName),
	}
	return fmt.Sprintf("Share (%s)", strings.Join(components, "\n"))
}
