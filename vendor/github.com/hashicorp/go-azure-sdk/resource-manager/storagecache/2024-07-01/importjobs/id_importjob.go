package importjobs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ImportJobId{})
}

var _ resourceids.ResourceId = &ImportJobId{}

// ImportJobId is a struct representing the Resource ID for a Import Job
type ImportJobId struct {
	SubscriptionId    string
	ResourceGroupName string
	AmlFilesystemName string
	ImportJobName     string
}

// NewImportJobID returns a new ImportJobId struct
func NewImportJobID(subscriptionId string, resourceGroupName string, amlFilesystemName string, importJobName string) ImportJobId {
	return ImportJobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AmlFilesystemName: amlFilesystemName,
		ImportJobName:     importJobName,
	}
}

// ParseImportJobID parses 'input' into a ImportJobId
func ParseImportJobID(input string) (*ImportJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImportJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImportJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseImportJobIDInsensitively parses 'input' case-insensitively into a ImportJobId
// note: this method should only be used for API response data and not user input
func ParseImportJobIDInsensitively(input string) (*ImportJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImportJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImportJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ImportJobId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AmlFilesystemName, ok = input.Parsed["amlFilesystemName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "amlFilesystemName", input)
	}

	if id.ImportJobName, ok = input.Parsed["importJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "importJobName", input)
	}

	return nil
}

// ValidateImportJobID checks that 'input' can be parsed as a Import Job ID
func ValidateImportJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseImportJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Import Job ID
func (id ImportJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageCache/amlFilesystems/%s/importJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AmlFilesystemName, id.ImportJobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Import Job ID
func (id ImportJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageCache", "Microsoft.StorageCache", "Microsoft.StorageCache"),
		resourceids.StaticSegment("staticAmlFilesystems", "amlFilesystems", "amlFilesystems"),
		resourceids.UserSpecifiedSegment("amlFilesystemName", "amlFilesystemName"),
		resourceids.StaticSegment("staticImportJobs", "importJobs", "importJobs"),
		resourceids.UserSpecifiedSegment("importJobName", "importJobName"),
	}
}

// String returns a human-readable description of this Import Job ID
func (id ImportJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Aml Filesystem Name: %q", id.AmlFilesystemName),
		fmt.Sprintf("Import Job Name: %q", id.ImportJobName),
	}
	return fmt.Sprintf("Import Job (%s)", strings.Join(components, "\n"))
}
