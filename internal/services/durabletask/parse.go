// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SchedulerID struct {
	SubscriptionId    string
	ResourceGroupName string
	SchedulerName     string
}

func NewSchedulerID(subscriptionId, resourceGroupName, schedulerName string) SchedulerID {
	return SchedulerID{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SchedulerName:     schedulerName,
	}
}

func (id SchedulerID) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DurableTask/schedulers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)
}

func ParseSchedulerID(input string) (*SchedulerID, error) {
	parser := resourceids.NewParserFromResourceIdType(&SchedulerID{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SchedulerID{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SchedulerID) FromParseResult(input resourceids.ParseResult) error {
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

func ValidateSchedulerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSchedulerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

type TaskHubID struct {
	SubscriptionId    string
	ResourceGroupName string
	SchedulerName     string
	TaskHubName       string
}

func NewTaskHubID(subscriptionId, resourceGroupName, schedulerName, taskHubName string) TaskHubID {
	return TaskHubID{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SchedulerName:     schedulerName,
		TaskHubName:       taskHubName,
	}
}

func (id TaskHubID) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DurableTask/schedulers/%s/taskHubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SchedulerName, id.TaskHubName)
}

func ParseTaskHubID(input string) (*TaskHubID, error) {
	parser := resourceids.NewParserFromResourceIdType(&TaskHubID{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TaskHubID{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TaskHubID) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TaskHubName, ok = input.Parsed["taskHubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "taskHubName", input)
	}

	return nil
}

func ValidateTaskHubID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTaskHubID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

type RetentionPolicyID struct {
	SubscriptionId    string
	ResourceGroupName string
	SchedulerName     string
}

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

func ParseRetentionPolicyID(input string) (*RetentionPolicyID, error) {
	parts := strings.Split(input, "/")
	if len(parts) < 11 {
		return nil, fmt.Errorf("invalid retention policy ID format")
	}

	id := RetentionPolicyID{
		SubscriptionId:    parts[2],
		ResourceGroupName: parts[4],
		SchedulerName:     parts[8],
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
