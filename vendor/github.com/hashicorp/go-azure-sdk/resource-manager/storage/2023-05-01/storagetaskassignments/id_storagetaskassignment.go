package storagetaskassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StorageTaskAssignmentId{})
}

var _ resourceids.ResourceId = &StorageTaskAssignmentId{}

// StorageTaskAssignmentId is a struct representing the Resource ID for a Storage Task Assignment
type StorageTaskAssignmentId struct {
	SubscriptionId            string
	ResourceGroupName         string
	StorageAccountName        string
	StorageTaskAssignmentName string
}

// NewStorageTaskAssignmentID returns a new StorageTaskAssignmentId struct
func NewStorageTaskAssignmentID(subscriptionId string, resourceGroupName string, storageAccountName string, storageTaskAssignmentName string) StorageTaskAssignmentId {
	return StorageTaskAssignmentId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		StorageAccountName:        storageAccountName,
		StorageTaskAssignmentName: storageTaskAssignmentName,
	}
}

// ParseStorageTaskAssignmentID parses 'input' into a StorageTaskAssignmentId
func ParseStorageTaskAssignmentID(input string) (*StorageTaskAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageTaskAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageTaskAssignmentId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageTaskAssignmentIDInsensitively parses 'input' case-insensitively into a StorageTaskAssignmentId
// note: this method should only be used for API response data and not user input
func ParseStorageTaskAssignmentIDInsensitively(input string) (*StorageTaskAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageTaskAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageTaskAssignmentId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageTaskAssignmentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.StorageTaskAssignmentName, ok = input.Parsed["storageTaskAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageTaskAssignmentName", input)
	}

	return nil
}

// ValidateStorageTaskAssignmentID checks that 'input' can be parsed as a Storage Task Assignment ID
func ValidateStorageTaskAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageTaskAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Task Assignment ID
func (id StorageTaskAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/storageTaskAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, id.StorageTaskAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Task Assignment ID
func (id StorageTaskAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("storageAccountName", "storageAccountValue"),
		resourceids.StaticSegment("staticStorageTaskAssignments", "storageTaskAssignments", "storageTaskAssignments"),
		resourceids.UserSpecifiedSegment("storageTaskAssignmentName", "storageTaskAssignmentValue"),
	}
}

// String returns a human-readable description of this Storage Task Assignment ID
func (id StorageTaskAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Account Name: %q", id.StorageAccountName),
		fmt.Sprintf("Storage Task Assignment Name: %q", id.StorageTaskAssignmentName),
	}
	return fmt.Sprintf("Storage Task Assignment (%s)", strings.Join(components, "\n"))
}
