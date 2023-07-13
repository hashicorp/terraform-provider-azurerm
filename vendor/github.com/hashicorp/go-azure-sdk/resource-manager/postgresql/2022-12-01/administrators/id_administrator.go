package administrators

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AdministratorId{}

// AdministratorId is a struct representing the Resource ID for a Administrator
type AdministratorId struct {
	SubscriptionId     string
	ResourceGroupName  string
	FlexibleServerName string
	ObjectId           string
}

// NewAdministratorID returns a new AdministratorId struct
func NewAdministratorID(subscriptionId string, resourceGroupName string, flexibleServerName string, objectId string) AdministratorId {
	return AdministratorId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		FlexibleServerName: flexibleServerName,
		ObjectId:           objectId,
	}
}

// ParseAdministratorID parses 'input' into a AdministratorId
func ParseAdministratorID(input string) (*AdministratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(AdministratorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AdministratorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FlexibleServerName, ok = parsed.Parsed["flexibleServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "flexibleServerName", *parsed)
	}

	if id.ObjectId, ok = parsed.Parsed["objectId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "objectId", *parsed)
	}

	return &id, nil
}

// ParseAdministratorIDInsensitively parses 'input' case-insensitively into a AdministratorId
// note: this method should only be used for API response data and not user input
func ParseAdministratorIDInsensitively(input string) (*AdministratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(AdministratorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AdministratorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FlexibleServerName, ok = parsed.Parsed["flexibleServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "flexibleServerName", *parsed)
	}

	if id.ObjectId, ok = parsed.Parsed["objectId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "objectId", *parsed)
	}

	return &id, nil
}

// ValidateAdministratorID checks that 'input' can be parsed as a Administrator ID
func ValidateAdministratorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAdministratorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Administrator ID
func (id AdministratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/flexibleServers/%s/administrators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName, id.ObjectId)
}

// Segments returns a slice of Resource ID Segments which comprise this Administrator ID
func (id AdministratorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticFlexibleServers", "flexibleServers", "flexibleServers"),
		resourceids.UserSpecifiedSegment("flexibleServerName", "flexibleServerValue"),
		resourceids.StaticSegment("staticAdministrators", "administrators", "administrators"),
		resourceids.UserSpecifiedSegment("objectId", "objectIdValue"),
	}
}

// String returns a human-readable description of this Administrator ID
func (id AdministratorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Flexible Server Name: %q", id.FlexibleServerName),
		fmt.Sprintf("Object: %q", id.ObjectId),
	}
	return fmt.Sprintf("Administrator (%s)", strings.Join(components, "\n"))
}
