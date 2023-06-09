package networkmanagerconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = Providers2NetworkManagerConnectionId{}

// Providers2NetworkManagerConnectionId is a struct representing the Resource ID for a Providers 2 Network Manager Connection
type Providers2NetworkManagerConnectionId struct {
	ManagementGroupId            string
	NetworkManagerConnectionName string
}

// NewProviders2NetworkManagerConnectionID returns a new Providers2NetworkManagerConnectionId struct
func NewProviders2NetworkManagerConnectionID(managementGroupId string, networkManagerConnectionName string) Providers2NetworkManagerConnectionId {
	return Providers2NetworkManagerConnectionId{
		ManagementGroupId:            managementGroupId,
		NetworkManagerConnectionName: networkManagerConnectionName,
	}
}

// ParseProviders2NetworkManagerConnectionID parses 'input' into a Providers2NetworkManagerConnectionId
func ParseProviders2NetworkManagerConnectionID(input string) (*Providers2NetworkManagerConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(Providers2NetworkManagerConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := Providers2NetworkManagerConnectionId{}

	if id.ManagementGroupId, ok = parsed.Parsed["managementGroupId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managementGroupId", *parsed)
	}

	if id.NetworkManagerConnectionName, ok = parsed.Parsed["networkManagerConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerConnectionName", *parsed)
	}

	return &id, nil
}

// ParseProviders2NetworkManagerConnectionIDInsensitively parses 'input' case-insensitively into a Providers2NetworkManagerConnectionId
// note: this method should only be used for API response data and not user input
func ParseProviders2NetworkManagerConnectionIDInsensitively(input string) (*Providers2NetworkManagerConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(Providers2NetworkManagerConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := Providers2NetworkManagerConnectionId{}

	if id.ManagementGroupId, ok = parsed.Parsed["managementGroupId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managementGroupId", *parsed)
	}

	if id.NetworkManagerConnectionName, ok = parsed.Parsed["networkManagerConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerConnectionName", *parsed)
	}

	return &id, nil
}

// ValidateProviders2NetworkManagerConnectionID checks that 'input' can be parsed as a Providers 2 Network Manager Connection ID
func ValidateProviders2NetworkManagerConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviders2NetworkManagerConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Providers 2 Network Manager Connection ID
func (id Providers2NetworkManagerConnectionId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Network/networkManagerConnections/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupId, id.NetworkManagerConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Providers 2 Network Manager Connection ID
func (id Providers2NetworkManagerConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagement", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("staticManagementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("managementGroupId", "managementGroupIdValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagerConnections", "networkManagerConnections", "networkManagerConnections"),
		resourceids.UserSpecifiedSegment("networkManagerConnectionName", "networkManagerConnectionValue"),
	}
}

// String returns a human-readable description of this Providers 2 Network Manager Connection ID
func (id Providers2NetworkManagerConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Management Group: %q", id.ManagementGroupId),
		fmt.Sprintf("Network Manager Connection Name: %q", id.NetworkManagerConnectionName),
	}
	return fmt.Sprintf("Providers 2 Network Manager Connection (%s)", strings.Join(components, "\n"))
}
