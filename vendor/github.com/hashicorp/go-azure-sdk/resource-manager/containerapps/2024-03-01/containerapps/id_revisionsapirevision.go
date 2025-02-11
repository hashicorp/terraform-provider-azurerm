package containerapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RevisionsApiRevisionId{})
}

var _ resourceids.ResourceId = &RevisionsApiRevisionId{}

// RevisionsApiRevisionId is a struct representing the Resource ID for a Revisions Api Revision
type RevisionsApiRevisionId struct {
	SubscriptionId    string
	ResourceGroupName string
	ContainerAppName  string
	RevisionName      string
}

// NewRevisionsApiRevisionID returns a new RevisionsApiRevisionId struct
func NewRevisionsApiRevisionID(subscriptionId string, resourceGroupName string, containerAppName string, revisionName string) RevisionsApiRevisionId {
	return RevisionsApiRevisionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ContainerAppName:  containerAppName,
		RevisionName:      revisionName,
	}
}

// ParseRevisionsApiRevisionID parses 'input' into a RevisionsApiRevisionId
func ParseRevisionsApiRevisionID(input string) (*RevisionsApiRevisionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RevisionsApiRevisionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RevisionsApiRevisionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRevisionsApiRevisionIDInsensitively parses 'input' case-insensitively into a RevisionsApiRevisionId
// note: this method should only be used for API response data and not user input
func ParseRevisionsApiRevisionIDInsensitively(input string) (*RevisionsApiRevisionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RevisionsApiRevisionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RevisionsApiRevisionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RevisionsApiRevisionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ContainerAppName, ok = input.Parsed["containerAppName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerAppName", input)
	}

	if id.RevisionName, ok = input.Parsed["revisionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "revisionName", input)
	}

	return nil
}

// ValidateRevisionsApiRevisionID checks that 'input' can be parsed as a Revisions Api Revision ID
func ValidateRevisionsApiRevisionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRevisionsApiRevisionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Revisions Api Revision ID
func (id RevisionsApiRevisionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/containerApps/%s/detectorProperties/revisionsApi/revisions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName, id.RevisionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Revisions Api Revision ID
func (id RevisionsApiRevisionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticContainerApps", "containerApps", "containerApps"),
		resourceids.UserSpecifiedSegment("containerAppName", "containerAppName"),
		resourceids.StaticSegment("staticDetectorProperties", "detectorProperties", "detectorProperties"),
		resourceids.StaticSegment("staticRevisionsApi", "revisionsApi", "revisionsApi"),
		resourceids.StaticSegment("staticRevisions", "revisions", "revisions"),
		resourceids.UserSpecifiedSegment("revisionName", "revisionName"),
	}
}

// String returns a human-readable description of this Revisions Api Revision ID
func (id RevisionsApiRevisionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container App Name: %q", id.ContainerAppName),
		fmt.Sprintf("Revision Name: %q", id.RevisionName),
	}
	return fmt.Sprintf("Revisions Api Revision (%s)", strings.Join(components, "\n"))
}
