package machinelearningcomputes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ComputeId{}

// ComputeId is a struct representing the Resource ID for a Compute
type ComputeId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	ComputeName       string
}

// NewComputeID returns a new ComputeId struct
func NewComputeID(subscriptionId string, resourceGroupName string, workspaceName string, computeName string) ComputeId {
	return ComputeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		ComputeName:       computeName,
	}
}

// ParseComputeID parses 'input' into a ComputeId
func ParseComputeID(input string) (*ComputeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ComputeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ComputeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.ComputeName, ok = parsed.Parsed["computeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "computeName", *parsed)
	}

	return &id, nil
}

// ParseComputeIDInsensitively parses 'input' case-insensitively into a ComputeId
// note: this method should only be used for API response data and not user input
func ParseComputeIDInsensitively(input string) (*ComputeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ComputeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ComputeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.ComputeName, ok = parsed.Parsed["computeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "computeName", *parsed)
	}

	return &id, nil
}

// ValidateComputeID checks that 'input' can be parsed as a Compute ID
func ValidateComputeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseComputeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Compute ID
func (id ComputeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/workspaces/%s/computes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.ComputeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Compute ID
func (id ComputeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticComputes", "computes", "computes"),
		resourceids.UserSpecifiedSegment("computeName", "computeValue"),
	}
}

// String returns a human-readable description of this Compute ID
func (id ComputeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Compute Name: %q", id.ComputeName),
	}
	return fmt.Sprintf("Compute (%s)", strings.Join(components, "\n"))
}
