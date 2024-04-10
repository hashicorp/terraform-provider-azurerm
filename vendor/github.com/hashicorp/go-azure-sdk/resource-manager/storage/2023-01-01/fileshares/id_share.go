package fileshares

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ShareId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ShareId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ShareId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseShareIDInsensitively parses 'input' case-insensitively into a ShareId
// note: this method should only be used for API response data and not user input
func ParseShareIDInsensitively(input string) (*ShareId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ShareId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ShareId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ShareId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageAccountName, ok = input.Parsed["storageAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageAccountName", input)
	}

	if id.ShareName, ok = input.Parsed["shareName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "shareName", input)
	}

	return nil
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
