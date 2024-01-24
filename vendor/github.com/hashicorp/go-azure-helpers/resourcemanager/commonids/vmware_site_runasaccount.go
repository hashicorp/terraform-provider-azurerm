// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &VMwareSiteRunAsAccountId{}

// VMwareSiteRunAsAccountId is a struct representing the Resource ID for a VMware Site Run As Account
type VMwareSiteRunAsAccountId struct {
	SubscriptionId    string
	ResourceGroupName string
	VMwareSiteName    string
	RunAsAccountName  string
}

// NewVMwareSiteRunAsAccountID returns a new VMwareSiteRunAsAccountId struct
func NewVMwareSiteRunAsAccountID(subscriptionId string, resourceGroupName string, vmwareSiteName string, runAsAccountName string) VMwareSiteRunAsAccountId {
	return VMwareSiteRunAsAccountId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VMwareSiteName:    vmwareSiteName,
		RunAsAccountName:  runAsAccountName,
	}
}

// ParseVMwareSiteRunAsAccountID parses 'input' into a VMwareSiteRunAsAccountId
func ParseVMwareSiteRunAsAccountID(input string) (*VMwareSiteRunAsAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VMwareSiteRunAsAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMwareSiteRunAsAccountId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}
	return &id, nil
}

// ParseVMwareSiteRunAsAccountIDInsensitively parses 'input' case-insensitively into a VMwareSiteRunAsAccountId
// note: this method should only be used for API response data and not user input
func ParseVMwareSiteRunAsAccountIDInsensitively(input string) (*VMwareSiteRunAsAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VMwareSiteRunAsAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMwareSiteRunAsAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VMwareSiteRunAsAccountId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VMwareSiteName, ok = input.Parsed["vmwareSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vmwareSiteName", input)
	}

	if id.RunAsAccountName, ok = input.Parsed["runAsAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "runAsAccountName", input)
	}

	return nil
}

// ValidateVMwareSiteRunAsAccountID checks that 'input' can be parsed as a VMware Site Run As Account ID
func ValidateVMwareSiteRunAsAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVMwareSiteRunAsAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted VMware Site Run As Account ID
func (id VMwareSiteRunAsAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OffAzure/vmwareSites/%s/runAsAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VMwareSiteName, id.RunAsAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this VMware Site Run As Account ID
func (id VMwareSiteRunAsAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOffAzure", "Microsoft.OffAzure", "Microsoft.OffAzure"),
		resourceids.StaticSegment("staticVMwareSites", "vmwareSites", "vmwareSites"),
		resourceids.UserSpecifiedSegment("vmwareSiteName", "vmwareSiteValue"),
		resourceids.StaticSegment("staticRunAsAccounts", "runAsAccounts", "runAsAccounts"),
		resourceids.UserSpecifiedSegment("runAsAccountName", "runAsAccountNameValue"),
	}
}

// String returns a human-readable description of this VMware Site RunAsAccount ID
func (id VMwareSiteRunAsAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("VMware Site Name: %q", id.VMwareSiteName),
		fmt.Sprintf("Run As Account Name: %q", id.RunAsAccountName),
	}
	return fmt.Sprintf("VMware Site Run As Account (%s)", strings.Join(components, "\n"))
}
