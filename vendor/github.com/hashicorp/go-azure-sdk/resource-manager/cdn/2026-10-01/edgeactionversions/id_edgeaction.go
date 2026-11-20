package edgeactionversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&EdgeActionId{})
}

var _ resourceids.ResourceId = &EdgeActionId{}

// EdgeActionId is a struct representing the Resource ID for a Edge Action
type EdgeActionId struct {
	SubscriptionId    string
	ResourceGroupName string
	EdgeActionName    string
}

// NewEdgeActionID returns a new EdgeActionId struct
func NewEdgeActionID(subscriptionId string, resourceGroupName string, edgeActionName string) EdgeActionId {
	return EdgeActionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		EdgeActionName:    edgeActionName,
	}
}

// ParseEdgeActionID parses 'input' into a EdgeActionId
func ParseEdgeActionID(input string) (*EdgeActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EdgeActionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EdgeActionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEdgeActionIDInsensitively parses 'input' case-insensitively into a EdgeActionId
// note: this method should only be used for API response data and not user input
func ParseEdgeActionIDInsensitively(input string) (*EdgeActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EdgeActionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EdgeActionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EdgeActionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.EdgeActionName, ok = input.Parsed["edgeActionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "edgeActionName", input)
	}

	return nil
}

// ValidateEdgeActionID checks that 'input' can be parsed as a Edge Action ID
func ValidateEdgeActionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEdgeActionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Edge Action ID
func (id EdgeActionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/edgeActions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Edge Action ID
func (id EdgeActionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCdn", "Microsoft.Cdn", "Microsoft.Cdn"),
		resourceids.StaticSegment("staticEdgeActions", "edgeActions", "edgeActions"),
		resourceids.UserSpecifiedSegment("edgeActionName", "edgeActionName"),
	}
}

// String returns a human-readable description of this Edge Action ID
func (id EdgeActionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Edge Action Name: %q", id.EdgeActionName),
	}
	return fmt.Sprintf("Edge Action (%s)", strings.Join(components, "\n"))
}
