package runasaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VMwareSiteId{})
}

var _ resourceids.ResourceId = &VMwareSiteId{}

// VMwareSiteId is a struct representing the Resource ID for a V Mware Site
type VMwareSiteId struct {
	SubscriptionId    string
	ResourceGroupName string
	VmwareSiteName    string
}

// NewVMwareSiteID returns a new VMwareSiteId struct
func NewVMwareSiteID(subscriptionId string, resourceGroupName string, vmwareSiteName string) VMwareSiteId {
	return VMwareSiteId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VmwareSiteName:    vmwareSiteName,
	}
}

// ParseVMwareSiteID parses 'input' into a VMwareSiteId
func ParseVMwareSiteID(input string) (*VMwareSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VMwareSiteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMwareSiteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVMwareSiteIDInsensitively parses 'input' case-insensitively into a VMwareSiteId
// note: this method should only be used for API response data and not user input
func ParseVMwareSiteIDInsensitively(input string) (*VMwareSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VMwareSiteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMwareSiteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VMwareSiteId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VmwareSiteName, ok = input.Parsed["vmwareSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vmwareSiteName", input)
	}

	return nil
}

// ValidateVMwareSiteID checks that 'input' can be parsed as a V Mware Site ID
func ValidateVMwareSiteID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVMwareSiteID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted V Mware Site ID
func (id VMwareSiteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OffAzure/vmwareSites/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VmwareSiteName)
}

// Segments returns a slice of Resource ID Segments which comprise this V Mware Site ID
func (id VMwareSiteId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOffAzure", "Microsoft.OffAzure", "Microsoft.OffAzure"),
		resourceids.StaticSegment("staticVmwareSites", "vmwareSites", "vmwareSites"),
		resourceids.UserSpecifiedSegment("vmwareSiteName", "vmwareSiteName"),
	}
}

// String returns a human-readable description of this V Mware Site ID
func (id VMwareSiteId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vmware Site Name: %q", id.VmwareSiteName),
	}
	return fmt.Sprintf("V Mware Site (%s)", strings.Join(components, "\n"))
}
