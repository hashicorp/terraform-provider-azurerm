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
	recaser.RegisterResourceId(&WorkerPoolInstanceId{})
}

var _ resourceids.ResourceId = &WorkerPoolInstanceId{}

// WorkerPoolInstanceId is a struct representing the Resource ID for a Worker Pool Instance
type WorkerPoolInstanceId struct {
	SubscriptionId         string
	ResourceGroupName      string
	HostingEnvironmentName string
	WorkerPoolName         string
	InstanceName           string
}

// NewWorkerPoolInstanceID returns a new WorkerPoolInstanceId struct
func NewWorkerPoolInstanceID(subscriptionId string, resourceGroupName string, hostingEnvironmentName string, workerPoolName string, instanceName string) WorkerPoolInstanceId {
	return WorkerPoolInstanceId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		HostingEnvironmentName: hostingEnvironmentName,
		WorkerPoolName:         workerPoolName,
		InstanceName:           instanceName,
	}
}

// ParseWorkerPoolInstanceID parses 'input' into a WorkerPoolInstanceId
func ParseWorkerPoolInstanceID(input string) (*WorkerPoolInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkerPoolInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkerPoolInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkerPoolInstanceIDInsensitively parses 'input' case-insensitively into a WorkerPoolInstanceId
// note: this method should only be used for API response data and not user input
func ParseWorkerPoolInstanceIDInsensitively(input string) (*WorkerPoolInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkerPoolInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkerPoolInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkerPoolInstanceId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.InstanceName, ok = input.Parsed["instanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instanceName", input)
	}

	return nil
}

// ValidateWorkerPoolInstanceID checks that 'input' can be parsed as a Worker Pool Instance ID
func ValidateWorkerPoolInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkerPoolInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Worker Pool Instance ID
func (id WorkerPoolInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/hostingEnvironments/%s/workerPools/%s/instances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostingEnvironmentName, id.WorkerPoolName, id.InstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Worker Pool Instance ID
func (id WorkerPoolInstanceId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceName", "instanceName"),
	}
}

// String returns a human-readable description of this Worker Pool Instance ID
func (id WorkerPoolInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Hosting Environment Name: %q", id.HostingEnvironmentName),
		fmt.Sprintf("Worker Pool Name: %q", id.WorkerPoolName),
		fmt.Sprintf("Instance Name: %q", id.InstanceName),
	}
	return fmt.Sprintf("Worker Pool Instance (%s)", strings.Join(components, "\n"))
}
