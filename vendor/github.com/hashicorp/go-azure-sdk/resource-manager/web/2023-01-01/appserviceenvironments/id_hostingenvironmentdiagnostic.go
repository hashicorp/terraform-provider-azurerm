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
	recaser.RegisterResourceId(&HostingEnvironmentDiagnosticId{})
}

var _ resourceids.ResourceId = &HostingEnvironmentDiagnosticId{}

// HostingEnvironmentDiagnosticId is a struct representing the Resource ID for a Hosting Environment Diagnostic
type HostingEnvironmentDiagnosticId struct {
	SubscriptionId         string
	ResourceGroupName      string
	HostingEnvironmentName string
	DiagnosticName         string
}

// NewHostingEnvironmentDiagnosticID returns a new HostingEnvironmentDiagnosticId struct
func NewHostingEnvironmentDiagnosticID(subscriptionId string, resourceGroupName string, hostingEnvironmentName string, diagnosticName string) HostingEnvironmentDiagnosticId {
	return HostingEnvironmentDiagnosticId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		HostingEnvironmentName: hostingEnvironmentName,
		DiagnosticName:         diagnosticName,
	}
}

// ParseHostingEnvironmentDiagnosticID parses 'input' into a HostingEnvironmentDiagnosticId
func ParseHostingEnvironmentDiagnosticID(input string) (*HostingEnvironmentDiagnosticId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HostingEnvironmentDiagnosticId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HostingEnvironmentDiagnosticId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHostingEnvironmentDiagnosticIDInsensitively parses 'input' case-insensitively into a HostingEnvironmentDiagnosticId
// note: this method should only be used for API response data and not user input
func ParseHostingEnvironmentDiagnosticIDInsensitively(input string) (*HostingEnvironmentDiagnosticId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HostingEnvironmentDiagnosticId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HostingEnvironmentDiagnosticId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HostingEnvironmentDiagnosticId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DiagnosticName, ok = input.Parsed["diagnosticName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "diagnosticName", input)
	}

	return nil
}

// ValidateHostingEnvironmentDiagnosticID checks that 'input' can be parsed as a Hosting Environment Diagnostic ID
func ValidateHostingEnvironmentDiagnosticID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHostingEnvironmentDiagnosticID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hosting Environment Diagnostic ID
func (id HostingEnvironmentDiagnosticId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/hostingEnvironments/%s/diagnostics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostingEnvironmentName, id.DiagnosticName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hosting Environment Diagnostic ID
func (id HostingEnvironmentDiagnosticId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticHostingEnvironments", "hostingEnvironments", "hostingEnvironments"),
		resourceids.UserSpecifiedSegment("hostingEnvironmentName", "hostingEnvironmentName"),
		resourceids.StaticSegment("staticDiagnostics", "diagnostics", "diagnostics"),
		resourceids.UserSpecifiedSegment("diagnosticName", "diagnosticName"),
	}
}

// String returns a human-readable description of this Hosting Environment Diagnostic ID
func (id HostingEnvironmentDiagnosticId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Hosting Environment Name: %q", id.HostingEnvironmentName),
		fmt.Sprintf("Diagnostic Name: %q", id.DiagnosticName),
	}
	return fmt.Sprintf("Hosting Environment Diagnostic (%s)", strings.Join(components, "\n"))
}
