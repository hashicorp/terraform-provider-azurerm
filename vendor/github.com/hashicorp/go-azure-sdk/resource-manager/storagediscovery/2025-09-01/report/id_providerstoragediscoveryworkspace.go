package report

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProviderStorageDiscoveryWorkspaceId{})
}

var _ resourceids.ResourceId = &ProviderStorageDiscoveryWorkspaceId{}

// ProviderStorageDiscoveryWorkspaceId is a struct representing the Resource ID for a Provider Storage Discovery Workspace
type ProviderStorageDiscoveryWorkspaceId struct {
	SubscriptionId                string
	ResourceGroupName             string
	StorageDiscoveryWorkspaceName string
}

// NewProviderStorageDiscoveryWorkspaceID returns a new ProviderStorageDiscoveryWorkspaceId struct
func NewProviderStorageDiscoveryWorkspaceID(subscriptionId string, resourceGroupName string, storageDiscoveryWorkspaceName string) ProviderStorageDiscoveryWorkspaceId {
	return ProviderStorageDiscoveryWorkspaceId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		StorageDiscoveryWorkspaceName: storageDiscoveryWorkspaceName,
	}
}

// ParseProviderStorageDiscoveryWorkspaceID parses 'input' into a ProviderStorageDiscoveryWorkspaceId
func ParseProviderStorageDiscoveryWorkspaceID(input string) (*ProviderStorageDiscoveryWorkspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderStorageDiscoveryWorkspaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderStorageDiscoveryWorkspaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviderStorageDiscoveryWorkspaceIDInsensitively parses 'input' case-insensitively into a ProviderStorageDiscoveryWorkspaceId
// note: this method should only be used for API response data and not user input
func ParseProviderStorageDiscoveryWorkspaceIDInsensitively(input string) (*ProviderStorageDiscoveryWorkspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderStorageDiscoveryWorkspaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderStorageDiscoveryWorkspaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProviderStorageDiscoveryWorkspaceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageDiscoveryWorkspaceName, ok = input.Parsed["storageDiscoveryWorkspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageDiscoveryWorkspaceName", input)
	}

	return nil
}

// ValidateProviderStorageDiscoveryWorkspaceID checks that 'input' can be parsed as a Provider Storage Discovery Workspace ID
func ValidateProviderStorageDiscoveryWorkspaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderStorageDiscoveryWorkspaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Storage Discovery Workspace ID
func (id ProviderStorageDiscoveryWorkspaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageDiscoveryWorkspaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Storage Discovery Workspace ID
func (id ProviderStorageDiscoveryWorkspaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageDiscovery", "Microsoft.StorageDiscovery", "Microsoft.StorageDiscovery"),
		resourceids.StaticSegment("staticStorageDiscoveryWorkspaces", "storageDiscoveryWorkspaces", "storageDiscoveryWorkspaces"),
		resourceids.UserSpecifiedSegment("storageDiscoveryWorkspaceName", "storageDiscoveryWorkspaceName"),
	}
}

// String returns a human-readable description of this Provider Storage Discovery Workspace ID
func (id ProviderStorageDiscoveryWorkspaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Discovery Workspace Name: %q", id.StorageDiscoveryWorkspaceName),
	}
	return fmt.Sprintf("Provider Storage Discovery Workspace (%s)", strings.Join(components, "\n"))
}
