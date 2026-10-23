package autoimportjobs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutoImportJobId{})
}

var _ resourceids.ResourceId = &AutoImportJobId{}

// AutoImportJobId is a struct representing the Resource ID for a Auto Import Job
type AutoImportJobId struct {
	SubscriptionId    string
	ResourceGroupName string
	AmlFilesystemName string
	AutoImportJobName string
}

// NewAutoImportJobID returns a new AutoImportJobId struct
func NewAutoImportJobID(subscriptionId string, resourceGroupName string, amlFilesystemName string, autoImportJobName string) AutoImportJobId {
	return AutoImportJobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AmlFilesystemName: amlFilesystemName,
		AutoImportJobName: autoImportJobName,
	}
}

// ParseAutoImportJobID parses 'input' into a AutoImportJobId
func ParseAutoImportJobID(input string) (*AutoImportJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutoImportJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutoImportJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutoImportJobIDInsensitively parses 'input' case-insensitively into a AutoImportJobId
// note: this method should only be used for API response data and not user input
func ParseAutoImportJobIDInsensitively(input string) (*AutoImportJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutoImportJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutoImportJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutoImportJobId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AutoImportJobName, ok = input.Parsed["autoImportJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "autoImportJobName", input)
	}

	return nil
}

// ValidateAutoImportJobID checks that 'input' can be parsed as a Auto Import Job ID
func ValidateAutoImportJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutoImportJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Auto Import Job ID
func (id AutoImportJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageCache/amlFilesystems/%s/autoImportJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AmlFilesystemName, id.AutoImportJobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Auto Import Job ID
func (id AutoImportJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageCache", "Microsoft.StorageCache", "Microsoft.StorageCache"),
		resourceids.StaticSegment("staticAmlFilesystems", "amlFilesystems", "amlFilesystems"),
		resourceids.UserSpecifiedSegment("amlFilesystemName", "amlFilesystemName"),
		resourceids.StaticSegment("staticAutoImportJobs", "autoImportJobs", "autoImportJobs"),
		resourceids.UserSpecifiedSegment("autoImportJobName", "autoImportJobName"),
	}
}

// String returns a human-readable description of this Auto Import Job ID
func (id AutoImportJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Aml Filesystem Name: %q", id.AmlFilesystemName),
		fmt.Sprintf("Auto Import Job Name: %q", id.AutoImportJobName),
	}
	return fmt.Sprintf("Auto Import Job (%s)", strings.Join(components, "\n"))
}
