package filesystems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FileSystemId{})
}

var _ resourceids.ResourceId = &FileSystemId{}

// FileSystemId is a struct representing the Resource ID for a File System
type FileSystemId struct {
	SubscriptionId    string
	ResourceGroupName string
	FileSystemName    string
}

// NewFileSystemID returns a new FileSystemId struct
func NewFileSystemID(subscriptionId string, resourceGroupName string, fileSystemName string) FileSystemId {
	return FileSystemId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FileSystemName:    fileSystemName,
	}
}

// ParseFileSystemID parses 'input' into a FileSystemId
func ParseFileSystemID(input string) (*FileSystemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FileSystemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FileSystemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFileSystemIDInsensitively parses 'input' case-insensitively into a FileSystemId
// note: this method should only be used for API response data and not user input
func ParseFileSystemIDInsensitively(input string) (*FileSystemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FileSystemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FileSystemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FileSystemId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FileSystemName, ok = input.Parsed["fileSystemName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "fileSystemName", input)
	}

	return nil
}

// ValidateFileSystemID checks that 'input' can be parsed as a File System ID
func ValidateFileSystemID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFileSystemID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted File System ID
func (id FileSystemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Qumulo.Storage/fileSystems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FileSystemName)
}

// Segments returns a slice of Resource ID Segments which comprise this File System ID
func (id FileSystemId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticQumuloStorage", "Qumulo.Storage", "Qumulo.Storage"),
		resourceids.StaticSegment("staticFileSystems", "fileSystems", "fileSystems"),
		resourceids.UserSpecifiedSegment("fileSystemName", "fileSystemName"),
	}
}

// String returns a human-readable description of this File System ID
func (id FileSystemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("File System Name: %q", id.FileSystemName),
	}
	return fmt.Sprintf("File System (%s)", strings.Join(components, "\n"))
}
