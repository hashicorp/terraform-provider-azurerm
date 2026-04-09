package remediations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProviderRemediationId{})
}

var _ resourceids.ResourceId = &ProviderRemediationId{}

// ProviderRemediationId is a struct representing the Resource ID for a Provider Remediation
type ProviderRemediationId struct {
	SubscriptionId    string
	ResourceGroupName string
	RemediationName   string
}

// NewProviderRemediationID returns a new ProviderRemediationId struct
func NewProviderRemediationID(subscriptionId string, resourceGroupName string, remediationName string) ProviderRemediationId {
	return ProviderRemediationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RemediationName:   remediationName,
	}
}

// ParseProviderRemediationID parses 'input' into a ProviderRemediationId
func ParseProviderRemediationID(input string) (*ProviderRemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderRemediationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderRemediationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviderRemediationIDInsensitively parses 'input' case-insensitively into a ProviderRemediationId
// note: this method should only be used for API response data and not user input
func ParseProviderRemediationIDInsensitively(input string) (*ProviderRemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderRemediationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderRemediationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProviderRemediationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RemediationName, ok = input.Parsed["remediationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "remediationName", input)
	}

	return nil
}

// ValidateProviderRemediationID checks that 'input' can be parsed as a Provider Remediation ID
func ValidateProviderRemediationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderRemediationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Remediation ID
func (id ProviderRemediationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.PolicyInsights/remediations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RemediationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Remediation ID
func (id ProviderRemediationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPolicyInsights", "Microsoft.PolicyInsights", "Microsoft.PolicyInsights"),
		resourceids.StaticSegment("staticRemediations", "remediations", "remediations"),
		resourceids.UserSpecifiedSegment("remediationName", "remediationName"),
	}
}

// String returns a human-readable description of this Provider Remediation ID
func (id ProviderRemediationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Remediation Name: %q", id.RemediationName),
	}
	return fmt.Sprintf("Provider Remediation (%s)", strings.Join(components, "\n"))
}
