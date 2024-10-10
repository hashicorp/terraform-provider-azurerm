package cognitiveservicesaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DeletedAccountId{})
}

var _ resourceids.ResourceId = &DeletedAccountId{}

// DeletedAccountId is a struct representing the Resource ID for a Deleted Account
type DeletedAccountId struct {
	SubscriptionId     string
	LocationName       string
	ResourceGroupName  string
	DeletedAccountName string
}

// NewDeletedAccountID returns a new DeletedAccountId struct
func NewDeletedAccountID(subscriptionId string, locationName string, resourceGroupName string, deletedAccountName string) DeletedAccountId {
	return DeletedAccountId{
		SubscriptionId:     subscriptionId,
		LocationName:       locationName,
		ResourceGroupName:  resourceGroupName,
		DeletedAccountName: deletedAccountName,
	}
}

// ParseDeletedAccountID parses 'input' into a DeletedAccountId
func ParseDeletedAccountID(input string) (*DeletedAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeletedAccountIDInsensitively parses 'input' case-insensitively into a DeletedAccountId
// note: this method should only be used for API response data and not user input
func ParseDeletedAccountIDInsensitively(input string) (*DeletedAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeletedAccountId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DeletedAccountName, ok = input.Parsed["deletedAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deletedAccountName", input)
	}

	return nil
}

// ValidateDeletedAccountID checks that 'input' can be parsed as a Deleted Account ID
func ValidateDeletedAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deleted Account ID
func (id DeletedAccountId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.CognitiveServices/locations/%s/resourceGroups/%s/deletedAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.ResourceGroupName, id.DeletedAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deleted Account ID
func (id DeletedAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCognitiveServices", "Microsoft.CognitiveServices", "Microsoft.CognitiveServices"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticDeletedAccounts", "deletedAccounts", "deletedAccounts"),
		resourceids.UserSpecifiedSegment("deletedAccountName", "deletedAccountName"),
	}
}

// String returns a human-readable description of this Deleted Account ID
func (id DeletedAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Deleted Account Name: %q", id.DeletedAccountName),
	}
	return fmt.Sprintf("Deleted Account (%s)", strings.Join(components, "\n"))
}
