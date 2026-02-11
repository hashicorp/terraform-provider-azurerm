package applicationgateways

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PredefinedPolicyId{})
}

var _ resourceids.ResourceId = &PredefinedPolicyId{}

// PredefinedPolicyId is a struct representing the Resource ID for a Predefined Policy
type PredefinedPolicyId struct {
	SubscriptionId       string
	PredefinedPolicyName string
}

// NewPredefinedPolicyID returns a new PredefinedPolicyId struct
func NewPredefinedPolicyID(subscriptionId string, predefinedPolicyName string) PredefinedPolicyId {
	return PredefinedPolicyId{
		SubscriptionId:       subscriptionId,
		PredefinedPolicyName: predefinedPolicyName,
	}
}

// ParsePredefinedPolicyID parses 'input' into a PredefinedPolicyId
func ParsePredefinedPolicyID(input string) (*PredefinedPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PredefinedPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PredefinedPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePredefinedPolicyIDInsensitively parses 'input' case-insensitively into a PredefinedPolicyId
// note: this method should only be used for API response data and not user input
func ParsePredefinedPolicyIDInsensitively(input string) (*PredefinedPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PredefinedPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PredefinedPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PredefinedPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.PredefinedPolicyName, ok = input.Parsed["predefinedPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "predefinedPolicyName", input)
	}

	return nil
}

// ValidatePredefinedPolicyID checks that 'input' can be parsed as a Predefined Policy ID
func ValidatePredefinedPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePredefinedPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Predefined Policy ID
func (id PredefinedPolicyId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default/predefinedPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PredefinedPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Predefined Policy ID
func (id PredefinedPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticApplicationGatewayAvailableSslOptions", "applicationGatewayAvailableSslOptions", "applicationGatewayAvailableSslOptions"),
		resourceids.StaticSegment("staticDefault", "default", "default"),
		resourceids.StaticSegment("staticPredefinedPolicies", "predefinedPolicies", "predefinedPolicies"),
		resourceids.UserSpecifiedSegment("predefinedPolicyName", "predefinedPolicyName"),
	}
}

// String returns a human-readable description of this Predefined Policy ID
func (id PredefinedPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Predefined Policy Name: %q", id.PredefinedPolicyName),
	}
	return fmt.Sprintf("Predefined Policy (%s)", strings.Join(components, "\n"))
}
