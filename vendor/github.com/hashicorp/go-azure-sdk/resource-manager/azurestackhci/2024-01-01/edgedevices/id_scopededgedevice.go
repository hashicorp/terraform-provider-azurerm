package edgedevices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedEdgeDeviceId{})
}

var _ resourceids.ResourceId = &ScopedEdgeDeviceId{}

// ScopedEdgeDeviceId is a struct representing the Resource ID for a Scoped Edge Device
type ScopedEdgeDeviceId struct {
	ResourceUri    string
	EdgeDeviceName string
}

// NewScopedEdgeDeviceID returns a new ScopedEdgeDeviceId struct
func NewScopedEdgeDeviceID(resourceUri string, edgeDeviceName string) ScopedEdgeDeviceId {
	return ScopedEdgeDeviceId{
		ResourceUri:    resourceUri,
		EdgeDeviceName: edgeDeviceName,
	}
}

// ParseScopedEdgeDeviceID parses 'input' into a ScopedEdgeDeviceId
func ParseScopedEdgeDeviceID(input string) (*ScopedEdgeDeviceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedEdgeDeviceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedEdgeDeviceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedEdgeDeviceIDInsensitively parses 'input' case-insensitively into a ScopedEdgeDeviceId
// note: this method should only be used for API response data and not user input
func ParseScopedEdgeDeviceIDInsensitively(input string) (*ScopedEdgeDeviceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedEdgeDeviceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedEdgeDeviceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedEdgeDeviceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ResourceUri, ok = input.Parsed["resourceUri"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", input)
	}

	if id.EdgeDeviceName, ok = input.Parsed["edgeDeviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "edgeDeviceName", input)
	}

	return nil
}

// ValidateScopedEdgeDeviceID checks that 'input' can be parsed as a Scoped Edge Device ID
func ValidateScopedEdgeDeviceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedEdgeDeviceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Edge Device ID
func (id ScopedEdgeDeviceId) ID() string {
	fmtString := "/%s/providers/Microsoft.AzureStackHCI/edgeDevices/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceUri, "/"), id.EdgeDeviceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Edge Device ID
func (id ScopedEdgeDeviceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceUri", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticEdgeDevices", "edgeDevices", "edgeDevices"),
		resourceids.UserSpecifiedSegment("edgeDeviceName", "edgeDeviceName"),
	}
}

// String returns a human-readable description of this Scoped Edge Device ID
func (id ScopedEdgeDeviceId) String() string {
	components := []string{
		fmt.Sprintf("Resource Uri: %q", id.ResourceUri),
		fmt.Sprintf("Edge Device Name: %q", id.EdgeDeviceName),
	}
	return fmt.Sprintf("Scoped Edge Device (%s)", strings.Join(components, "\n"))
}
