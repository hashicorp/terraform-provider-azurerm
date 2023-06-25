package pool

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = BatchAccountId{}

// BatchAccountId is a struct representing the Resource ID for a Batch Account
type BatchAccountId struct {
	SubscriptionId    string
	ResourceGroupName string
	BatchAccountName  string
}

// NewBatchAccountID returns a new BatchAccountId struct
func NewBatchAccountID(subscriptionId string, resourceGroupName string, batchAccountName string) BatchAccountId {
	return BatchAccountId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		BatchAccountName:  batchAccountName,
	}
}

// ParseBatchAccountID parses 'input' into a BatchAccountId
func ParseBatchAccountID(input string) (*BatchAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(BatchAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BatchAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.BatchAccountName, ok = parsed.Parsed["batchAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "batchAccountName", *parsed)
	}

	return &id, nil
}

// ParseBatchAccountIDInsensitively parses 'input' case-insensitively into a BatchAccountId
// note: this method should only be used for API response data and not user input
func ParseBatchAccountIDInsensitively(input string) (*BatchAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(BatchAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BatchAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.BatchAccountName, ok = parsed.Parsed["batchAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "batchAccountName", *parsed)
	}

	return &id, nil
}

// ValidateBatchAccountID checks that 'input' can be parsed as a Batch Account ID
func ValidateBatchAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBatchAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Batch Account ID
func (id BatchAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Batch/batchAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BatchAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Batch Account ID
func (id BatchAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBatch", "Microsoft.Batch", "Microsoft.Batch"),
		resourceids.StaticSegment("staticBatchAccounts", "batchAccounts", "batchAccounts"),
		resourceids.UserSpecifiedSegment("batchAccountName", "batchAccountValue"),
	}
}

// String returns a human-readable description of this Batch Account ID
func (id BatchAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Batch Account Name: %q", id.BatchAccountName),
	}
	return fmt.Sprintf("Batch Account (%s)", strings.Join(components, "\n"))
}
