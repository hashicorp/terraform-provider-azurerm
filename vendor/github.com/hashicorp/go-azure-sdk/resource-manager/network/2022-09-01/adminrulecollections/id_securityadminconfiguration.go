package adminrulecollections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SecurityAdminConfigurationId{}

// SecurityAdminConfigurationId is a struct representing the Resource ID for a Security Admin Configuration
type SecurityAdminConfigurationId struct {
	SubscriptionId                 string
	ResourceGroupName              string
	NetworkManagerName             string
	SecurityAdminConfigurationName string
}

// NewSecurityAdminConfigurationID returns a new SecurityAdminConfigurationId struct
func NewSecurityAdminConfigurationID(subscriptionId string, resourceGroupName string, networkManagerName string, securityAdminConfigurationName string) SecurityAdminConfigurationId {
	return SecurityAdminConfigurationId{
		SubscriptionId:                 subscriptionId,
		ResourceGroupName:              resourceGroupName,
		NetworkManagerName:             networkManagerName,
		SecurityAdminConfigurationName: securityAdminConfigurationName,
	}
}

// ParseSecurityAdminConfigurationID parses 'input' into a SecurityAdminConfigurationId
func ParseSecurityAdminConfigurationID(input string) (*SecurityAdminConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(SecurityAdminConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SecurityAdminConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.SecurityAdminConfigurationName, ok = parsed.Parsed["securityAdminConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "securityAdminConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseSecurityAdminConfigurationIDInsensitively parses 'input' case-insensitively into a SecurityAdminConfigurationId
// note: this method should only be used for API response data and not user input
func ParseSecurityAdminConfigurationIDInsensitively(input string) (*SecurityAdminConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(SecurityAdminConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SecurityAdminConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.SecurityAdminConfigurationName, ok = parsed.Parsed["securityAdminConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "securityAdminConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateSecurityAdminConfigurationID checks that 'input' can be parsed as a Security Admin Configuration ID
func ValidateSecurityAdminConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecurityAdminConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Security Admin Configuration ID
func (id SecurityAdminConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/securityAdminConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.SecurityAdminConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Security Admin Configuration ID
func (id SecurityAdminConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerValue"),
		resourceids.StaticSegment("staticSecurityAdminConfigurations", "securityAdminConfigurations", "securityAdminConfigurations"),
		resourceids.UserSpecifiedSegment("securityAdminConfigurationName", "securityAdminConfigurationValue"),
	}
}

// String returns a human-readable description of this Security Admin Configuration ID
func (id SecurityAdminConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Security Admin Configuration Name: %q", id.SecurityAdminConfigurationName),
	}
	return fmt.Sprintf("Security Admin Configuration (%s)", strings.Join(components, "\n"))
}
