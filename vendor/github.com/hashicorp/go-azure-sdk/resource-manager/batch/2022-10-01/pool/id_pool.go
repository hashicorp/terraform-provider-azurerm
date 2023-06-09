package pool

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PoolId{}

// PoolId is a struct representing the Resource ID for a Pool
type PoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	BatchAccountName  string
	PoolName          string
}

// NewPoolID returns a new PoolId struct
func NewPoolID(subscriptionId string, resourceGroupName string, batchAccountName string, poolName string) PoolId {
	return PoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		BatchAccountName:  batchAccountName,
		PoolName:          poolName,
	}
}

// ParsePoolID parses 'input' into a PoolId
func ParsePoolID(input string) (*PoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(PoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.BatchAccountName, ok = parsed.Parsed["batchAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "batchAccountName", *parsed)
	}

	if id.PoolName, ok = parsed.Parsed["poolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "poolName", *parsed)
	}

	return &id, nil
}

// ParsePoolIDInsensitively parses 'input' case-insensitively into a PoolId
// note: this method should only be used for API response data and not user input
func ParsePoolIDInsensitively(input string) (*PoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(PoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.BatchAccountName, ok = parsed.Parsed["batchAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "batchAccountName", *parsed)
	}

	if id.PoolName, ok = parsed.Parsed["poolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "poolName", *parsed)
	}

	return &id, nil
}

// ValidatePoolID checks that 'input' can be parsed as a Pool ID
func ValidatePoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Pool ID
func (id PoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Batch/batchAccounts/%s/pools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BatchAccountName, id.PoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Pool ID
func (id PoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBatch", "Microsoft.Batch", "Microsoft.Batch"),
		resourceids.StaticSegment("staticBatchAccounts", "batchAccounts", "batchAccounts"),
		resourceids.UserSpecifiedSegment("batchAccountName", "batchAccountValue"),
		resourceids.StaticSegment("staticPools", "pools", "pools"),
		resourceids.UserSpecifiedSegment("poolName", "poolValue"),
	}
}

// String returns a human-readable description of this Pool ID
func (id PoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Batch Account Name: %q", id.BatchAccountName),
		fmt.Sprintf("Pool Name: %q", id.PoolName),
	}
	return fmt.Sprintf("Pool (%s)", strings.Join(components, "\n"))
}
