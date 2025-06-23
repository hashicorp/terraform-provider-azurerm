package standbycontainergrouppools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StandbyContainerGroupPoolId{})
}

var _ resourceids.ResourceId = &StandbyContainerGroupPoolId{}

// StandbyContainerGroupPoolId is a struct representing the Resource ID for a Standby Container Group Pool
type StandbyContainerGroupPoolId struct {
	SubscriptionId                string
	ResourceGroupName             string
	StandbyContainerGroupPoolName string
}

// NewStandbyContainerGroupPoolID returns a new StandbyContainerGroupPoolId struct
func NewStandbyContainerGroupPoolID(subscriptionId string, resourceGroupName string, standbyContainerGroupPoolName string) StandbyContainerGroupPoolId {
	return StandbyContainerGroupPoolId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		StandbyContainerGroupPoolName: standbyContainerGroupPoolName,
	}
}

// ParseStandbyContainerGroupPoolID parses 'input' into a StandbyContainerGroupPoolId
func ParseStandbyContainerGroupPoolID(input string) (*StandbyContainerGroupPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StandbyContainerGroupPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StandbyContainerGroupPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStandbyContainerGroupPoolIDInsensitively parses 'input' case-insensitively into a StandbyContainerGroupPoolId
// note: this method should only be used for API response data and not user input
func ParseStandbyContainerGroupPoolIDInsensitively(input string) (*StandbyContainerGroupPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StandbyContainerGroupPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StandbyContainerGroupPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StandbyContainerGroupPoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StandbyContainerGroupPoolName, ok = input.Parsed["standbyContainerGroupPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "standbyContainerGroupPoolName", input)
	}

	return nil
}

// ValidateStandbyContainerGroupPoolID checks that 'input' can be parsed as a Standby Container Group Pool ID
func ValidateStandbyContainerGroupPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStandbyContainerGroupPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Standby Container Group Pool ID
func (id StandbyContainerGroupPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StandbyPool/standbyContainerGroupPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StandbyContainerGroupPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Standby Container Group Pool ID
func (id StandbyContainerGroupPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStandbyPool", "Microsoft.StandbyPool", "Microsoft.StandbyPool"),
		resourceids.StaticSegment("staticStandbyContainerGroupPools", "standbyContainerGroupPools", "standbyContainerGroupPools"),
		resourceids.UserSpecifiedSegment("standbyContainerGroupPoolName", "standbyContainerGroupPoolName"),
	}
}

// String returns a human-readable description of this Standby Container Group Pool ID
func (id StandbyContainerGroupPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Standby Container Group Pool Name: %q", id.StandbyContainerGroupPoolName),
	}
	return fmt.Sprintf("Standby Container Group Pool (%s)", strings.Join(components, "\n"))
}
