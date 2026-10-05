package updateruns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&UpdateId{})
}

var _ resourceids.ResourceId = &UpdateId{}

// UpdateId is a struct representing the Resource ID for a Update
type UpdateId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
	UpdateName        string
}

// NewUpdateID returns a new UpdateId struct
func NewUpdateID(subscriptionId string, resourceGroupName string, clusterName string, updateName string) UpdateId {
	return UpdateId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
		UpdateName:        updateName,
	}
}

// ParseUpdateID parses 'input' into a UpdateId
func ParseUpdateID(input string) (*UpdateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UpdateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UpdateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseUpdateIDInsensitively parses 'input' case-insensitively into a UpdateId
// note: this method should only be used for API response data and not user input
func ParseUpdateIDInsensitively(input string) (*UpdateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UpdateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UpdateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *UpdateId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.UpdateName, ok = input.Parsed["updateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "updateName", input)
	}

	return nil
}

// ValidateUpdateID checks that 'input' can be parsed as a Update ID
func ValidateUpdateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUpdateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Update ID
func (id UpdateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/clusters/%s/updates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.UpdateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Update ID
func (id UpdateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticUpdates", "updates", "updates"),
		resourceids.UserSpecifiedSegment("updateName", "updateName"),
	}
}

// String returns a human-readable description of this Update ID
func (id UpdateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Update Name: %q", id.UpdateName),
	}
	return fmt.Sprintf("Update (%s)", strings.Join(components, "\n"))
}
