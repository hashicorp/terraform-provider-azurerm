package managedprivateendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ManagedPrivateEndpointId{}

// ManagedPrivateEndpointId is a struct representing the Resource ID for a Managed Private Endpoint
type ManagedPrivateEndpointId struct {
	BaseURI                    string
	ManagedVirtualNetworkName  string
	ManagedPrivateEndpointName string
}

// NewManagedPrivateEndpointID returns a new ManagedPrivateEndpointId struct
func NewManagedPrivateEndpointID(baseURI string, managedVirtualNetworkName string, managedPrivateEndpointName string) ManagedPrivateEndpointId {
	return ManagedPrivateEndpointId{
		BaseURI:                    strings.TrimSuffix(baseURI, "/"),
		ManagedVirtualNetworkName:  managedVirtualNetworkName,
		ManagedPrivateEndpointName: managedPrivateEndpointName,
	}
}

// ParseManagedPrivateEndpointID parses 'input' into a ManagedPrivateEndpointId
func ParseManagedPrivateEndpointID(input string) (*ManagedPrivateEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedPrivateEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedPrivateEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedPrivateEndpointIDInsensitively parses 'input' case-insensitively into a ManagedPrivateEndpointId
// note: this method should only be used for API response data and not user input
func ParseManagedPrivateEndpointIDInsensitively(input string) (*ManagedPrivateEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedPrivateEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedPrivateEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedPrivateEndpointId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.ManagedVirtualNetworkName, ok = input.Parsed["managedVirtualNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedVirtualNetworkName", input)
	}

	if id.ManagedPrivateEndpointName, ok = input.Parsed["managedPrivateEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedPrivateEndpointName", input)
	}

	return nil
}

// ValidateManagedPrivateEndpointID checks that 'input' can be parsed as a Managed Private Endpoint ID
func ValidateManagedPrivateEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedPrivateEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Private Endpoint ID
func (id ManagedPrivateEndpointId) ID() string {
	fmtString := "%s/managedVirtualNetworks/%s/managedPrivateEndpoints/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.ManagedVirtualNetworkName, id.ManagedPrivateEndpointName)
}

// Path returns the formatted Managed Private Endpoint ID without the BaseURI
func (id ManagedPrivateEndpointId) Path() string {
	fmtString := "/managedVirtualNetworks/%s/managedPrivateEndpoints/%s"
	return fmt.Sprintf(fmtString, id.ManagedVirtualNetworkName, id.ManagedPrivateEndpointName)
}

// PathElements returns the values of Managed Private Endpoint ID Segments without the BaseURI
func (id ManagedPrivateEndpointId) PathElements() []any {
	return []any{id.ManagedVirtualNetworkName, id.ManagedPrivateEndpointName}
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Private Endpoint ID
func (id ManagedPrivateEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticManagedVirtualNetworks", "managedVirtualNetworks", "managedVirtualNetworks"),
		resourceids.UserSpecifiedSegment("managedVirtualNetworkName", "managedVirtualNetworkName"),
		resourceids.StaticSegment("staticManagedPrivateEndpoints", "managedPrivateEndpoints", "managedPrivateEndpoints"),
		resourceids.UserSpecifiedSegment("managedPrivateEndpointName", "managedPrivateEndpointName"),
	}
}

// String returns a human-readable description of this Managed Private Endpoint ID
func (id ManagedPrivateEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Managed Virtual Network Name: %q", id.ManagedVirtualNetworkName),
		fmt.Sprintf("Managed Private Endpoint Name: %q", id.ManagedPrivateEndpointName),
	}
	return fmt.Sprintf("Managed Private Endpoint (%s)", strings.Join(components, "\n"))
}
