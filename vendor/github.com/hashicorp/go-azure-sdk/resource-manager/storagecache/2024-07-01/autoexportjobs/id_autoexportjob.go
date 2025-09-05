package autoexportjobs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutoExportJobId{})
}

var _ resourceids.ResourceId = &AutoExportJobId{}

// AutoExportJobId is a struct representing the Resource ID for a Auto Export Job
type AutoExportJobId struct {
	SubscriptionId    string
	ResourceGroupName string
	AmlFilesystemName string
	AutoExportJobName string
}

// NewAutoExportJobID returns a new AutoExportJobId struct
func NewAutoExportJobID(subscriptionId string, resourceGroupName string, amlFilesystemName string, autoExportJobName string) AutoExportJobId {
	return AutoExportJobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AmlFilesystemName: amlFilesystemName,
		AutoExportJobName: autoExportJobName,
	}
}

// ParseAutoExportJobID parses 'input' into a AutoExportJobId
func ParseAutoExportJobID(input string) (*AutoExportJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutoExportJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutoExportJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutoExportJobIDInsensitively parses 'input' case-insensitively into a AutoExportJobId
// note: this method should only be used for API response data and not user input
func ParseAutoExportJobIDInsensitively(input string) (*AutoExportJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutoExportJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutoExportJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutoExportJobId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AutoExportJobName, ok = input.Parsed["autoExportJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "autoExportJobName", input)
	}

	return nil
}

// ValidateAutoExportJobID checks that 'input' can be parsed as a Auto Export Job ID
func ValidateAutoExportJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutoExportJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Auto Export Job ID
func (id AutoExportJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageCache/amlFilesystems/%s/autoExportJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AmlFilesystemName, id.AutoExportJobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Auto Export Job ID
func (id AutoExportJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageCache", "Microsoft.StorageCache", "Microsoft.StorageCache"),
		resourceids.StaticSegment("staticAmlFilesystems", "amlFilesystems", "amlFilesystems"),
		resourceids.UserSpecifiedSegment("amlFilesystemName", "amlFilesystemName"),
		resourceids.StaticSegment("staticAutoExportJobs", "autoExportJobs", "autoExportJobs"),
		resourceids.UserSpecifiedSegment("autoExportJobName", "autoExportJobName"),
	}
}

// String returns a human-readable description of this Auto Export Job ID
func (id AutoExportJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Aml Filesystem Name: %q", id.AmlFilesystemName),
		fmt.Sprintf("Auto Export Job Name: %q", id.AutoExportJobName),
	}
	return fmt.Sprintf("Auto Export Job (%s)", strings.Join(components, "\n"))
}
