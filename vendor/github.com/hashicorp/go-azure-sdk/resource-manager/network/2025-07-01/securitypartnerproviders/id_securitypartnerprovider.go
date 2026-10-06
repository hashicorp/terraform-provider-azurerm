package securitypartnerproviders

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SecurityPartnerProviderId{})
}

var _ resourceids.ResourceId = &SecurityPartnerProviderId{}

// SecurityPartnerProviderId is a struct representing the Resource ID for a Security Partner Provider
type SecurityPartnerProviderId struct {
	SubscriptionId              string
	ResourceGroupName           string
	SecurityPartnerProviderName string
}

// NewSecurityPartnerProviderID returns a new SecurityPartnerProviderId struct
func NewSecurityPartnerProviderID(subscriptionId string, resourceGroupName string, securityPartnerProviderName string) SecurityPartnerProviderId {
	return SecurityPartnerProviderId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		SecurityPartnerProviderName: securityPartnerProviderName,
	}
}

// ParseSecurityPartnerProviderID parses 'input' into a SecurityPartnerProviderId
func ParseSecurityPartnerProviderID(input string) (*SecurityPartnerProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecurityPartnerProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecurityPartnerProviderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSecurityPartnerProviderIDInsensitively parses 'input' case-insensitively into a SecurityPartnerProviderId
// note: this method should only be used for API response data and not user input
func ParseSecurityPartnerProviderIDInsensitively(input string) (*SecurityPartnerProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecurityPartnerProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecurityPartnerProviderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SecurityPartnerProviderId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SecurityPartnerProviderName, ok = input.Parsed["securityPartnerProviderName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "securityPartnerProviderName", input)
	}

	return nil
}

// ValidateSecurityPartnerProviderID checks that 'input' can be parsed as a Security Partner Provider ID
func ValidateSecurityPartnerProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecurityPartnerProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Security Partner Provider ID
func (id SecurityPartnerProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/securityPartnerProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SecurityPartnerProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Security Partner Provider ID
func (id SecurityPartnerProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticSecurityPartnerProviders", "securityPartnerProviders", "securityPartnerProviders"),
		resourceids.UserSpecifiedSegment("securityPartnerProviderName", "securityPartnerProviderName"),
	}
}

// String returns a human-readable description of this Security Partner Provider ID
func (id SecurityPartnerProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Security Partner Provider Name: %q", id.SecurityPartnerProviderName),
	}
	return fmt.Sprintf("Security Partner Provider (%s)", strings.Join(components, "\n"))
}
