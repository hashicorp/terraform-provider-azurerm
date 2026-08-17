package managedclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MeshRevisionProfileId{})
}

var _ resourceids.ResourceId = &MeshRevisionProfileId{}

// MeshRevisionProfileId is a struct representing the Resource ID for a Mesh Revision Profile
type MeshRevisionProfileId struct {
	SubscriptionId          string
	LocationName            string
	MeshRevisionProfileName string
}

// NewMeshRevisionProfileID returns a new MeshRevisionProfileId struct
func NewMeshRevisionProfileID(subscriptionId string, locationName string, meshRevisionProfileName string) MeshRevisionProfileId {
	return MeshRevisionProfileId{
		SubscriptionId:          subscriptionId,
		LocationName:            locationName,
		MeshRevisionProfileName: meshRevisionProfileName,
	}
}

// ParseMeshRevisionProfileID parses 'input' into a MeshRevisionProfileId
func ParseMeshRevisionProfileID(input string) (*MeshRevisionProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MeshRevisionProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MeshRevisionProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMeshRevisionProfileIDInsensitively parses 'input' case-insensitively into a MeshRevisionProfileId
// note: this method should only be used for API response data and not user input
func ParseMeshRevisionProfileIDInsensitively(input string) (*MeshRevisionProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MeshRevisionProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MeshRevisionProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MeshRevisionProfileId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.MeshRevisionProfileName, ok = input.Parsed["meshRevisionProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "meshRevisionProfileName", input)
	}

	return nil
}

// ValidateMeshRevisionProfileID checks that 'input' can be parsed as a Mesh Revision Profile ID
func ValidateMeshRevisionProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMeshRevisionProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mesh Revision Profile ID
func (id MeshRevisionProfileId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.ContainerService/locations/%s/meshRevisionProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.MeshRevisionProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Mesh Revision Profile ID
func (id MeshRevisionProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticMeshRevisionProfiles", "meshRevisionProfiles", "meshRevisionProfiles"),
		resourceids.UserSpecifiedSegment("meshRevisionProfileName", "meshRevisionProfileName"),
	}
}

// String returns a human-readable description of this Mesh Revision Profile ID
func (id MeshRevisionProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Mesh Revision Profile Name: %q", id.MeshRevisionProfileName),
	}
	return fmt.Sprintf("Mesh Revision Profile (%s)", strings.Join(components, "\n"))
}
