// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = HyperVSiteRunAsAccountId{}

// HyperVSiteRunAsAccountId is a struct representing the Resource ID for a Hyper V Site Run As Account
type HyperVSiteRunAsAccountId struct {
	SubscriptionId    string
	ResourceGroupName string
	HyperVSiteName    string
	RunAsAccountName  string
}

// NewHyperVSiteRunAsAccountID returns a new HyperVSiteRunAsAccountId struct
func NewHyperVSiteRunAsAccountID(subscriptionId string, resourceGroupName string, hyperVSiteName string, runAsAccountName string) HyperVSiteRunAsAccountId {
	return HyperVSiteRunAsAccountId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HyperVSiteName:    hyperVSiteName,
		RunAsAccountName:  runAsAccountName,
	}
}

// ParseHyperVSiteRunAsAccountID parses 'input' into a HyperVSiteRunAsAccountId
func ParseHyperVSiteRunAsAccountID(input string) (*HyperVSiteRunAsAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(HyperVSiteRunAsAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HyperVSiteRunAsAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HyperVSiteName, ok = parsed.Parsed["hyperVSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hyperVSiteName", *parsed)
	}

	if id.RunAsAccountName, ok = parsed.Parsed["runAsAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "runAsAccountName", *parsed)
	}

	return &id, nil
}

// ParseHyperVSiteRunAsAccountIDInsensitively parses 'input' case-insensitively into a HyperVSiteRunAsAccountId
// note: this method should only be used for API response data and not user input
func ParseHyperVSiteRunAsAccountIDInsensitively(input string) (*HyperVSiteRunAsAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(HyperVSiteRunAsAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HyperVSiteRunAsAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HyperVSiteName, ok = parsed.Parsed["hyperVSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hyperVSiteName", *parsed)
	}

	if id.RunAsAccountName, ok = parsed.Parsed["runAsAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "runAsAccountName", *parsed)
	}

	return &id, nil
}

// ValidateHyperVSiteRunAsAccountID checks that 'input' can be parsed as a Hyper V Site Run As Account ID
func ValidateHyperVSiteRunAsAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHyperVSiteRunAsAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hyper V Site Run As Account ID
func (id HyperVSiteRunAsAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OffAzure/hyperVSites/%s/runAsAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HyperVSiteName, id.RunAsAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hyper V Site Run As Account ID
func (id HyperVSiteRunAsAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOffAzure", "Microsoft.OffAzure", "Microsoft.OffAzure"),
		resourceids.StaticSegment("staticHyperVSites", "hyperVSites", "hyperVSites"),
		resourceids.UserSpecifiedSegment("hyperVSiteName", "hyperVSiteValue"),
		resourceids.StaticSegment("staticRunAsAccounts", "runAsAccounts", "runAsAccounts"),
		resourceids.UserSpecifiedSegment("runAsAccountName", "runAsAccountNameValue"),
	}
}

// String returns a human-readable description of this Hyper V Site RunAsAccount ID
func (id HyperVSiteRunAsAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Hyper V Site Name: %q", id.HyperVSiteName),
		fmt.Sprintf("Run As Account Name: %q", id.RunAsAccountName),
	}
	return fmt.Sprintf("Hyper V Site Run As Account (%s)", strings.Join(components, "\n"))
}
