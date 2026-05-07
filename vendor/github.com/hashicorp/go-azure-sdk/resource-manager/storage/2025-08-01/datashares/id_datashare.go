package datashares

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DataShareId{})
}

var _ resourceids.ResourceId = &DataShareId{}

// DataShareId is a struct representing the Resource ID for a Data Share
type DataShareId struct {
	SubscriptionId     string
	ResourceGroupName  string
	StorageAccountName string
	DataShareName      string
}

// NewDataShareID returns a new DataShareId struct
func NewDataShareID(subscriptionId string, resourceGroupName string, storageAccountName string, dataShareName string) DataShareId {
	return DataShareId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		StorageAccountName: storageAccountName,
		DataShareName:      dataShareName,
	}
}

// ParseDataShareID parses 'input' into a DataShareId
func ParseDataShareID(input string) (*DataShareId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataShareId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataShareId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataShareIDInsensitively parses 'input' case-insensitively into a DataShareId
// note: this method should only be used for API response data and not user input
func ParseDataShareIDInsensitively(input string) (*DataShareId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataShareId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataShareId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataShareId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DataShareName, ok = input.Parsed["dataShareName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataShareName", input)
	}

	return nil
}

// ValidateDataShareID checks that 'input' can be parsed as a Data Share ID
func ValidateDataShareID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataShareID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Share ID
func (id DataShareId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/dataShares/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, id.DataShareName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Share ID
func (id DataShareId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("storageAccountName", "storageAccountName"),
		resourceids.StaticSegment("staticDataShares", "dataShares", "dataShares"),
		resourceids.UserSpecifiedSegment("dataShareName", "dataShareName"),
	}
}

// String returns a human-readable description of this Data Share ID
func (id DataShareId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Account Name: %q", id.StorageAccountName),
		fmt.Sprintf("Data Share Name: %q", id.DataShareName),
	}
	return fmt.Sprintf("Data Share (%s)", strings.Join(components, "\n"))
}
