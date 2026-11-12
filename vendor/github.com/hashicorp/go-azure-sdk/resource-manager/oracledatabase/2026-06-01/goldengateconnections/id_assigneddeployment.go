package goldengateconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AssignedDeploymentId{})
}

var _ resourceids.ResourceId = &AssignedDeploymentId{}

// AssignedDeploymentId is a struct representing the Resource ID for a Assigned Deployment
type AssignedDeploymentId struct {
	SubscriptionId           string
	ResourceGroupName        string
	GoldenGateConnectionName string
	AssignmentId             string
}

// NewAssignedDeploymentID returns a new AssignedDeploymentId struct
func NewAssignedDeploymentID(subscriptionId string, resourceGroupName string, goldenGateConnectionName string, assignmentId string) AssignedDeploymentId {
	return AssignedDeploymentId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		GoldenGateConnectionName: goldenGateConnectionName,
		AssignmentId:             assignmentId,
	}
}

// ParseAssignedDeploymentID parses 'input' into a AssignedDeploymentId
func ParseAssignedDeploymentID(input string) (*AssignedDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AssignedDeploymentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssignedDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAssignedDeploymentIDInsensitively parses 'input' case-insensitively into a AssignedDeploymentId
// note: this method should only be used for API response data and not user input
func ParseAssignedDeploymentIDInsensitively(input string) (*AssignedDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AssignedDeploymentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssignedDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AssignedDeploymentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.GoldenGateConnectionName, ok = input.Parsed["goldenGateConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "goldenGateConnectionName", input)
	}

	if id.AssignmentId, ok = input.Parsed["assignmentId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "assignmentId", input)
	}

	return nil
}

// ValidateAssignedDeploymentID checks that 'input' can be parsed as a Assigned Deployment ID
func ValidateAssignedDeploymentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAssignedDeploymentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Assigned Deployment ID
func (id AssignedDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/goldenGateConnections/%s/assignedDeployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GoldenGateConnectionName, id.AssignmentId)
}

// Segments returns a slice of Resource ID Segments which comprise this Assigned Deployment ID
func (id AssignedDeploymentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticGoldenGateConnections", "goldenGateConnections", "goldenGateConnections"),
		resourceids.UserSpecifiedSegment("goldenGateConnectionName", "goldenGateConnectionName"),
		resourceids.StaticSegment("staticAssignedDeployments", "assignedDeployments", "assignedDeployments"),
		resourceids.UserSpecifiedSegment("assignmentId", "assignmentId"),
	}
}

// String returns a human-readable description of this Assigned Deployment ID
func (id AssignedDeploymentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Golden Gate Connection Name: %q", id.GoldenGateConnectionName),
		fmt.Sprintf("Assignment: %q", id.AssignmentId),
	}
	return fmt.Sprintf("Assigned Deployment (%s)", strings.Join(components, "\n"))
}
