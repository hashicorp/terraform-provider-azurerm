package archiveversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ArchiveId{})
}

var _ resourceids.ResourceId = &ArchiveId{}

// ArchiveId is a struct representing the Resource ID for a Archive
type ArchiveId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	PackageName       string
	ArchiveName       string
}

// NewArchiveID returns a new ArchiveId struct
func NewArchiveID(subscriptionId string, resourceGroupName string, registryName string, packageName string, archiveName string) ArchiveId {
	return ArchiveId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		PackageName:       packageName,
		ArchiveName:       archiveName,
	}
}

// ParseArchiveID parses 'input' into a ArchiveId
func ParseArchiveID(input string) (*ArchiveId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ArchiveId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ArchiveId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseArchiveIDInsensitively parses 'input' case-insensitively into a ArchiveId
// note: this method should only be used for API response data and not user input
func ParseArchiveIDInsensitively(input string) (*ArchiveId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ArchiveId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ArchiveId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ArchiveId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RegistryName, ok = input.Parsed["registryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "registryName", input)
	}

	if id.PackageName, ok = input.Parsed["packageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "packageName", input)
	}

	if id.ArchiveName, ok = input.Parsed["archiveName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "archiveName", input)
	}

	return nil
}

// ValidateArchiveID checks that 'input' can be parsed as a Archive ID
func ValidateArchiveID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseArchiveID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Archive ID
func (id ArchiveId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/packages/%s/archives/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.PackageName, id.ArchiveName)
}

// Segments returns a slice of Resource ID Segments which comprise this Archive ID
func (id ArchiveId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticPackages", "packages", "packages"),
		resourceids.UserSpecifiedSegment("packageName", "packageName"),
		resourceids.StaticSegment("staticArchives", "archives", "archives"),
		resourceids.UserSpecifiedSegment("archiveName", "archiveName"),
	}
}

// String returns a human-readable description of this Archive ID
func (id ArchiveId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Package Name: %q", id.PackageName),
		fmt.Sprintf("Archive Name: %q", id.ArchiveName),
	}
	return fmt.Sprintf("Archive (%s)", strings.Join(components, "\n"))
}
