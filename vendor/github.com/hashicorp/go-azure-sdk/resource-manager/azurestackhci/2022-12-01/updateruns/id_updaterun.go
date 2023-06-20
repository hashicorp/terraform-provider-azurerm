package updateruns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = UpdateRunId{}

// UpdateRunId is a struct representing the Resource ID for a Update Run
type UpdateRunId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
	UpdateName        string
	UpdateRunName     string
}

// NewUpdateRunID returns a new UpdateRunId struct
func NewUpdateRunID(subscriptionId string, resourceGroupName string, clusterName string, updateName string, updateRunName string) UpdateRunId {
	return UpdateRunId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
		UpdateName:        updateName,
		UpdateRunName:     updateRunName,
	}
}

// ParseUpdateRunID parses 'input' into a UpdateRunId
func ParseUpdateRunID(input string) (*UpdateRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(UpdateRunId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UpdateRunId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clusterName", *parsed)
	}

	if id.UpdateName, ok = parsed.Parsed["updateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "updateName", *parsed)
	}

	if id.UpdateRunName, ok = parsed.Parsed["updateRunName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "updateRunName", *parsed)
	}

	return &id, nil
}

// ParseUpdateRunIDInsensitively parses 'input' case-insensitively into a UpdateRunId
// note: this method should only be used for API response data and not user input
func ParseUpdateRunIDInsensitively(input string) (*UpdateRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(UpdateRunId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UpdateRunId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clusterName", *parsed)
	}

	if id.UpdateName, ok = parsed.Parsed["updateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "updateName", *parsed)
	}

	if id.UpdateRunName, ok = parsed.Parsed["updateRunName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "updateRunName", *parsed)
	}

	return &id, nil
}

// ValidateUpdateRunID checks that 'input' can be parsed as a Update Run ID
func ValidateUpdateRunID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUpdateRunID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Update Run ID
func (id UpdateRunId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/clusters/%s/updates/%s/updateRuns/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.UpdateName, id.UpdateRunName)
}

// Segments returns a slice of Resource ID Segments which comprise this Update Run ID
func (id UpdateRunId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
		resourceids.StaticSegment("staticUpdates", "updates", "updates"),
		resourceids.UserSpecifiedSegment("updateName", "updateValue"),
		resourceids.StaticSegment("staticUpdateRuns", "updateRuns", "updateRuns"),
		resourceids.UserSpecifiedSegment("updateRunName", "updateRunValue"),
	}
}

// String returns a human-readable description of this Update Run ID
func (id UpdateRunId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Update Name: %q", id.UpdateName),
		fmt.Sprintf("Update Run Name: %q", id.UpdateRunName),
	}
	return fmt.Sprintf("Update Run (%s)", strings.Join(components, "\n"))
}
