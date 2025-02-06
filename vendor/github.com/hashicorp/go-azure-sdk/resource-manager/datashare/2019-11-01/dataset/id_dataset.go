package dataset

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DataSetId{})
}

var _ resourceids.ResourceId = &DataSetId{}

// DataSetId is a struct representing the Resource ID for a Data Set
type DataSetId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	ShareName         string
	DataSetName       string
}

// NewDataSetID returns a new DataSetId struct
func NewDataSetID(subscriptionId string, resourceGroupName string, accountName string, shareName string, dataSetName string) DataSetId {
	return DataSetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		ShareName:         shareName,
		DataSetName:       dataSetName,
	}
}

// ParseDataSetID parses 'input' into a DataSetId
func ParseDataSetID(input string) (*DataSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataSetIDInsensitively parses 'input' case-insensitively into a DataSetId
// note: this method should only be used for API response data and not user input
func ParseDataSetIDInsensitively(input string) (*DataSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataSetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AccountName, ok = input.Parsed["accountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accountName", input)
	}

	if id.ShareName, ok = input.Parsed["shareName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "shareName", input)
	}

	if id.DataSetName, ok = input.Parsed["dataSetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataSetName", input)
	}

	return nil
}

// ValidateDataSetID checks that 'input' can be parsed as a Data Set ID
func ValidateDataSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Set ID
func (id DataSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataShare/accounts/%s/shares/%s/dataSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName, id.DataSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Set ID
func (id DataSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataShare", "Microsoft.DataShare", "Microsoft.DataShare"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountName"),
		resourceids.StaticSegment("staticShares", "shares", "shares"),
		resourceids.UserSpecifiedSegment("shareName", "shareName"),
		resourceids.StaticSegment("staticDataSets", "dataSets", "dataSets"),
		resourceids.UserSpecifiedSegment("dataSetName", "dataSetName"),
	}
}

// String returns a human-readable description of this Data Set ID
func (id DataSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Share Name: %q", id.ShareName),
		fmt.Sprintf("Data Set Name: %q", id.DataSetName),
	}
	return fmt.Sprintf("Data Set (%s)", strings.Join(components, "\n"))
}
