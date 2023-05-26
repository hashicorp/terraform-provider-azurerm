package securitymlanalyticssettings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SecurityMLAnalyticsSettingId{}

// SecurityMLAnalyticsSettingId is a struct representing the Resource ID for a Security M L Analytics Setting
type SecurityMLAnalyticsSettingId struct {
	SubscriptionId                 string
	ResourceGroupName              string
	WorkspaceName                  string
	SecurityMLAnalyticsSettingName string
}

// NewSecurityMLAnalyticsSettingID returns a new SecurityMLAnalyticsSettingId struct
func NewSecurityMLAnalyticsSettingID(subscriptionId string, resourceGroupName string, workspaceName string, securityMLAnalyticsSettingName string) SecurityMLAnalyticsSettingId {
	return SecurityMLAnalyticsSettingId{
		SubscriptionId:                 subscriptionId,
		ResourceGroupName:              resourceGroupName,
		WorkspaceName:                  workspaceName,
		SecurityMLAnalyticsSettingName: securityMLAnalyticsSettingName,
	}
}

// ParseSecurityMLAnalyticsSettingID parses 'input' into a SecurityMLAnalyticsSettingId
func ParseSecurityMLAnalyticsSettingID(input string) (*SecurityMLAnalyticsSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(SecurityMLAnalyticsSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SecurityMLAnalyticsSettingId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.SecurityMLAnalyticsSettingName, ok = parsed.Parsed["securityMLAnalyticsSettingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "securityMLAnalyticsSettingName", *parsed)
	}

	return &id, nil
}

// ParseSecurityMLAnalyticsSettingIDInsensitively parses 'input' case-insensitively into a SecurityMLAnalyticsSettingId
// note: this method should only be used for API response data and not user input
func ParseSecurityMLAnalyticsSettingIDInsensitively(input string) (*SecurityMLAnalyticsSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(SecurityMLAnalyticsSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SecurityMLAnalyticsSettingId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.SecurityMLAnalyticsSettingName, ok = parsed.Parsed["securityMLAnalyticsSettingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "securityMLAnalyticsSettingName", *parsed)
	}

	return &id, nil
}

// ValidateSecurityMLAnalyticsSettingID checks that 'input' can be parsed as a Security M L Analytics Setting ID
func ValidateSecurityMLAnalyticsSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecurityMLAnalyticsSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Security M L Analytics Setting ID
func (id SecurityMLAnalyticsSettingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/securityMLAnalyticsSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.SecurityMLAnalyticsSettingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Security M L Analytics Setting ID
func (id SecurityMLAnalyticsSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurityInsights", "Microsoft.SecurityInsights", "Microsoft.SecurityInsights"),
		resourceids.StaticSegment("staticSecurityMLAnalyticsSettings", "securityMLAnalyticsSettings", "securityMLAnalyticsSettings"),
		resourceids.UserSpecifiedSegment("securityMLAnalyticsSettingName", "securityMLAnalyticsSettingValue"),
	}
}

// String returns a human-readable description of this Security M L Analytics Setting ID
func (id SecurityMLAnalyticsSettingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Security M L Analytics Setting Name: %q", id.SecurityMLAnalyticsSettingName),
	}
	return fmt.Sprintf("Security M L Analytics Setting (%s)", strings.Join(components, "\n"))
}
