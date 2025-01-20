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
	recaser.RegisterResourceId(&RemediationId{})
}

var _ resourceids.ResourceId = &RemediationId{}

// RemediationId is a struct representing the Resource ID for a Remediation
type RemediationId struct {
	SubscriptionId  string
	RemediationName string
}

// NewRemediationID returns a new RemediationId struct
func NewRemediationID(subscriptionId string, remediationName string) RemediationId {
	return RemediationId{
		SubscriptionId:  subscriptionId,
		RemediationName: remediationName,
	}
}

// ParseRemediationID parses 'input' into a RemediationId
func ParseRemediationID(input string) (*RemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RemediationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RemediationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRemediationIDInsensitively parses 'input' case-insensitively into a RemediationId
// note: this method should only be used for API response data and not user input
func ParseRemediationIDInsensitively(input string) (*RemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RemediationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RemediationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RemediationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.RemediationName, ok = input.Parsed["remediationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "remediationName", input)
	}

	return nil
}

// ValidateRemediationID checks that 'input' can be parsed as a Remediation ID
func ValidateRemediationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRemediationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Remediation ID
func (id RemediationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.PolicyInsights/remediations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.RemediationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Remediation ID
func (id RemediationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPolicyInsights", "Microsoft.PolicyInsights", "Microsoft.PolicyInsights"),
		resourceids.StaticSegment("staticRemediations", "remediations", "remediations"),
		resourceids.UserSpecifiedSegment("remediationName", "remediationName"),
	}
}

// String returns a human-readable description of this Remediation ID
func (id RemediationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Remediation Name: %q", id.RemediationName),
	}
	return fmt.Sprintf("Remediation (%s)", strings.Join(components, "\n"))
}
