package instancefailovergroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&InstanceFailoverGroupId{})
}

var _ resourceids.ResourceId = &InstanceFailoverGroupId{}

// InstanceFailoverGroupId is a struct representing the Resource ID for a Instance Failover Group
type InstanceFailoverGroupId struct {
	SubscriptionId            string
	ResourceGroupName         string
	LocationName              string
	InstanceFailoverGroupName string
}

// NewInstanceFailoverGroupID returns a new InstanceFailoverGroupId struct
func NewInstanceFailoverGroupID(subscriptionId string, resourceGroupName string, locationName string, instanceFailoverGroupName string) InstanceFailoverGroupId {
	return InstanceFailoverGroupId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		LocationName:              locationName,
		InstanceFailoverGroupName: instanceFailoverGroupName,
	}
}

// ParseInstanceFailoverGroupID parses 'input' into a InstanceFailoverGroupId
func ParseInstanceFailoverGroupID(input string) (*InstanceFailoverGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InstanceFailoverGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InstanceFailoverGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseInstanceFailoverGroupIDInsensitively parses 'input' case-insensitively into a InstanceFailoverGroupId
// note: this method should only be used for API response data and not user input
func ParseInstanceFailoverGroupIDInsensitively(input string) (*InstanceFailoverGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InstanceFailoverGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InstanceFailoverGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *InstanceFailoverGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.InstanceFailoverGroupName, ok = input.Parsed["instanceFailoverGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instanceFailoverGroupName", input)
	}

	return nil
}

// ValidateInstanceFailoverGroupID checks that 'input' can be parsed as a Instance Failover Group ID
func ValidateInstanceFailoverGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInstanceFailoverGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Instance Failover Group ID
func (id InstanceFailoverGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/locations/%s/instanceFailoverGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocationName, id.InstanceFailoverGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Instance Failover Group ID
func (id InstanceFailoverGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticInstanceFailoverGroups", "instanceFailoverGroups", "instanceFailoverGroups"),
		resourceids.UserSpecifiedSegment("instanceFailoverGroupName", "instanceFailoverGroupName"),
	}
}

// String returns a human-readable description of this Instance Failover Group ID
func (id InstanceFailoverGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Instance Failover Group Name: %q", id.InstanceFailoverGroupName),
	}
	return fmt.Sprintf("Instance Failover Group (%s)", strings.Join(components, "\n"))
}
