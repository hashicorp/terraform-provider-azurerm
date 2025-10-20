package dataflow

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DataflowId{})
}

var _ resourceids.ResourceId = &DataflowId{}

// DataflowId is a struct representing the Resource ID for a Dataflow
type DataflowId struct {
	SubscriptionId      string
	ResourceGroupName   string
	InstanceName        string
	DataflowProfileName string
	DataflowName        string
}

// NewDataflowID returns a new DataflowId struct
func NewDataflowID(subscriptionId string, resourceGroupName string, instanceName string, dataflowProfileName string, dataflowName string) DataflowId {
	return DataflowId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		InstanceName:        instanceName,
		DataflowProfileName: dataflowProfileName,
		DataflowName:        dataflowName,
	}
}

// ParseDataflowID parses 'input' into a DataflowId
func ParseDataflowID(input string) (*DataflowId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataflowId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataflowId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataflowIDInsensitively parses 'input' case-insensitively into a DataflowId
// note: this method should only be used for API response data and not user input
func ParseDataflowIDInsensitively(input string) (*DataflowId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataflowId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataflowId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataflowId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DataflowProfileName, ok = input.Parsed["dataflowProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataflowProfileName", input)
	}

	if id.DataflowName, ok = input.Parsed["dataflowName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataflowName", input)
	}

	return nil
}

// ValidateDataflowID checks that 'input' can be parsed as a Dataflow ID
func ValidateDataflowID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataflowID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dataflow ID
func (id DataflowId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTOperations/instances/%s/dataflowProfiles/%s/dataflows/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.InstanceName, id.DataflowProfileName, id.DataflowName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dataflow ID
func (id DataflowId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftIoTOperations", "Microsoft.IoTOperations", "Microsoft.IoTOperations"),
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceName", "instanceName"),
		resourceids.StaticSegment("staticDataflowProfiles", "dataflowProfiles", "dataflowProfiles"),
		resourceids.UserSpecifiedSegment("dataflowProfileName", "dataflowProfileName"),
		resourceids.StaticSegment("staticDataflows", "dataflows", "dataflows"),
		resourceids.UserSpecifiedSegment("dataflowName", "dataflowName"),
	}
}

// String returns a human-readable description of this Dataflow ID
func (id DataflowId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Instance Name: %q", id.InstanceName),
		fmt.Sprintf("Dataflow Profile Name: %q", id.DataflowProfileName),
		fmt.Sprintf("Dataflow Name: %q", id.DataflowName),
	}
	return fmt.Sprintf("Dataflow (%s)", strings.Join(components, "\n"))
}
