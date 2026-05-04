package applicationdefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApplicationDefinitionId{})
}

var _ resourceids.ResourceId = &ApplicationDefinitionId{}

// ApplicationDefinitionId is a struct representing the Resource ID for a Application Definition
type ApplicationDefinitionId struct {
	SubscriptionId            string
	ResourceGroupName         string
	ApplicationDefinitionName string
}

// NewApplicationDefinitionID returns a new ApplicationDefinitionId struct
func NewApplicationDefinitionID(subscriptionId string, resourceGroupName string, applicationDefinitionName string) ApplicationDefinitionId {
	return ApplicationDefinitionId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		ApplicationDefinitionName: applicationDefinitionName,
	}
}

// ParseApplicationDefinitionID parses 'input' into a ApplicationDefinitionId
func ParseApplicationDefinitionID(input string) (*ApplicationDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationDefinitionIDInsensitively parses 'input' case-insensitively into a ApplicationDefinitionId
// note: this method should only be used for API response data and not user input
func ParseApplicationDefinitionIDInsensitively(input string) (*ApplicationDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ApplicationDefinitionName, ok = input.Parsed["applicationDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationDefinitionName", input)
	}

	return nil
}

// ValidateApplicationDefinitionID checks that 'input' can be parsed as a Application Definition ID
func ValidateApplicationDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Definition ID
func (id ApplicationDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Solutions/applicationDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Definition ID
func (id ApplicationDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSolutions", "Microsoft.Solutions", "Microsoft.Solutions"),
		resourceids.StaticSegment("staticApplicationDefinitions", "applicationDefinitions", "applicationDefinitions"),
		resourceids.UserSpecifiedSegment("applicationDefinitionName", "applicationDefinitionName"),
	}
}

// String returns a human-readable description of this Application Definition ID
func (id ApplicationDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Definition Name: %q", id.ApplicationDefinitionName),
	}
	return fmt.Sprintf("Application Definition (%s)", strings.Join(components, "\n"))
}
