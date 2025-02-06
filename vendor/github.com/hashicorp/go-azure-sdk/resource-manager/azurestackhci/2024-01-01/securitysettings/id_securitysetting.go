package securitysettings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SecuritySettingId{})
}

var _ resourceids.ResourceId = &SecuritySettingId{}

// SecuritySettingId is a struct representing the Resource ID for a Security Setting
type SecuritySettingId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ClusterName         string
	SecuritySettingName string
}

// NewSecuritySettingID returns a new SecuritySettingId struct
func NewSecuritySettingID(subscriptionId string, resourceGroupName string, clusterName string, securitySettingName string) SecuritySettingId {
	return SecuritySettingId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ClusterName:         clusterName,
		SecuritySettingName: securitySettingName,
	}
}

// ParseSecuritySettingID parses 'input' into a SecuritySettingId
func ParseSecuritySettingID(input string) (*SecuritySettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecuritySettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecuritySettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSecuritySettingIDInsensitively parses 'input' case-insensitively into a SecuritySettingId
// note: this method should only be used for API response data and not user input
func ParseSecuritySettingIDInsensitively(input string) (*SecuritySettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecuritySettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecuritySettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SecuritySettingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ClusterName, ok = input.Parsed["clusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clusterName", input)
	}

	if id.SecuritySettingName, ok = input.Parsed["securitySettingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "securitySettingName", input)
	}

	return nil
}

// ValidateSecuritySettingID checks that 'input' can be parsed as a Security Setting ID
func ValidateSecuritySettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecuritySettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Security Setting ID
func (id SecuritySettingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/clusters/%s/securitySettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.SecuritySettingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Security Setting ID
func (id SecuritySettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticSecuritySettings", "securitySettings", "securitySettings"),
		resourceids.UserSpecifiedSegment("securitySettingName", "securitySettingName"),
	}
}

// String returns a human-readable description of this Security Setting ID
func (id SecuritySettingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Security Setting Name: %q", id.SecuritySettingName),
	}
	return fmt.Sprintf("Security Setting (%s)", strings.Join(components, "\n"))
}
