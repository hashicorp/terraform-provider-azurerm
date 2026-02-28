package taskhubs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SchedulerId{})
}

var _ resourceids.ResourceId = &SchedulerId{}

// SchedulerId is a struct representing the Resource ID for a Scheduler
type SchedulerId struct {
	SubscriptionId    string
	ResourceGroupName string
	SchedulerName     string
}

// NewSchedulerID returns a new SchedulerId struct
func NewSchedulerID(subscriptionId string, resourceGroupName string, schedulerName string) SchedulerId {
	return SchedulerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SchedulerName:     schedulerName,
	}
}

// ParseSchedulerID parses 'input' into a SchedulerId
func ParseSchedulerID(input string) (*SchedulerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SchedulerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SchedulerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSchedulerIDInsensitively parses 'input' case-insensitively into a SchedulerId
// note: this method should only be used for API response data and not user input
func ParseSchedulerIDInsensitively(input string) (*SchedulerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SchedulerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SchedulerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SchedulerId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateSchedulerID checks that 'input' can be parsed as a Scheduler ID
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

// ID returns the formatted Scheduler ID
func (id SchedulerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DurableTask/schedulers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scheduler ID
func (id SchedulerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDurableTask", "Microsoft.DurableTask", "Microsoft.DurableTask"),
		resourceids.StaticSegment("staticSchedulers", "schedulers", "schedulers"),
		resourceids.UserSpecifiedSegment("schedulerName", "schedulerName"),
	}
}

// String returns a human-readable description of this Scheduler ID
func (id SchedulerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Scheduler Name: %q", id.SchedulerName),
	}
	return fmt.Sprintf("Scheduler (%s)", strings.Join(components, "\n"))
}
