package managedprivateendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ManagedVirtualNetworkId{}

// ManagedVirtualNetworkId is a struct representing the Resource ID for a Managed Virtual Network
type ManagedVirtualNetworkId struct {
	BaseURI                   string
	ManagedVirtualNetworkName string
}

// NewManagedVirtualNetworkID returns a new ManagedVirtualNetworkId struct
func NewManagedVirtualNetworkID(baseURI string, managedVirtualNetworkName string) ManagedVirtualNetworkId {
	return ManagedVirtualNetworkId{
		BaseURI:                   strings.TrimSuffix(baseURI, "/"),
		ManagedVirtualNetworkName: managedVirtualNetworkName,
	}
}

// ParseManagedVirtualNetworkID parses 'input' into a ManagedVirtualNetworkId
func ParseManagedVirtualNetworkID(input string) (*ManagedVirtualNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedVirtualNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedVirtualNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedVirtualNetworkIDInsensitively parses 'input' case-insensitively into a ManagedVirtualNetworkId
// note: this method should only be used for API response data and not user input
func ParseManagedVirtualNetworkIDInsensitively(input string) (*ManagedVirtualNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedVirtualNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedVirtualNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedVirtualNetworkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.ManagedVirtualNetworkName, ok = input.Parsed["managedVirtualNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedVirtualNetworkName", input)
	}

	return nil
}

// ValidateManagedVirtualNetworkID checks that 'input' can be parsed as a Managed Virtual Network ID
func ValidateManagedVirtualNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedVirtualNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Virtual Network ID
func (id ManagedVirtualNetworkId) ID() string {
	fmtString := "%s/managedVirtualNetworks/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.ManagedVirtualNetworkName)
}

// Path returns the formatted Managed Virtual Network ID without the BaseURI
func (id ManagedVirtualNetworkId) Path() string {
	fmtString := "/managedVirtualNetworks/%s"
	return fmt.Sprintf(fmtString, id.ManagedVirtualNetworkName)
}

// PathElements returns the values of Managed Virtual Network ID Segments without the BaseURI
func (id ManagedVirtualNetworkId) PathElements() []any {
	return []any{id.ManagedVirtualNetworkName}
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Virtual Network ID
func (id ManagedVirtualNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticManagedVirtualNetworks", "managedVirtualNetworks", "managedVirtualNetworks"),
		resourceids.UserSpecifiedSegment("managedVirtualNetworkName", "managedVirtualNetworkName"),
	}
}

// String returns a human-readable description of this Managed Virtual Network ID
func (id ManagedVirtualNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Managed Virtual Network Name: %q", id.ManagedVirtualNetworkName),
	}
	return fmt.Sprintf("Managed Virtual Network (%s)", strings.Join(components, "\n"))
}
