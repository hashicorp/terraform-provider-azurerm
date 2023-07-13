// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VMwareSiteJobId{}

// VMwareSiteJobId is a struct representing the Resource ID for a VMware Site Job
type VMwareSiteJobId struct {
	SubscriptionId    string
	ResourceGroupName string
	VMwareSiteName    string
	JobName           string
}

// NewVMwareSiteJobID returns a new VMwareSiteJobId struct
func NewVMwareSiteJobID(subscriptionId string, resourceGroupName string, vmwareSiteName string, jobName string) VMwareSiteJobId {
	return VMwareSiteJobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VMwareSiteName:    vmwareSiteName,
		JobName:           jobName,
	}
}

// ParseVMwareSiteJobID parses 'input' into a VMwareSiteJobId
func ParseVMwareSiteJobID(input string) (*VMwareSiteJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(VMwareSiteJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VMwareSiteJobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VMwareSiteName, ok = parsed.Parsed["vmwareSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vmwareSiteName", *parsed)
	}

	if id.JobName, ok = parsed.Parsed["jobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobName", *parsed)
	}

	return &id, nil
}

// ParseVMwareSiteJobIDInsensitively parses 'input' case-insensitively into a VMwareSiteJobId
// note: this method should only be used for API response data and not user input
func ParseVMwareSiteJobIDInsensitively(input string) (*VMwareSiteJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(VMwareSiteJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VMwareSiteJobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VMwareSiteName, ok = parsed.Parsed["vmwareSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vmwareSiteName", *parsed)
	}

	if id.JobName, ok = parsed.Parsed["jobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobName", *parsed)
	}

	return &id, nil
}

// ValidateVMwareSiteJobID checks that 'input' can be parsed as a VMware Site Job ID
func ValidateVMwareSiteJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVMwareSiteJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted VMware Site Job ID
func (id VMwareSiteJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OffAzure/vmwareSites/%s/jobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VMwareSiteName, id.JobName)
}

// Segments returns a slice of Resource ID Segments which comprise this VMware Site Job ID
func (id VMwareSiteJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOffAzure", "Microsoft.OffAzure", "Microsoft.OffAzure"),
		resourceids.StaticSegment("staticVMwareSites", "vmwareSites", "vmwareSites"),
		resourceids.UserSpecifiedSegment("vmwareSiteName", "vmwareSiteNameValue"),
		resourceids.StaticSegment("staticJobs", "jobs", "jobs"),
		resourceids.UserSpecifiedSegment("jobName", "jobNameValue"),
	}
}

// String returns a human-readable description of this VMware Site Job ID
func (id VMwareSiteJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("VMware Site Name: %q", id.VMwareSiteName),
		fmt.Sprintf("Job Name: %q", id.JobName),
	}
	return fmt.Sprintf("VMware Site Job (%s)", strings.Join(components, "\n"))
}
