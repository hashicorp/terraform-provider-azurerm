package fhirservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FhirServiceId{}

// FhirServiceId is a struct representing the Resource ID for a Fhir Service
type FhirServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	FhirServiceName   string
}

// NewFhirServiceID returns a new FhirServiceId struct
func NewFhirServiceID(subscriptionId string, resourceGroupName string, workspaceName string, fhirServiceName string) FhirServiceId {
	return FhirServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		FhirServiceName:   fhirServiceName,
	}
}

// ParseFhirServiceID parses 'input' into a FhirServiceId
func ParseFhirServiceID(input string) (*FhirServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(FhirServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FhirServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.FhirServiceName, ok = parsed.Parsed["fhirServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fhirServiceName", *parsed)
	}

	return &id, nil
}

// ParseFhirServiceIDInsensitively parses 'input' case-insensitively into a FhirServiceId
// note: this method should only be used for API response data and not user input
func ParseFhirServiceIDInsensitively(input string) (*FhirServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(FhirServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FhirServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.FhirServiceName, ok = parsed.Parsed["fhirServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fhirServiceName", *parsed)
	}

	return &id, nil
}

// ValidateFhirServiceID checks that 'input' can be parsed as a Fhir Service ID
func ValidateFhirServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFhirServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Fhir Service ID
func (id FhirServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthcareApis/workspaces/%s/fhirServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.FhirServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Fhir Service ID
func (id FhirServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHealthcareApis", "Microsoft.HealthcareApis", "Microsoft.HealthcareApis"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticFhirServices", "fhirServices", "fhirServices"),
		resourceids.UserSpecifiedSegment("fhirServiceName", "fhirServiceValue"),
	}
}

// String returns a human-readable description of this Fhir Service ID
func (id FhirServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Fhir Service Name: %q", id.FhirServiceName),
	}
	return fmt.Sprintf("Fhir Service (%s)", strings.Join(components, "\n"))
}
