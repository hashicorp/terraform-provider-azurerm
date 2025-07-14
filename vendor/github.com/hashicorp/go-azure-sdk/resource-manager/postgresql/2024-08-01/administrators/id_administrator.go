package administrators

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AdministratorId{})
}

var _ resourceids.ResourceId = &AdministratorId{}

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
	parser := resourceids.NewParserFromResourceIdType(&AdministratorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AdministratorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAdministratorIDInsensitively parses 'input' case-insensitively into a AdministratorId
// note: this method should only be used for API response data and not user input
func ParseAdministratorIDInsensitively(input string) (*AdministratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AdministratorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AdministratorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AdministratorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FlexibleServerName, ok = input.Parsed["flexibleServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "flexibleServerName", input)
	}

	if id.ObjectId, ok = input.Parsed["objectId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "objectId", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("flexibleServerName", "flexibleServerName"),
		resourceids.StaticSegment("staticAdministrators", "administrators", "administrators"),
		resourceids.UserSpecifiedSegment("objectId", "objectId"),
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
