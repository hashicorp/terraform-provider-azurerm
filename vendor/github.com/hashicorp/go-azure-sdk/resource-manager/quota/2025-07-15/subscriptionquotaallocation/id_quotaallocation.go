package subscriptionquotaallocation

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&QuotaAllocationId{})
}

var _ resourceids.ResourceId = &QuotaAllocationId{}

// QuotaAllocationId is a struct representing the Resource ID for a Quota Allocation
type QuotaAllocationId struct {
	ManagementGroupId    string
	SubscriptionId       string
	GroupQuotaName       string
	ResourceProviderName string
	QuotaAllocationName  string
}

// NewQuotaAllocationID returns a new QuotaAllocationId struct
func NewQuotaAllocationID(managementGroupId string, subscriptionId string, groupQuotaName string, resourceProviderName string, quotaAllocationName string) QuotaAllocationId {
	return QuotaAllocationId{
		ManagementGroupId:    managementGroupId,
		SubscriptionId:       subscriptionId,
		GroupQuotaName:       groupQuotaName,
		ResourceProviderName: resourceProviderName,
		QuotaAllocationName:  quotaAllocationName,
	}
}

// ParseQuotaAllocationID parses 'input' into a QuotaAllocationId
func ParseQuotaAllocationID(input string) (*QuotaAllocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&QuotaAllocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := QuotaAllocationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseQuotaAllocationIDInsensitively parses 'input' case-insensitively into a QuotaAllocationId
// note: this method should only be used for API response data and not user input
func ParseQuotaAllocationIDInsensitively(input string) (*QuotaAllocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&QuotaAllocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := QuotaAllocationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *QuotaAllocationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ManagementGroupId, ok = input.Parsed["managementGroupId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managementGroupId", input)
	}

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.GroupQuotaName, ok = input.Parsed["groupQuotaName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "groupQuotaName", input)
	}

	if id.ResourceProviderName, ok = input.Parsed["resourceProviderName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceProviderName", input)
	}

	if id.QuotaAllocationName, ok = input.Parsed["quotaAllocationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "quotaAllocationName", input)
	}

	return nil
}

// ValidateQuotaAllocationID checks that 'input' can be parsed as a Quota Allocation ID
func ValidateQuotaAllocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseQuotaAllocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Quota Allocation ID
func (id QuotaAllocationId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/subscriptions/%s/providers/Microsoft.Quota/groupQuotas/%s/resourceProviders/%s/quotaAllocations/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupId, id.SubscriptionId, id.GroupQuotaName, id.ResourceProviderName, id.QuotaAllocationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Quota Allocation ID
func (id QuotaAllocationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagement", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("staticManagementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("managementGroupId", "managementGroupId"),
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftQuota", "Microsoft.Quota", "Microsoft.Quota"),
		resourceids.StaticSegment("staticGroupQuotas", "groupQuotas", "groupQuotas"),
		resourceids.UserSpecifiedSegment("groupQuotaName", "groupQuotaName"),
		resourceids.StaticSegment("staticResourceProviders", "resourceProviders", "resourceProviders"),
		resourceids.UserSpecifiedSegment("resourceProviderName", "resourceProviderName"),
		resourceids.StaticSegment("staticQuotaAllocations", "quotaAllocations", "quotaAllocations"),
		resourceids.UserSpecifiedSegment("quotaAllocationName", "quotaAllocationName"),
	}
}

// String returns a human-readable description of this Quota Allocation ID
func (id QuotaAllocationId) String() string {
	components := []string{
		fmt.Sprintf("Management Group: %q", id.ManagementGroupId),
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Group Quota Name: %q", id.GroupQuotaName),
		fmt.Sprintf("Resource Provider Name: %q", id.ResourceProviderName),
		fmt.Sprintf("Quota Allocation Name: %q", id.QuotaAllocationName),
	}
	return fmt.Sprintf("Quota Allocation (%s)", strings.Join(components, "\n"))
}
