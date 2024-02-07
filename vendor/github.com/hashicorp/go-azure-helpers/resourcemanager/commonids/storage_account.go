// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &StorageAccountId{}

// StorageAccountId is a struct representing the Resource ID for a Storage Account
type StorageAccountId struct {
	SubscriptionId     string
	ResourceGroupName  string
	StorageAccountName string
}

// NewStorageAccountID returns a new StorageAccountId struct
func NewStorageAccountID(subscriptionId string, resourceGroupName string, storageAccountName string) StorageAccountId {
	return StorageAccountId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		StorageAccountName: storageAccountName,
	}
}

// ParseStorageAccountID parses 'input' into a StorageAccountId
func ParseStorageAccountID(input string) (*StorageAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageAccountId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageAccountIDInsensitively parses 'input' case-insensitively into a StorageAccountId
// note: this method should only be used for API response data and not user input
func ParseStorageAccountIDInsensitively(input string) (*StorageAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageAccountId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateStorageAccountID checks that 'input' can be parsed as a Storage Account ID
func ValidateStorageAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Account ID
func (id StorageAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Account ID
func (id StorageAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("storageAccountName", "storageAccountValue"),
	}
}

// String returns a human-readable description of this Storage Account ID
func (id StorageAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Account Name: %q", id.StorageAccountName),
	}
	return fmt.Sprintf("Storage Account (%s)", strings.Join(components, "\n"))
}
