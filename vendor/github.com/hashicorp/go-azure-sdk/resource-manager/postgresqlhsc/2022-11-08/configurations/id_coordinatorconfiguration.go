package configurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CoordinatorConfigurationId{})
}

var _ resourceids.ResourceId = &CoordinatorConfigurationId{}

// CoordinatorConfigurationId is a struct representing the Resource ID for a Coordinator Configuration
type CoordinatorConfigurationId struct {
	SubscriptionId               string
	ResourceGroupName            string
	ServerGroupsv2Name           string
	CoordinatorConfigurationName string
}

// NewCoordinatorConfigurationID returns a new CoordinatorConfigurationId struct
func NewCoordinatorConfigurationID(subscriptionId string, resourceGroupName string, serverGroupsv2Name string, coordinatorConfigurationName string) CoordinatorConfigurationId {
	return CoordinatorConfigurationId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		ServerGroupsv2Name:           serverGroupsv2Name,
		CoordinatorConfigurationName: coordinatorConfigurationName,
	}
}

// ParseCoordinatorConfigurationID parses 'input' into a CoordinatorConfigurationId
func ParseCoordinatorConfigurationID(input string) (*CoordinatorConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CoordinatorConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CoordinatorConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCoordinatorConfigurationIDInsensitively parses 'input' case-insensitively into a CoordinatorConfigurationId
// note: this method should only be used for API response data and not user input
func ParseCoordinatorConfigurationIDInsensitively(input string) (*CoordinatorConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CoordinatorConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CoordinatorConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CoordinatorConfigurationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServerGroupsv2Name, ok = input.Parsed["serverGroupsv2Name"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serverGroupsv2Name", input)
	}

	if id.CoordinatorConfigurationName, ok = input.Parsed["coordinatorConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "coordinatorConfigurationName", input)
	}

	return nil
}

// ValidateCoordinatorConfigurationID checks that 'input' can be parsed as a Coordinator Configuration ID
func ValidateCoordinatorConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCoordinatorConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Coordinator Configuration ID
func (id CoordinatorConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/%s/coordinatorConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerGroupsv2Name, id.CoordinatorConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Coordinator Configuration ID
func (id CoordinatorConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticServerGroupsv2", "serverGroupsv2", "serverGroupsv2"),
		resourceids.UserSpecifiedSegment("serverGroupsv2Name", "serverGroupsv2Name"),
		resourceids.StaticSegment("staticCoordinatorConfigurations", "coordinatorConfigurations", "coordinatorConfigurations"),
		resourceids.UserSpecifiedSegment("coordinatorConfigurationName", "coordinatorConfigurationName"),
	}
}

// String returns a human-readable description of this Coordinator Configuration ID
func (id CoordinatorConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Groupsv 2 Name: %q", id.ServerGroupsv2Name),
		fmt.Sprintf("Coordinator Configuration Name: %q", id.CoordinatorConfigurationName),
	}
	return fmt.Sprintf("Coordinator Configuration (%s)", strings.Join(components, "\n"))
}
