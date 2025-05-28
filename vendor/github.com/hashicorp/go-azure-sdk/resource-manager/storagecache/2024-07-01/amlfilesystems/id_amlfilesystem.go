package amlfilesystems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AmlFilesystemId{})
}

var _ resourceids.ResourceId = &AmlFilesystemId{}

// AmlFilesystemId is a struct representing the Resource ID for a Aml Filesystem
type AmlFilesystemId struct {
	SubscriptionId    string
	ResourceGroupName string
	AmlFilesystemName string
}

// NewAmlFilesystemID returns a new AmlFilesystemId struct
func NewAmlFilesystemID(subscriptionId string, resourceGroupName string, amlFilesystemName string) AmlFilesystemId {
	return AmlFilesystemId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AmlFilesystemName: amlFilesystemName,
	}
}

// ParseAmlFilesystemID parses 'input' into a AmlFilesystemId
func ParseAmlFilesystemID(input string) (*AmlFilesystemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AmlFilesystemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AmlFilesystemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAmlFilesystemIDInsensitively parses 'input' case-insensitively into a AmlFilesystemId
// note: this method should only be used for API response data and not user input
func ParseAmlFilesystemIDInsensitively(input string) (*AmlFilesystemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AmlFilesystemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AmlFilesystemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AmlFilesystemId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateAmlFilesystemID checks that 'input' can be parsed as a Aml Filesystem ID
func ValidateAmlFilesystemID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAmlFilesystemID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Aml Filesystem ID
func (id AmlFilesystemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageCache/amlFilesystems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AmlFilesystemName)
}

// Segments returns a slice of Resource ID Segments which comprise this Aml Filesystem ID
func (id AmlFilesystemId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageCache", "Microsoft.StorageCache", "Microsoft.StorageCache"),
		resourceids.StaticSegment("staticAmlFilesystems", "amlFilesystems", "amlFilesystems"),
		resourceids.UserSpecifiedSegment("amlFilesystemName", "amlFilesystemName"),
	}
}

// String returns a human-readable description of this Aml Filesystem ID
func (id AmlFilesystemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Aml Filesystem Name: %q", id.AmlFilesystemName),
	}
	return fmt.Sprintf("Aml Filesystem (%s)", strings.Join(components, "\n"))
}
