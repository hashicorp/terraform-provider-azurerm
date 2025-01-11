package managedinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&InstancePoolId{})
}

var _ resourceids.ResourceId = &InstancePoolId{}

// InstancePoolId is a struct representing the Resource ID for a Instance Pool
type InstancePoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	InstancePoolName  string
}

// NewInstancePoolID returns a new InstancePoolId struct
func NewInstancePoolID(subscriptionId string, resourceGroupName string, instancePoolName string) InstancePoolId {
	return InstancePoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		InstancePoolName:  instancePoolName,
	}
}

// ParseInstancePoolID parses 'input' into a InstancePoolId
func ParseInstancePoolID(input string) (*InstancePoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InstancePoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InstancePoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseInstancePoolIDInsensitively parses 'input' case-insensitively into a InstancePoolId
// note: this method should only be used for API response data and not user input
func ParseInstancePoolIDInsensitively(input string) (*InstancePoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InstancePoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InstancePoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *InstancePoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.InstancePoolName, ok = input.Parsed["instancePoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instancePoolName", input)
	}

	return nil
}

// ValidateInstancePoolID checks that 'input' can be parsed as a Instance Pool ID
func ValidateInstancePoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInstancePoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Instance Pool ID
func (id InstancePoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/instancePools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.InstancePoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Instance Pool ID
func (id InstancePoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticInstancePools", "instancePools", "instancePools"),
		resourceids.UserSpecifiedSegment("instancePoolName", "instancePoolName"),
	}
}

// String returns a human-readable description of this Instance Pool ID
func (id InstancePoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Instance Pool Name: %q", id.InstancePoolName),
	}
	return fmt.Sprintf("Instance Pool (%s)", strings.Join(components, "\n"))
}
