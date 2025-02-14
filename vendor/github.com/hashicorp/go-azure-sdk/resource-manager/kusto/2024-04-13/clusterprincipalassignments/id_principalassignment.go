package clusterprincipalassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PrincipalAssignmentId{})
}

var _ resourceids.ResourceId = &PrincipalAssignmentId{}

// PrincipalAssignmentId is a struct representing the Resource ID for a Principal Assignment
type PrincipalAssignmentId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ClusterName             string
	PrincipalAssignmentName string
}

// NewPrincipalAssignmentID returns a new PrincipalAssignmentId struct
func NewPrincipalAssignmentID(subscriptionId string, resourceGroupName string, clusterName string, principalAssignmentName string) PrincipalAssignmentId {
	return PrincipalAssignmentId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ClusterName:             clusterName,
		PrincipalAssignmentName: principalAssignmentName,
	}
}

// ParsePrincipalAssignmentID parses 'input' into a PrincipalAssignmentId
func ParsePrincipalAssignmentID(input string) (*PrincipalAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrincipalAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrincipalAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrincipalAssignmentIDInsensitively parses 'input' case-insensitively into a PrincipalAssignmentId
// note: this method should only be used for API response data and not user input
func ParsePrincipalAssignmentIDInsensitively(input string) (*PrincipalAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrincipalAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrincipalAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrincipalAssignmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ClusterName, ok = input.Parsed["clusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clusterName", input)
	}

	if id.PrincipalAssignmentName, ok = input.Parsed["principalAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "principalAssignmentName", input)
	}

	return nil
}

// ValidatePrincipalAssignmentID checks that 'input' can be parsed as a Principal Assignment ID
func ValidatePrincipalAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrincipalAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Principal Assignment ID
func (id PrincipalAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/principalAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.PrincipalAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Principal Assignment ID
func (id PrincipalAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKusto", "Microsoft.Kusto", "Microsoft.Kusto"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticPrincipalAssignments", "principalAssignments", "principalAssignments"),
		resourceids.UserSpecifiedSegment("principalAssignmentName", "principalAssignmentName"),
	}
}

// String returns a human-readable description of this Principal Assignment ID
func (id PrincipalAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Principal Assignment Name: %q", id.PrincipalAssignmentName),
	}
	return fmt.Sprintf("Principal Assignment (%s)", strings.Join(components, "\n"))
}
