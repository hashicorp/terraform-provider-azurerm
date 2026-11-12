package goldengatedeployments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AssignedConnectionId{})
}

var _ resourceids.ResourceId = &AssignedConnectionId{}

// AssignedConnectionId is a struct representing the Resource ID for a Assigned Connection
type AssignedConnectionId struct {
	SubscriptionId           string
	ResourceGroupName        string
	GoldenGateDeploymentName string
	AssignmentId             string
}

// NewAssignedConnectionID returns a new AssignedConnectionId struct
func NewAssignedConnectionID(subscriptionId string, resourceGroupName string, goldenGateDeploymentName string, assignmentId string) AssignedConnectionId {
	return AssignedConnectionId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		GoldenGateDeploymentName: goldenGateDeploymentName,
		AssignmentId:             assignmentId,
	}
}

// ParseAssignedConnectionID parses 'input' into a AssignedConnectionId
func ParseAssignedConnectionID(input string) (*AssignedConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AssignedConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssignedConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAssignedConnectionIDInsensitively parses 'input' case-insensitively into a AssignedConnectionId
// note: this method should only be used for API response data and not user input
func ParseAssignedConnectionIDInsensitively(input string) (*AssignedConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AssignedConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssignedConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AssignedConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.GoldenGateDeploymentName, ok = input.Parsed["goldenGateDeploymentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "goldenGateDeploymentName", input)
	}

	if id.AssignmentId, ok = input.Parsed["assignmentId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "assignmentId", input)
	}

	return nil
}

// ValidateAssignedConnectionID checks that 'input' can be parsed as a Assigned Connection ID
func ValidateAssignedConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAssignedConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Assigned Connection ID
func (id AssignedConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/goldenGateDeployments/%s/assignedConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GoldenGateDeploymentName, id.AssignmentId)
}

// Segments returns a slice of Resource ID Segments which comprise this Assigned Connection ID
func (id AssignedConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticGoldenGateDeployments", "goldenGateDeployments", "goldenGateDeployments"),
		resourceids.UserSpecifiedSegment("goldenGateDeploymentName", "goldenGateDeploymentName"),
		resourceids.StaticSegment("staticAssignedConnections", "assignedConnections", "assignedConnections"),
		resourceids.UserSpecifiedSegment("assignmentId", "assignmentId"),
	}
}

// String returns a human-readable description of this Assigned Connection ID
func (id AssignedConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Golden Gate Deployment Name: %q", id.GoldenGateDeploymentName),
		fmt.Sprintf("Assignment: %q", id.AssignmentId),
	}
	return fmt.Sprintf("Assigned Connection (%s)", strings.Join(components, "\n"))
}
