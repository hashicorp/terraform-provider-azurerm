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
	recaser.RegisterResourceId(&BuildLinkedBackendId{})
}

var _ resourceids.ResourceId = &BuildLinkedBackendId{}

// BuildLinkedBackendId is a struct representing the Resource ID for a Build Linked Backend
type BuildLinkedBackendId struct {
	SubscriptionId    string
	ResourceGroupName string
	StaticSiteName    string
	BuildName         string
	LinkedBackendName string
}

// NewBuildLinkedBackendID returns a new BuildLinkedBackendId struct
func NewBuildLinkedBackendID(subscriptionId string, resourceGroupName string, staticSiteName string, buildName string, linkedBackendName string) BuildLinkedBackendId {
	return BuildLinkedBackendId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StaticSiteName:    staticSiteName,
		BuildName:         buildName,
		LinkedBackendName: linkedBackendName,
	}
}

// ParseBuildLinkedBackendID parses 'input' into a BuildLinkedBackendId
func ParseBuildLinkedBackendID(input string) (*BuildLinkedBackendId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuildLinkedBackendId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuildLinkedBackendId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBuildLinkedBackendIDInsensitively parses 'input' case-insensitively into a BuildLinkedBackendId
// note: this method should only be used for API response data and not user input
func ParseBuildLinkedBackendIDInsensitively(input string) (*BuildLinkedBackendId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuildLinkedBackendId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuildLinkedBackendId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BuildLinkedBackendId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.LinkedBackendName, ok = input.Parsed["linkedBackendName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkedBackendName", input)
	}

	return nil
}

// ValidateBuildLinkedBackendID checks that 'input' can be parsed as a Build Linked Backend ID
func ValidateBuildLinkedBackendID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBuildLinkedBackendID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Build Linked Backend ID
func (id BuildLinkedBackendId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/staticSites/%s/builds/%s/linkedBackends/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName, id.BuildName, id.LinkedBackendName)
}

// Segments returns a slice of Resource ID Segments which comprise this Build Linked Backend ID
func (id BuildLinkedBackendId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticLinkedBackends", "linkedBackends", "linkedBackends"),
		resourceids.UserSpecifiedSegment("linkedBackendName", "linkedBackendName"),
	}
}

// String returns a human-readable description of this Build Linked Backend ID
func (id BuildLinkedBackendId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Static Site Name: %q", id.StaticSiteName),
		fmt.Sprintf("Build Name: %q", id.BuildName),
		fmt.Sprintf("Linked Backend Name: %q", id.LinkedBackendName),
	}
	return fmt.Sprintf("Build Linked Backend (%s)", strings.Join(components, "\n"))
}
