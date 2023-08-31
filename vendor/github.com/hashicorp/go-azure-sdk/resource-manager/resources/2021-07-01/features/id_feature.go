package features

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FeatureId{}

// FeatureId is a struct representing the Resource ID for a Feature
type FeatureId struct {
	SubscriptionId string
	ProviderName   string
	FeatureName    string
}

// NewFeatureID returns a new FeatureId struct
func NewFeatureID(subscriptionId string, providerName string, featureName string) FeatureId {
	return FeatureId{
		SubscriptionId: subscriptionId,
		ProviderName:   providerName,
		FeatureName:    featureName,
	}
}

// ParseFeatureID parses 'input' into a FeatureId
func ParseFeatureID(input string) (*FeatureId, error) {
	parser := resourceids.NewParserFromResourceIdType(FeatureId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FeatureId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	if id.FeatureName, ok = parsed.Parsed["featureName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "featureName", *parsed)
	}

	return &id, nil
}

// ParseFeatureIDInsensitively parses 'input' case-insensitively into a FeatureId
// note: this method should only be used for API response data and not user input
func ParseFeatureIDInsensitively(input string) (*FeatureId, error) {
	parser := resourceids.NewParserFromResourceIdType(FeatureId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FeatureId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	if id.FeatureName, ok = parsed.Parsed["featureName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "featureName", *parsed)
	}

	return &id, nil
}

// ValidateFeatureID checks that 'input' can be parsed as a Feature ID
func ValidateFeatureID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFeatureID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Feature ID
func (id FeatureId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Features/providers/%s/features/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ProviderName, id.FeatureName)
}

// Segments returns a slice of Resource ID Segments which comprise this Feature ID
func (id FeatureId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftFeatures", "Microsoft.Features", "Microsoft.Features"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.UserSpecifiedSegment("providerName", "providerValue"),
		resourceids.StaticSegment("staticFeatures", "features", "features"),
		resourceids.UserSpecifiedSegment("featureName", "featureValue"),
	}
}

// String returns a human-readable description of this Feature ID
func (id FeatureId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
		fmt.Sprintf("Feature Name: %q", id.FeatureName),
	}
	return fmt.Sprintf("Feature (%s)", strings.Join(components, "\n"))
}
