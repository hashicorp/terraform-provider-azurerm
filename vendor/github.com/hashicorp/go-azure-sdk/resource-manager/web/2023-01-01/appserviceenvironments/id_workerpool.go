package appserviceenvironments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkerPoolId{})
}

var _ resourceids.ResourceId = &WorkerPoolId{}

// WorkerPoolId is a struct representing the Resource ID for a Worker Pool
type WorkerPoolId struct {
	SubscriptionId         string
	ResourceGroupName      string
	HostingEnvironmentName string
	WorkerPoolName         string
}

// NewWorkerPoolID returns a new WorkerPoolId struct
func NewWorkerPoolID(subscriptionId string, resourceGroupName string, hostingEnvironmentName string, workerPoolName string) WorkerPoolId {
	return WorkerPoolId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		HostingEnvironmentName: hostingEnvironmentName,
		WorkerPoolName:         workerPoolName,
	}
}

// ParseWorkerPoolID parses 'input' into a WorkerPoolId
func ParseWorkerPoolID(input string) (*WorkerPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkerPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkerPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkerPoolIDInsensitively parses 'input' case-insensitively into a WorkerPoolId
// note: this method should only be used for API response data and not user input
func ParseWorkerPoolIDInsensitively(input string) (*WorkerPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkerPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkerPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkerPoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.HostingEnvironmentName, ok = input.Parsed["hostingEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostingEnvironmentName", input)
	}

	if id.WorkerPoolName, ok = input.Parsed["workerPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workerPoolName", input)
	}

	return nil
}

// ValidateWorkerPoolID checks that 'input' can be parsed as a Worker Pool ID
func ValidateWorkerPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkerPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Worker Pool ID
func (id WorkerPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/hostingEnvironments/%s/workerPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostingEnvironmentName, id.WorkerPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Worker Pool ID
func (id WorkerPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticHostingEnvironments", "hostingEnvironments", "hostingEnvironments"),
		resourceids.UserSpecifiedSegment("hostingEnvironmentName", "hostingEnvironmentName"),
		resourceids.StaticSegment("staticWorkerPools", "workerPools", "workerPools"),
		resourceids.UserSpecifiedSegment("workerPoolName", "workerPoolName"),
	}
}

// String returns a human-readable description of this Worker Pool ID
func (id WorkerPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Hosting Environment Name: %q", id.HostingEnvironmentName),
		fmt.Sprintf("Worker Pool Name: %q", id.WorkerPoolName),
	}
	return fmt.Sprintf("Worker Pool (%s)", strings.Join(components, "\n"))
}
