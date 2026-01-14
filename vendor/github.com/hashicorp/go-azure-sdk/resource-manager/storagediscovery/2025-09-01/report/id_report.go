package report

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReportId{})
}

var _ resourceids.ResourceId = &ReportId{}

// ReportId is a struct representing the Resource ID for a Report
type ReportId struct {
	SubscriptionId                string
	ResourceGroupName             string
	StorageDiscoveryWorkspaceName string
	ReportName                    string
}

// NewReportID returns a new ReportId struct
func NewReportID(subscriptionId string, resourceGroupName string, storageDiscoveryWorkspaceName string, reportName string) ReportId {
	return ReportId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		StorageDiscoveryWorkspaceName: storageDiscoveryWorkspaceName,
		ReportName:                    reportName,
	}
}

// ParseReportID parses 'input' into a ReportId
func ParseReportID(input string) (*ReportId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReportId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReportId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReportIDInsensitively parses 'input' case-insensitively into a ReportId
// note: this method should only be used for API response data and not user input
func ParseReportIDInsensitively(input string) (*ReportId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReportId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReportId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReportId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageDiscoveryWorkspaceName, ok = input.Parsed["storageDiscoveryWorkspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageDiscoveryWorkspaceName", input)
	}

	if id.ReportName, ok = input.Parsed["reportName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "reportName", input)
	}

	return nil
}

// ValidateReportID checks that 'input' can be parsed as a Report ID
func ValidateReportID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReportID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Report ID
func (id ReportId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/%s/reports/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageDiscoveryWorkspaceName, id.ReportName)
}

// Segments returns a slice of Resource ID Segments which comprise this Report ID
func (id ReportId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageDiscovery", "Microsoft.StorageDiscovery", "Microsoft.StorageDiscovery"),
		resourceids.StaticSegment("staticStorageDiscoveryWorkspaces", "storageDiscoveryWorkspaces", "storageDiscoveryWorkspaces"),
		resourceids.UserSpecifiedSegment("storageDiscoveryWorkspaceName", "storageDiscoveryWorkspaceName"),
		resourceids.StaticSegment("staticReports", "reports", "reports"),
		resourceids.UserSpecifiedSegment("reportName", "reportName"),
	}
}

// String returns a human-readable description of this Report ID
func (id ReportId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Discovery Workspace Name: %q", id.StorageDiscoveryWorkspaceName),
		fmt.Sprintf("Report Name: %q", id.ReportName),
	}
	return fmt.Sprintf("Report (%s)", strings.Join(components, "\n"))
}
