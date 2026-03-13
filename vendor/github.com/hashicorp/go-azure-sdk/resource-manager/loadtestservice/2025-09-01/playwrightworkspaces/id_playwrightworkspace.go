package playwrightworkspaces

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PlaywrightWorkspaceId{})
}

var _ resourceids.ResourceId = &PlaywrightWorkspaceId{}

// PlaywrightWorkspaceId is a struct representing the Resource ID for a Playwright Workspace
type PlaywrightWorkspaceId struct {
	SubscriptionId          string
	ResourceGroupName       string
	PlaywrightWorkspaceName string
}

// NewPlaywrightWorkspaceID returns a new PlaywrightWorkspaceId struct
func NewPlaywrightWorkspaceID(subscriptionId string, resourceGroupName string, playwrightWorkspaceName string) PlaywrightWorkspaceId {
	return PlaywrightWorkspaceId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		PlaywrightWorkspaceName: playwrightWorkspaceName,
	}
}

// ParsePlaywrightWorkspaceID parses 'input' into a PlaywrightWorkspaceId
func ParsePlaywrightWorkspaceID(input string) (*PlaywrightWorkspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PlaywrightWorkspaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PlaywrightWorkspaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePlaywrightWorkspaceIDInsensitively parses 'input' case-insensitively into a PlaywrightWorkspaceId
// note: this method should only be used for API response data and not user input
func ParsePlaywrightWorkspaceIDInsensitively(input string) (*PlaywrightWorkspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PlaywrightWorkspaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PlaywrightWorkspaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PlaywrightWorkspaceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PlaywrightWorkspaceName, ok = input.Parsed["playwrightWorkspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "playwrightWorkspaceName", input)
	}

	return nil
}

// ValidatePlaywrightWorkspaceID checks that 'input' can be parsed as a Playwright Workspace ID
func ValidatePlaywrightWorkspaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePlaywrightWorkspaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Playwright Workspace ID
func (id PlaywrightWorkspaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.LoadTestService/playwrightWorkspaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PlaywrightWorkspaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Playwright Workspace ID
func (id PlaywrightWorkspaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLoadTestService", "Microsoft.LoadTestService", "Microsoft.LoadTestService"),
		resourceids.StaticSegment("staticPlaywrightWorkspaces", "playwrightWorkspaces", "playwrightWorkspaces"),
		resourceids.UserSpecifiedSegment("playwrightWorkspaceName", "playwrightWorkspaceName"),
	}
}

// String returns a human-readable description of this Playwright Workspace ID
func (id PlaywrightWorkspaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Playwright Workspace Name: %q", id.PlaywrightWorkspaceName),
	}
	return fmt.Sprintf("Playwright Workspace (%s)", strings.Join(components, "\n"))
}
