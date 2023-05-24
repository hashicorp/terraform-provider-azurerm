package metadata

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = MetadataId{}

// MetadataId is a struct representing the Resource ID for a Metadata
type MetadataId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	MetadataName      string
}

// NewMetadataID returns a new MetadataId struct
func NewMetadataID(subscriptionId string, resourceGroupName string, workspaceName string, metadataName string) MetadataId {
	return MetadataId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		MetadataName:      metadataName,
	}
}

// ParseMetadataID parses 'input' into a MetadataId
func ParseMetadataID(input string) (*MetadataId, error) {
	parser := resourceids.NewParserFromResourceIdType(MetadataId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MetadataId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.MetadataName, ok = parsed.Parsed["metadataName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "metadataName", *parsed)
	}

	return &id, nil
}

// ParseMetadataIDInsensitively parses 'input' case-insensitively into a MetadataId
// note: this method should only be used for API response data and not user input
func ParseMetadataIDInsensitively(input string) (*MetadataId, error) {
	parser := resourceids.NewParserFromResourceIdType(MetadataId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MetadataId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.MetadataName, ok = parsed.Parsed["metadataName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "metadataName", *parsed)
	}

	return &id, nil
}

// ValidateMetadataID checks that 'input' can be parsed as a Metadata ID
func ValidateMetadataID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMetadataID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Metadata ID
func (id MetadataId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/metadata/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.MetadataName)
}

// Segments returns a slice of Resource ID Segments which comprise this Metadata ID
func (id MetadataId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurityInsights", "Microsoft.SecurityInsights", "Microsoft.SecurityInsights"),
		resourceids.StaticSegment("staticMetadata", "metadata", "metadata"),
		resourceids.UserSpecifiedSegment("metadataName", "metadataValue"),
	}
}

// String returns a human-readable description of this Metadata ID
func (id MetadataId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Metadata Name: %q", id.MetadataName),
	}
	return fmt.Sprintf("Metadata (%s)", strings.Join(components, "\n"))
}
