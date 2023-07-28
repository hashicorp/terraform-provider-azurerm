package environmentversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = EnvironmentVersionId{}

// EnvironmentVersionId is a struct representing the Resource ID for a Environment Version
type EnvironmentVersionId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	EnvironmentName   string
	VersionName       string
}

// NewEnvironmentVersionID returns a new EnvironmentVersionId struct
func NewEnvironmentVersionID(subscriptionId string, resourceGroupName string, workspaceName string, environmentName string, versionName string) EnvironmentVersionId {
	return EnvironmentVersionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		EnvironmentName:   environmentName,
		VersionName:       versionName,
	}
}

// ParseEnvironmentVersionID parses 'input' into a EnvironmentVersionId
func ParseEnvironmentVersionID(input string) (*EnvironmentVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(EnvironmentVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EnvironmentVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.EnvironmentName, ok = parsed.Parsed["environmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "environmentName", *parsed)
	}

	if id.VersionName, ok = parsed.Parsed["versionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionName", *parsed)
	}

	return &id, nil
}

// ParseEnvironmentVersionIDInsensitively parses 'input' case-insensitively into a EnvironmentVersionId
// note: this method should only be used for API response data and not user input
func ParseEnvironmentVersionIDInsensitively(input string) (*EnvironmentVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(EnvironmentVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EnvironmentVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.EnvironmentName, ok = parsed.Parsed["environmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "environmentName", *parsed)
	}

	if id.VersionName, ok = parsed.Parsed["versionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionName", *parsed)
	}

	return &id, nil
}

// ValidateEnvironmentVersionID checks that 'input' can be parsed as a Environment Version ID
func ValidateEnvironmentVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEnvironmentVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Environment Version ID
func (id EnvironmentVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/workspaces/%s/environments/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.EnvironmentName, id.VersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Environment Version ID
func (id EnvironmentVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticEnvironments", "environments", "environments"),
		resourceids.UserSpecifiedSegment("environmentName", "environmentValue"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionValue"),
	}
}

// String returns a human-readable description of this Environment Version ID
func (id EnvironmentVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Environment Name: %q", id.EnvironmentName),
		fmt.Sprintf("Version Name: %q", id.VersionName),
	}
	return fmt.Sprintf("Environment Version (%s)", strings.Join(components, "\n"))
}
