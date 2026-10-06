package projectpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProjectPolicyId{})
}

var _ resourceids.ResourceId = &ProjectPolicyId{}

// ProjectPolicyId is a struct representing the Resource ID for a Project Policy
type ProjectPolicyId struct {
	SubscriptionId    string
	ResourceGroupName string
	DevCenterName     string
	ProjectPolicyName string
}

// NewProjectPolicyID returns a new ProjectPolicyId struct
func NewProjectPolicyID(subscriptionId string, resourceGroupName string, devCenterName string, projectPolicyName string) ProjectPolicyId {
	return ProjectPolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DevCenterName:     devCenterName,
		ProjectPolicyName: projectPolicyName,
	}
}

// ParseProjectPolicyID parses 'input' into a ProjectPolicyId
func ParseProjectPolicyID(input string) (*ProjectPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProjectPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProjectPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProjectPolicyIDInsensitively parses 'input' case-insensitively into a ProjectPolicyId
// note: this method should only be used for API response data and not user input
func ParseProjectPolicyIDInsensitively(input string) (*ProjectPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProjectPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProjectPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProjectPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DevCenterName, ok = input.Parsed["devCenterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", input)
	}

	if id.ProjectPolicyName, ok = input.Parsed["projectPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "projectPolicyName", input)
	}

	return nil
}

// ValidateProjectPolicyID checks that 'input' can be parsed as a Project Policy ID
func ValidateProjectPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProjectPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Project Policy ID
func (id ProjectPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/projectPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.ProjectPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Project Policy ID
func (id ProjectPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticDevCenters", "devCenters", "devCenters"),
		resourceids.UserSpecifiedSegment("devCenterName", "devCenterName"),
		resourceids.StaticSegment("staticProjectPolicies", "projectPolicies", "projectPolicies"),
		resourceids.UserSpecifiedSegment("projectPolicyName", "projectPolicyName"),
	}
}

// String returns a human-readable description of this Project Policy ID
func (id ProjectPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Project Policy Name: %q", id.ProjectPolicyName),
	}
	return fmt.Sprintf("Project Policy (%s)", strings.Join(components, "\n"))
}
