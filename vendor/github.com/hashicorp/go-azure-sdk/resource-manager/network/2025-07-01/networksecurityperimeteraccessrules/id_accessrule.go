package networksecurityperimeteraccessrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AccessRuleId{})
}

var _ resourceids.ResourceId = &AccessRuleId{}

// AccessRuleId is a struct representing the Resource ID for a Access Rule
type AccessRuleId struct {
	SubscriptionId               string
	ResourceGroupName            string
	NetworkSecurityPerimeterName string
	ProfileName                  string
	AccessRuleName               string
}

// NewAccessRuleID returns a new AccessRuleId struct
func NewAccessRuleID(subscriptionId string, resourceGroupName string, networkSecurityPerimeterName string, profileName string, accessRuleName string) AccessRuleId {
	return AccessRuleId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		NetworkSecurityPerimeterName: networkSecurityPerimeterName,
		ProfileName:                  profileName,
		AccessRuleName:               accessRuleName,
	}
}

// ParseAccessRuleID parses 'input' into a AccessRuleId
func ParseAccessRuleID(input string) (*AccessRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccessRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccessRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAccessRuleIDInsensitively parses 'input' case-insensitively into a AccessRuleId
// note: this method should only be used for API response data and not user input
func ParseAccessRuleIDInsensitively(input string) (*AccessRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccessRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccessRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AccessRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkSecurityPerimeterName, ok = input.Parsed["networkSecurityPerimeterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkSecurityPerimeterName", input)
	}

	if id.ProfileName, ok = input.Parsed["profileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "profileName", input)
	}

	if id.AccessRuleName, ok = input.Parsed["accessRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accessRuleName", input)
	}

	return nil
}

// ValidateAccessRuleID checks that 'input' can be parsed as a Access Rule ID
func ValidateAccessRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccessRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Access Rule ID
func (id AccessRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityPerimeters/%s/profiles/%s/accessRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityPerimeterName, id.ProfileName, id.AccessRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Access Rule ID
func (id AccessRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkSecurityPerimeters", "networkSecurityPerimeters", "networkSecurityPerimeters"),
		resourceids.UserSpecifiedSegment("networkSecurityPerimeterName", "networkSecurityPerimeterName"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileName"),
		resourceids.StaticSegment("staticAccessRules", "accessRules", "accessRules"),
		resourceids.UserSpecifiedSegment("accessRuleName", "accessRuleName"),
	}
}

// String returns a human-readable description of this Access Rule ID
func (id AccessRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Security Perimeter Name: %q", id.NetworkSecurityPerimeterName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Access Rule Name: %q", id.AccessRuleName),
	}
	return fmt.Sprintf("Access Rule (%s)", strings.Join(components, "\n"))
}
