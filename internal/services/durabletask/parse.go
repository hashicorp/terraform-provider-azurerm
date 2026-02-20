// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// RetentionPolicyID represents the unique identifier for a Durable Task Scheduler Retention Policy
// Note: The retention policy doesn't have a separate resource ID in Azure - it uses the scheduler ID
// with a "default" retention policy name, but we create a synthetic ID for Terraform state management
type RetentionPolicyID struct {
	SubscriptionId    string
	ResourceGroupName string
	SchedulerName     string
}

var _ resourceids.ResourceId = &RetentionPolicyID{}

func NewRetentionPolicyID(subscriptionId, resourceGroupName, schedulerName string) RetentionPolicyID {
	return RetentionPolicyID{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SchedulerName:     schedulerName,
	}
}

func (id RetentionPolicyID) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DurableTask/schedulers/%s/retentionPolicies/default"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)
}

func (id RetentionPolicyID) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Scheduler Name: %q", id.SchedulerName),
	}
	return fmt.Sprintf("Retention Policy (%s)", strings.Join(components, "\n"))
}

func (id RetentionPolicyID) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDurableTask", "Microsoft.DurableTask", "Microsoft.DurableTask"),
		resourceids.StaticSegment("staticSchedulers", "schedulers", "schedulers"),
		resourceids.UserSpecifiedSegment("schedulerName", "schedulerName"),
		resourceids.StaticSegment("staticRetentionPolicies", "retentionPolicies", "retentionPolicies"),
		resourceids.StaticSegment("staticDefault", "default", "default"),
	}
}

func (id *RetentionPolicyID) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SchedulerName, ok = input.Parsed["schedulerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "schedulerName", input)
	}

	return nil
}

func ParseRetentionPolicyID(input string) (*RetentionPolicyID, error) {
	parser := resourceids.NewParserFromResourceIdType(&RetentionPolicyID{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RetentionPolicyID{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func ValidateRetentionPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRetentionPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
