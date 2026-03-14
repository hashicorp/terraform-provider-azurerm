package dataflowendpoint

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DataflowEndpointId{})
}

var _ resourceids.ResourceId = &DataflowEndpointId{}

// DataflowEndpointId is a struct representing the Resource ID for a Dataflow Endpoint
type DataflowEndpointId struct {
	SubscriptionId       string
	ResourceGroupName    string
	InstanceName         string
	DataflowEndpointName string
}

// NewDataflowEndpointID returns a new DataflowEndpointId struct
func NewDataflowEndpointID(subscriptionId string, resourceGroupName string, instanceName string, dataflowEndpointName string) DataflowEndpointId {
	return DataflowEndpointId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		InstanceName:         instanceName,
		DataflowEndpointName: dataflowEndpointName,
	}
}

// ParseDataflowEndpointID parses 'input' into a DataflowEndpointId
func ParseDataflowEndpointID(input string) (*DataflowEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataflowEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataflowEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataflowEndpointIDInsensitively parses 'input' case-insensitively into a DataflowEndpointId
// note: this method should only be used for API response data and not user input
func ParseDataflowEndpointIDInsensitively(input string) (*DataflowEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataflowEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataflowEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataflowEndpointId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.InstanceName, ok = input.Parsed["instanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instanceName", input)
	}

	if id.DataflowEndpointName, ok = input.Parsed["dataflowEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataflowEndpointName", input)
	}

	return nil
}

// ValidateDataflowEndpointID checks that 'input' can be parsed as a Dataflow Endpoint ID
func ValidateDataflowEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataflowEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dataflow Endpoint ID
func (id DataflowEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTOperations/instances/%s/dataflowEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.InstanceName, id.DataflowEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dataflow Endpoint ID
func (id DataflowEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftIoTOperations", "Microsoft.IoTOperations", "Microsoft.IoTOperations"),
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceName", "instanceName"),
		resourceids.StaticSegment("staticDataflowEndpoints", "dataflowEndpoints", "dataflowEndpoints"),
		resourceids.UserSpecifiedSegment("dataflowEndpointName", "dataflowEndpointName"),
	}
}

// String returns a human-readable description of this Dataflow Endpoint ID
func (id DataflowEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Instance Name: %q", id.InstanceName),
		fmt.Sprintf("Dataflow Endpoint Name: %q", id.DataflowEndpointName),
	}
	return fmt.Sprintf("Dataflow Endpoint (%s)", strings.Join(components, "\n"))
}
