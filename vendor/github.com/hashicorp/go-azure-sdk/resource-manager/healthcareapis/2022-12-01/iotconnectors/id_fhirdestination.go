package iotconnectors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FhirDestinationId{}

// FhirDestinationId is a struct representing the Resource ID for a Fhir Destination
type FhirDestinationId struct {
	SubscriptionId      string
	ResourceGroupName   string
	WorkspaceName       string
	IotConnectorName    string
	FhirDestinationName string
}

// NewFhirDestinationID returns a new FhirDestinationId struct
func NewFhirDestinationID(subscriptionId string, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string) FhirDestinationId {
	return FhirDestinationId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		WorkspaceName:       workspaceName,
		IotConnectorName:    iotConnectorName,
		FhirDestinationName: fhirDestinationName,
	}
}

// ParseFhirDestinationID parses 'input' into a FhirDestinationId
func ParseFhirDestinationID(input string) (*FhirDestinationId, error) {
	parser := resourceids.NewParserFromResourceIdType(FhirDestinationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FhirDestinationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.IotConnectorName, ok = parsed.Parsed["iotConnectorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "iotConnectorName", *parsed)
	}

	if id.FhirDestinationName, ok = parsed.Parsed["fhirDestinationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fhirDestinationName", *parsed)
	}

	return &id, nil
}

// ParseFhirDestinationIDInsensitively parses 'input' case-insensitively into a FhirDestinationId
// note: this method should only be used for API response data and not user input
func ParseFhirDestinationIDInsensitively(input string) (*FhirDestinationId, error) {
	parser := resourceids.NewParserFromResourceIdType(FhirDestinationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FhirDestinationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.IotConnectorName, ok = parsed.Parsed["iotConnectorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "iotConnectorName", *parsed)
	}

	if id.FhirDestinationName, ok = parsed.Parsed["fhirDestinationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fhirDestinationName", *parsed)
	}

	return &id, nil
}

// ValidateFhirDestinationID checks that 'input' can be parsed as a Fhir Destination ID
func ValidateFhirDestinationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFhirDestinationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Fhir Destination ID
func (id FhirDestinationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthcareApis/workspaces/%s/iotConnectors/%s/fhirDestinations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.IotConnectorName, id.FhirDestinationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Fhir Destination ID
func (id FhirDestinationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHealthcareApis", "Microsoft.HealthcareApis", "Microsoft.HealthcareApis"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticIotConnectors", "iotConnectors", "iotConnectors"),
		resourceids.UserSpecifiedSegment("iotConnectorName", "iotConnectorValue"),
		resourceids.StaticSegment("staticFhirDestinations", "fhirDestinations", "fhirDestinations"),
		resourceids.UserSpecifiedSegment("fhirDestinationName", "fhirDestinationValue"),
	}
}

// String returns a human-readable description of this Fhir Destination ID
func (id FhirDestinationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Iot Connector Name: %q", id.IotConnectorName),
		fmt.Sprintf("Fhir Destination Name: %q", id.FhirDestinationName),
	}
	return fmt.Sprintf("Fhir Destination (%s)", strings.Join(components, "\n"))
}
