package schedules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PoolId{})
}

var _ resourceids.ResourceId = &PoolId{}

// PoolId is a struct representing the Resource ID for a Pool
type PoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProjectName       string
	PoolName          string
}

// NewPoolID returns a new PoolId struct
func NewPoolID(subscriptionId string, resourceGroupName string, projectName string, poolName string) PoolId {
	return PoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProjectName:       projectName,
		PoolName:          poolName,
	}
}

// ParsePoolID parses 'input' into a PoolId
func ParsePoolID(input string) (*PoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePoolIDInsensitively parses 'input' case-insensitively into a PoolId
// note: this method should only be used for API response data and not user input
func ParsePoolIDInsensitively(input string) (*PoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProjectName, ok = input.Parsed["projectName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "projectName", input)
	}

	if id.PoolName, ok = input.Parsed["poolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "poolName", input)
	}

	return nil
}

// ValidatePoolID checks that 'input' can be parsed as a Pool ID
func ValidatePoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Pool ID
func (id PoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/pools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.PoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Pool ID
func (id PoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
		resourceids.StaticSegment("staticPools", "pools", "pools"),
		resourceids.UserSpecifiedSegment("poolName", "poolName"),
	}
}

// String returns a human-readable description of this Pool ID
func (id PoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Pool Name: %q", id.PoolName),
	}
	return fmt.Sprintf("Pool (%s)", strings.Join(components, "\n"))
}
