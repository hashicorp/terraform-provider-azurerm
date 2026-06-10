package attachednetworkconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AttachedNetworkId{})
}

var _ resourceids.ResourceId = &AttachedNetworkId{}

// AttachedNetworkId is a struct representing the Resource ID for a Attached Network
type AttachedNetworkId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ProjectName         string
	AttachedNetworkName string
}

// NewAttachedNetworkID returns a new AttachedNetworkId struct
func NewAttachedNetworkID(subscriptionId string, resourceGroupName string, projectName string, attachedNetworkName string) AttachedNetworkId {
	return AttachedNetworkId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ProjectName:         projectName,
		AttachedNetworkName: attachedNetworkName,
	}
}

// ParseAttachedNetworkID parses 'input' into a AttachedNetworkId
func ParseAttachedNetworkID(input string) (*AttachedNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AttachedNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AttachedNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAttachedNetworkIDInsensitively parses 'input' case-insensitively into a AttachedNetworkId
// note: this method should only be used for API response data and not user input
func ParseAttachedNetworkIDInsensitively(input string) (*AttachedNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AttachedNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AttachedNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AttachedNetworkId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AttachedNetworkName, ok = input.Parsed["attachedNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "attachedNetworkName", input)
	}

	return nil
}

// ValidateAttachedNetworkID checks that 'input' can be parsed as a Attached Network ID
func ValidateAttachedNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAttachedNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Attached Network ID
func (id AttachedNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/attachedNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.AttachedNetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Attached Network ID
func (id AttachedNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
		resourceids.StaticSegment("staticAttachedNetworks", "attachedNetworks", "attachedNetworks"),
		resourceids.UserSpecifiedSegment("attachedNetworkName", "attachedNetworkName"),
	}
}

// String returns a human-readable description of this Attached Network ID
func (id AttachedNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Attached Network Name: %q", id.AttachedNetworkName),
	}
	return fmt.Sprintf("Attached Network (%s)", strings.Join(components, "\n"))
}
