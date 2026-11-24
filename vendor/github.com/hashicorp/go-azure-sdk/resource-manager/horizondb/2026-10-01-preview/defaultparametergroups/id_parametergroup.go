package defaultparametergroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ParameterGroupId{})
}

var _ resourceids.ResourceId = &ParameterGroupId{}

// ParameterGroupId is a struct representing the Resource ID for a Parameter Group
type ParameterGroupId struct {
	SubscriptionId     string
	LocationName       string
	ParameterGroupName string
}

// NewParameterGroupID returns a new ParameterGroupId struct
func NewParameterGroupID(subscriptionId string, locationName string, parameterGroupName string) ParameterGroupId {
	return ParameterGroupId{
		SubscriptionId:     subscriptionId,
		LocationName:       locationName,
		ParameterGroupName: parameterGroupName,
	}
}

// ParseParameterGroupID parses 'input' into a ParameterGroupId
func ParseParameterGroupID(input string) (*ParameterGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ParameterGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ParameterGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseParameterGroupIDInsensitively parses 'input' case-insensitively into a ParameterGroupId
// note: this method should only be used for API response data and not user input
func ParseParameterGroupIDInsensitively(input string) (*ParameterGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ParameterGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ParameterGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ParameterGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.ParameterGroupName, ok = input.Parsed["parameterGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "parameterGroupName", input)
	}

	return nil
}

// ValidateParameterGroupID checks that 'input' can be parsed as a Parameter Group ID
func ValidateParameterGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseParameterGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Parameter Group ID
func (id ParameterGroupId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.HorizonDb/locations/%s/parameterGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.ParameterGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Parameter Group ID
func (id ParameterGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHorizonDb", "Microsoft.HorizonDb", "Microsoft.HorizonDb"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticParameterGroups", "parameterGroups", "parameterGroups"),
		resourceids.UserSpecifiedSegment("parameterGroupName", "parameterGroupName"),
	}
}

// String returns a human-readable description of this Parameter Group ID
func (id ParameterGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Parameter Group Name: %q", id.ParameterGroupName),
	}
	return fmt.Sprintf("Parameter Group (%s)", strings.Join(components, "\n"))
}
