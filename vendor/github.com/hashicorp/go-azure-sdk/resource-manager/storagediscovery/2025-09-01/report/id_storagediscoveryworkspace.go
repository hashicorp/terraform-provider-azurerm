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
	recaser.RegisterResourceId(&StorageDiscoveryWorkspaceId{})
}

var _ resourceids.ResourceId = &StorageDiscoveryWorkspaceId{}

// StorageDiscoveryWorkspaceId is a struct representing the Resource ID for a Storage Discovery Workspace
type StorageDiscoveryWorkspaceId struct {
	SubscriptionId                string
	StorageDiscoveryWorkspaceName string
}

// NewStorageDiscoveryWorkspaceID returns a new StorageDiscoveryWorkspaceId struct
func NewStorageDiscoveryWorkspaceID(subscriptionId string, storageDiscoveryWorkspaceName string) StorageDiscoveryWorkspaceId {
	return StorageDiscoveryWorkspaceId{
		SubscriptionId:                subscriptionId,
		StorageDiscoveryWorkspaceName: storageDiscoveryWorkspaceName,
	}
}

// ParseStorageDiscoveryWorkspaceID parses 'input' into a StorageDiscoveryWorkspaceId
func ParseStorageDiscoveryWorkspaceID(input string) (*StorageDiscoveryWorkspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageDiscoveryWorkspaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageDiscoveryWorkspaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageDiscoveryWorkspaceIDInsensitively parses 'input' case-insensitively into a StorageDiscoveryWorkspaceId
// note: this method should only be used for API response data and not user input
func ParseStorageDiscoveryWorkspaceIDInsensitively(input string) (*StorageDiscoveryWorkspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageDiscoveryWorkspaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageDiscoveryWorkspaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageDiscoveryWorkspaceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.StorageDiscoveryWorkspaceName, ok = input.Parsed["storageDiscoveryWorkspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageDiscoveryWorkspaceName", input)
	}

	return nil
}

// ValidateStorageDiscoveryWorkspaceID checks that 'input' can be parsed as a Storage Discovery Workspace ID
func ValidateStorageDiscoveryWorkspaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageDiscoveryWorkspaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Discovery Workspace ID
func (id StorageDiscoveryWorkspaceId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.StorageDiscoveryWorkspaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Discovery Workspace ID
func (id StorageDiscoveryWorkspaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageDiscovery", "Microsoft.StorageDiscovery", "Microsoft.StorageDiscovery"),
		resourceids.StaticSegment("staticStorageDiscoveryWorkspaces", "storageDiscoveryWorkspaces", "storageDiscoveryWorkspaces"),
		resourceids.UserSpecifiedSegment("storageDiscoveryWorkspaceName", "storageDiscoveryWorkspaceName"),
	}
}

// String returns a human-readable description of this Storage Discovery Workspace ID
func (id StorageDiscoveryWorkspaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Storage Discovery Workspace Name: %q", id.StorageDiscoveryWorkspaceName),
	}
	return fmt.Sprintf("Storage Discovery Workspace (%s)", strings.Join(components, "\n"))
}
