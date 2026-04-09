package staticsites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BuildDatabaseConnectionId{})
}

var _ resourceids.ResourceId = &BuildDatabaseConnectionId{}

// BuildDatabaseConnectionId is a struct representing the Resource ID for a Build Database Connection
type BuildDatabaseConnectionId struct {
	SubscriptionId         string
	ResourceGroupName      string
	StaticSiteName         string
	BuildName              string
	DatabaseConnectionName string
}

// NewBuildDatabaseConnectionID returns a new BuildDatabaseConnectionId struct
func NewBuildDatabaseConnectionID(subscriptionId string, resourceGroupName string, staticSiteName string, buildName string, databaseConnectionName string) BuildDatabaseConnectionId {
	return BuildDatabaseConnectionId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		StaticSiteName:         staticSiteName,
		BuildName:              buildName,
		DatabaseConnectionName: databaseConnectionName,
	}
}

// ParseBuildDatabaseConnectionID parses 'input' into a BuildDatabaseConnectionId
func ParseBuildDatabaseConnectionID(input string) (*BuildDatabaseConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuildDatabaseConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuildDatabaseConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBuildDatabaseConnectionIDInsensitively parses 'input' case-insensitively into a BuildDatabaseConnectionId
// note: this method should only be used for API response data and not user input
func ParseBuildDatabaseConnectionIDInsensitively(input string) (*BuildDatabaseConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuildDatabaseConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuildDatabaseConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BuildDatabaseConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StaticSiteName, ok = input.Parsed["staticSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "staticSiteName", input)
	}

	if id.BuildName, ok = input.Parsed["buildName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "buildName", input)
	}

	if id.DatabaseConnectionName, ok = input.Parsed["databaseConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseConnectionName", input)
	}

	return nil
}

// ValidateBuildDatabaseConnectionID checks that 'input' can be parsed as a Build Database Connection ID
func ValidateBuildDatabaseConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBuildDatabaseConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Build Database Connection ID
func (id BuildDatabaseConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/staticSites/%s/builds/%s/databaseConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName, id.BuildName, id.DatabaseConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Build Database Connection ID
func (id BuildDatabaseConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticStaticSites", "staticSites", "staticSites"),
		resourceids.UserSpecifiedSegment("staticSiteName", "staticSiteName"),
		resourceids.StaticSegment("staticBuilds", "builds", "builds"),
		resourceids.UserSpecifiedSegment("buildName", "buildName"),
		resourceids.StaticSegment("staticDatabaseConnections", "databaseConnections", "databaseConnections"),
		resourceids.UserSpecifiedSegment("databaseConnectionName", "databaseConnectionName"),
	}
}

// String returns a human-readable description of this Build Database Connection ID
func (id BuildDatabaseConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Static Site Name: %q", id.StaticSiteName),
		fmt.Sprintf("Build Name: %q", id.BuildName),
		fmt.Sprintf("Database Connection Name: %q", id.DatabaseConnectionName),
	}
	return fmt.Sprintf("Build Database Connection (%s)", strings.Join(components, "\n"))
}
