package iotconnectors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&IotConnectorId{})
}

var _ resourceids.ResourceId = &IotConnectorId{}

// IotConnectorId is a struct representing the Resource ID for a Iot Connector
type IotConnectorId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	IotConnectorName  string
}

// NewIotConnectorID returns a new IotConnectorId struct
func NewIotConnectorID(subscriptionId string, resourceGroupName string, workspaceName string, iotConnectorName string) IotConnectorId {
	return IotConnectorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		IotConnectorName:  iotConnectorName,
	}
}

// ParseIotConnectorID parses 'input' into a IotConnectorId
func ParseIotConnectorID(input string) (*IotConnectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IotConnectorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IotConnectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIotConnectorIDInsensitively parses 'input' case-insensitively into a IotConnectorId
// note: this method should only be used for API response data and not user input
func ParseIotConnectorIDInsensitively(input string) (*IotConnectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IotConnectorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IotConnectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IotConnectorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.WorkspaceName, ok = input.Parsed["workspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", input)
	}

	if id.IotConnectorName, ok = input.Parsed["iotConnectorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "iotConnectorName", input)
	}

	return nil
}

// ValidateIotConnectorID checks that 'input' can be parsed as a Iot Connector ID
func ValidateIotConnectorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIotConnectorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Iot Connector ID
func (id IotConnectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthcareApis/workspaces/%s/iotConnectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.IotConnectorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Iot Connector ID
func (id IotConnectorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHealthcareApis", "Microsoft.HealthcareApis", "Microsoft.HealthcareApis"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticIotConnectors", "iotConnectors", "iotConnectors"),
		resourceids.UserSpecifiedSegment("iotConnectorName", "iotConnectorName"),
	}
}

// String returns a human-readable description of this Iot Connector ID
func (id IotConnectorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Iot Connector Name: %q", id.IotConnectorName),
	}
	return fmt.Sprintf("Iot Connector (%s)", strings.Join(components, "\n"))
}
