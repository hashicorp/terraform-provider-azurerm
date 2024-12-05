package securitypolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SecurityPolicyId{})
}

var _ resourceids.ResourceId = &SecurityPolicyId{}

// SecurityPolicyId is a struct representing the Resource ID for a Security Policy
type SecurityPolicyId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ProfileName        string
	SecurityPolicyName string
}

// NewSecurityPolicyID returns a new SecurityPolicyId struct
func NewSecurityPolicyID(subscriptionId string, resourceGroupName string, profileName string, securityPolicyName string) SecurityPolicyId {
	return SecurityPolicyId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ProfileName:        profileName,
		SecurityPolicyName: securityPolicyName,
	}
}

// ParseSecurityPolicyID parses 'input' into a SecurityPolicyId
func ParseSecurityPolicyID(input string) (*SecurityPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecurityPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecurityPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSecurityPolicyIDInsensitively parses 'input' case-insensitively into a SecurityPolicyId
// note: this method should only be used for API response data and not user input
func ParseSecurityPolicyIDInsensitively(input string) (*SecurityPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecurityPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecurityPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SecurityPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProfileName, ok = input.Parsed["profileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "profileName", input)
	}

	if id.SecurityPolicyName, ok = input.Parsed["securityPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "securityPolicyName", input)
	}

	return nil
}

// ValidateSecurityPolicyID checks that 'input' can be parsed as a Security Policy ID
func ValidateSecurityPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecurityPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Security Policy ID
func (id SecurityPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/securityPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.SecurityPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Security Policy ID
func (id SecurityPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCDN", "Microsoft.CDN", "Microsoft.CDN"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileName"),
		resourceids.StaticSegment("staticSecurityPolicies", "securityPolicies", "securityPolicies"),
		resourceids.UserSpecifiedSegment("securityPolicyName", "securityPolicyName"),
	}
}

// String returns a human-readable description of this Security Policy ID
func (id SecurityPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Security Policy Name: %q", id.SecurityPolicyName),
	}
	return fmt.Sprintf("Security Policy (%s)", strings.Join(components, "\n"))
}
