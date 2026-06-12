package containerinstance

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RevisionId{})
}

var _ resourceids.ResourceId = &RevisionId{}

// RevisionId is a struct representing the Resource ID for a Revision
type RevisionId struct {
	SubscriptionId            string
	ResourceGroupName         string
	ContainerGroupProfileName string
	RevisionName              string
}

// NewRevisionID returns a new RevisionId struct
func NewRevisionID(subscriptionId string, resourceGroupName string, containerGroupProfileName string, revisionName string) RevisionId {
	return RevisionId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		ContainerGroupProfileName: containerGroupProfileName,
		RevisionName:              revisionName,
	}
}

// ParseRevisionID parses 'input' into a RevisionId
func ParseRevisionID(input string) (*RevisionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RevisionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RevisionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRevisionIDInsensitively parses 'input' case-insensitively into a RevisionId
// note: this method should only be used for API response data and not user input
func ParseRevisionIDInsensitively(input string) (*RevisionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RevisionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RevisionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RevisionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ContainerGroupProfileName, ok = input.Parsed["containerGroupProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerGroupProfileName", input)
	}

	if id.RevisionName, ok = input.Parsed["revisionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "revisionName", input)
	}

	return nil
}

// ValidateRevisionID checks that 'input' can be parsed as a Revision ID
func ValidateRevisionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRevisionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Revision ID
func (id RevisionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerInstance/containerGroupProfiles/%s/revisions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerGroupProfileName, id.RevisionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Revision ID
func (id RevisionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerInstance", "Microsoft.ContainerInstance", "Microsoft.ContainerInstance"),
		resourceids.StaticSegment("staticContainerGroupProfiles", "containerGroupProfiles", "containerGroupProfiles"),
		resourceids.UserSpecifiedSegment("containerGroupProfileName", "containerGroupProfileName"),
		resourceids.StaticSegment("staticRevisions", "revisions", "revisions"),
		resourceids.UserSpecifiedSegment("revisionName", "revisionName"),
	}
}

// String returns a human-readable description of this Revision ID
func (id RevisionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container Group Profile Name: %q", id.ContainerGroupProfileName),
		fmt.Sprintf("Revision Name: %q", id.RevisionName),
	}
	return fmt.Sprintf("Revision (%s)", strings.Join(components, "\n"))
}
