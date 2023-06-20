package referencedatasets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ReferenceDataSetId{}

// ReferenceDataSetId is a struct representing the Resource ID for a Reference Data Set
type ReferenceDataSetId struct {
	SubscriptionId       string
	ResourceGroupName    string
	EnvironmentName      string
	ReferenceDataSetName string
}

// NewReferenceDataSetID returns a new ReferenceDataSetId struct
func NewReferenceDataSetID(subscriptionId string, resourceGroupName string, environmentName string, referenceDataSetName string) ReferenceDataSetId {
	return ReferenceDataSetId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		EnvironmentName:      environmentName,
		ReferenceDataSetName: referenceDataSetName,
	}
}

// ParseReferenceDataSetID parses 'input' into a ReferenceDataSetId
func ParseReferenceDataSetID(input string) (*ReferenceDataSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReferenceDataSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReferenceDataSetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.EnvironmentName, ok = parsed.Parsed["environmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "environmentName", *parsed)
	}

	if id.ReferenceDataSetName, ok = parsed.Parsed["referenceDataSetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "referenceDataSetName", *parsed)
	}

	return &id, nil
}

// ParseReferenceDataSetIDInsensitively parses 'input' case-insensitively into a ReferenceDataSetId
// note: this method should only be used for API response data and not user input
func ParseReferenceDataSetIDInsensitively(input string) (*ReferenceDataSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReferenceDataSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReferenceDataSetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.EnvironmentName, ok = parsed.Parsed["environmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "environmentName", *parsed)
	}

	if id.ReferenceDataSetName, ok = parsed.Parsed["referenceDataSetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "referenceDataSetName", *parsed)
	}

	return &id, nil
}

// ValidateReferenceDataSetID checks that 'input' can be parsed as a Reference Data Set ID
func ValidateReferenceDataSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReferenceDataSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Reference Data Set ID
func (id ReferenceDataSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.TimeSeriesInsights/environments/%s/referenceDataSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.EnvironmentName, id.ReferenceDataSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Reference Data Set ID
func (id ReferenceDataSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftTimeSeriesInsights", "Microsoft.TimeSeriesInsights", "Microsoft.TimeSeriesInsights"),
		resourceids.StaticSegment("staticEnvironments", "environments", "environments"),
		resourceids.UserSpecifiedSegment("environmentName", "environmentValue"),
		resourceids.StaticSegment("staticReferenceDataSets", "referenceDataSets", "referenceDataSets"),
		resourceids.UserSpecifiedSegment("referenceDataSetName", "referenceDataSetValue"),
	}
}

// String returns a human-readable description of this Reference Data Set ID
func (id ReferenceDataSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Environment Name: %q", id.EnvironmentName),
		fmt.Sprintf("Reference Data Set Name: %q", id.ReferenceDataSetName),
	}
	return fmt.Sprintf("Reference Data Set (%s)", strings.Join(components, "\n"))
}
