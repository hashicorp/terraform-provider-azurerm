package dicomservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DicomServiceId{})
}

var _ resourceids.ResourceId = &DicomServiceId{}

// DicomServiceId is a struct representing the Resource ID for a Dicom Service
type DicomServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	DicomServiceName  string
}

// NewDicomServiceID returns a new DicomServiceId struct
func NewDicomServiceID(subscriptionId string, resourceGroupName string, workspaceName string, dicomServiceName string) DicomServiceId {
	return DicomServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		DicomServiceName:  dicomServiceName,
	}
}

// ParseDicomServiceID parses 'input' into a DicomServiceId
func ParseDicomServiceID(input string) (*DicomServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DicomServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DicomServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDicomServiceIDInsensitively parses 'input' case-insensitively into a DicomServiceId
// note: this method should only be used for API response data and not user input
func ParseDicomServiceIDInsensitively(input string) (*DicomServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DicomServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DicomServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DicomServiceId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DicomServiceName, ok = input.Parsed["dicomServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dicomServiceName", input)
	}

	return nil
}

// ValidateDicomServiceID checks that 'input' can be parsed as a Dicom Service ID
func ValidateDicomServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDicomServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dicom Service ID
func (id DicomServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthcareApis/workspaces/%s/dicomServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.DicomServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dicom Service ID
func (id DicomServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHealthcareApis", "Microsoft.HealthcareApis", "Microsoft.HealthcareApis"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticDicomServices", "dicomServices", "dicomServices"),
		resourceids.UserSpecifiedSegment("dicomServiceName", "dicomServiceName"),
	}
}

// String returns a human-readable description of this Dicom Service ID
func (id DicomServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Dicom Service Name: %q", id.DicomServiceName),
	}
	return fmt.Sprintf("Dicom Service (%s)", strings.Join(components, "\n"))
}
