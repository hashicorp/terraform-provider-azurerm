package restorepoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RestorePointId{})
}

var _ resourceids.ResourceId = &RestorePointId{}

// RestorePointId is a struct representing the Resource ID for a Restore Point
type RestorePointId struct {
	SubscriptionId             string
	ResourceGroupName          string
	RestorePointCollectionName string
	RestorePointName           string
}

// NewRestorePointID returns a new RestorePointId struct
func NewRestorePointID(subscriptionId string, resourceGroupName string, restorePointCollectionName string, restorePointName string) RestorePointId {
	return RestorePointId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		RestorePointCollectionName: restorePointCollectionName,
		RestorePointName:           restorePointName,
	}
}

// ParseRestorePointID parses 'input' into a RestorePointId
func ParseRestorePointID(input string) (*RestorePointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RestorePointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RestorePointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRestorePointIDInsensitively parses 'input' case-insensitively into a RestorePointId
// note: this method should only be used for API response data and not user input
func ParseRestorePointIDInsensitively(input string) (*RestorePointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RestorePointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RestorePointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RestorePointId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RestorePointCollectionName, ok = input.Parsed["restorePointCollectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "restorePointCollectionName", input)
	}

	if id.RestorePointName, ok = input.Parsed["restorePointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "restorePointName", input)
	}

	return nil
}

// ValidateRestorePointID checks that 'input' can be parsed as a Restore Point ID
func ValidateRestorePointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRestorePointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Restore Point ID
func (id RestorePointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/restorePointCollections/%s/restorePoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RestorePointCollectionName, id.RestorePointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Restore Point ID
func (id RestorePointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticRestorePointCollections", "restorePointCollections", "restorePointCollections"),
		resourceids.UserSpecifiedSegment("restorePointCollectionName", "restorePointCollectionName"),
		resourceids.StaticSegment("staticRestorePoints", "restorePoints", "restorePoints"),
		resourceids.UserSpecifiedSegment("restorePointName", "restorePointName"),
	}
}

// String returns a human-readable description of this Restore Point ID
func (id RestorePointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Restore Point Collection Name: %q", id.RestorePointCollectionName),
		fmt.Sprintf("Restore Point Name: %q", id.RestorePointName),
	}
	return fmt.Sprintf("Restore Point (%s)", strings.Join(components, "\n"))
}
