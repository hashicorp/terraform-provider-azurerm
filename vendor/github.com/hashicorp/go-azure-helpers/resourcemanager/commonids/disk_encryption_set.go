// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &DiskEncryptionSetId{}

// DiskEncryptionSetId is a struct representing the Resource ID for a Disk Encryption Set
type DiskEncryptionSetId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DiskEncryptionSetName string
}

// NewDiskEncryptionSetID returns a new DiskEncryptionSetId struct
func NewDiskEncryptionSetID(subscriptionId string, resourceGroupName string, diskEncryptionSetName string) DiskEncryptionSetId {
	return DiskEncryptionSetId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DiskEncryptionSetName: diskEncryptionSetName,
	}
}

// ParseDiskEncryptionSetID parses 'input' into a DiskEncryptionSetId
func ParseDiskEncryptionSetID(input string) (*DiskEncryptionSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DiskEncryptionSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DiskEncryptionSetId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDiskEncryptionSetIDInsensitively parses 'input' case-insensitively into a DiskEncryptionSetId
// note: this method should only be used for API response data and not user input
func ParseDiskEncryptionSetIDInsensitively(input string) (*DiskEncryptionSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DiskEncryptionSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DiskEncryptionSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DiskEncryptionSetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DiskEncryptionSetName, ok = input.Parsed["diskEncryptionSetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "diskEncryptionSetName", input)
	}

	return nil
}

// ValidateDiskEncryptionSetID checks that 'input' can be parsed as a Disk Encryption Set ID
func ValidateDiskEncryptionSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDiskEncryptionSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Disk Encryption Set ID
func (id DiskEncryptionSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/diskEncryptionSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DiskEncryptionSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Disk Encryption Set ID
func (id DiskEncryptionSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticDiskEncryptionSets", "diskEncryptionSets", "diskEncryptionSets"),
		resourceids.UserSpecifiedSegment("diskEncryptionSetName", "diskEncryptionSetValue"),
	}
}

// String returns a human-readable description of this Disk Encryption Set ID
func (id DiskEncryptionSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Disk Encryption Set Name: %q", id.DiskEncryptionSetName),
	}
	return fmt.Sprintf("Disk Encryption Set (%s)", strings.Join(components, "\n"))
}
