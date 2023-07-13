// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = HyperVSiteJobId{}

// HyperVSiteJobId is a struct representing the Resource ID for a Hyper V Site Job
type HyperVSiteJobId struct {
	SubscriptionId    string
	ResourceGroupName string
	HyperVSiteName    string
	JobName           string
}

// NewHyperVSiteJobID returns a new HyperVSiteJobId struct
func NewHyperVSiteJobID(subscriptionId string, resourceGroupName string, hyperVSiteName string, jobName string) HyperVSiteJobId {
	return HyperVSiteJobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HyperVSiteName:    hyperVSiteName,
		JobName:           jobName,
	}
}

// ParseHyperVSiteJobID parses 'input' into a HyperVSiteJobId
func ParseHyperVSiteJobID(input string) (*HyperVSiteJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(HyperVSiteJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HyperVSiteJobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HyperVSiteName, ok = parsed.Parsed["hyperVSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hyperVSiteName", *parsed)
	}

	if id.JobName, ok = parsed.Parsed["jobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobName", *parsed)
	}

	return &id, nil
}

// ParseHyperVSiteJobIDInsensitively parses 'input' case-insensitively into a HyperVSiteJobId
// note: this method should only be used for API response data and not user input
func ParseHyperVSiteJobIDInsensitively(input string) (*HyperVSiteJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(HyperVSiteJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HyperVSiteJobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HyperVSiteName, ok = parsed.Parsed["hyperVSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hyperVSiteName", *parsed)
	}

	if id.JobName, ok = parsed.Parsed["jobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobName", *parsed)
	}

	return &id, nil
}

// ValidateHyperVSiteJobID checks that 'input' can be parsed as a Hyper V Site Job ID
func ValidateHyperVSiteJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHyperVSiteJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hyper V Site Job ID
func (id HyperVSiteJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OffAzure/hyperVSites/%s/jobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HyperVSiteName, id.JobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hyper V Site Job ID
func (id HyperVSiteJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOffAzure", "Microsoft.OffAzure", "Microsoft.OffAzure"),
		resourceids.StaticSegment("staticHyperVSites", "hyperVSites", "hyperVSites"),
		resourceids.UserSpecifiedSegment("hyperVSiteName", "hyperVSiteValue"),
		resourceids.StaticSegment("staticJobs", "jobs", "jobs"),
		resourceids.UserSpecifiedSegment("jobName", "jobNameValue"),
	}
}

// String returns a human-readable description of this Hyper V Site Job ID
func (id HyperVSiteJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Hyper V Site Name: %q", id.HyperVSiteName),
		fmt.Sprintf("Job Name: %q", id.JobName),
	}
	return fmt.Sprintf("Hyper V Site Job (%s)", strings.Join(components, "\n"))
}
