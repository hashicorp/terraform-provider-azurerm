package raiblocklists

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RaiBlocklistItemId{})
}

var _ resourceids.ResourceId = &RaiBlocklistItemId{}

// RaiBlocklistItemId is a struct representing the Resource ID for a Rai Blocklist Item
type RaiBlocklistItemId struct {
	SubscriptionId       string
	ResourceGroupName    string
	AccountName          string
	RaiBlocklistName     string
	RaiBlocklistItemName string
}

// NewRaiBlocklistItemID returns a new RaiBlocklistItemId struct
func NewRaiBlocklistItemID(subscriptionId string, resourceGroupName string, accountName string, raiBlocklistName string, raiBlocklistItemName string) RaiBlocklistItemId {
	return RaiBlocklistItemId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		AccountName:          accountName,
		RaiBlocklistName:     raiBlocklistName,
		RaiBlocklistItemName: raiBlocklistItemName,
	}
}

// ParseRaiBlocklistItemID parses 'input' into a RaiBlocklistItemId
func ParseRaiBlocklistItemID(input string) (*RaiBlocklistItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RaiBlocklistItemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RaiBlocklistItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRaiBlocklistItemIDInsensitively parses 'input' case-insensitively into a RaiBlocklistItemId
// note: this method should only be used for API response data and not user input
func ParseRaiBlocklistItemIDInsensitively(input string) (*RaiBlocklistItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RaiBlocklistItemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RaiBlocklistItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RaiBlocklistItemId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AccountName, ok = input.Parsed["accountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accountName", input)
	}

	if id.RaiBlocklistName, ok = input.Parsed["raiBlocklistName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "raiBlocklistName", input)
	}

	if id.RaiBlocklistItemName, ok = input.Parsed["raiBlocklistItemName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "raiBlocklistItemName", input)
	}

	return nil
}

// ValidateRaiBlocklistItemID checks that 'input' can be parsed as a Rai Blocklist Item ID
func ValidateRaiBlocklistItemID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRaiBlocklistItemID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rai Blocklist Item ID
func (id RaiBlocklistItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s/raiBlocklists/%s/raiBlocklistItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.RaiBlocklistName, id.RaiBlocklistItemName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rai Blocklist Item ID
func (id RaiBlocklistItemId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCognitiveServices", "Microsoft.CognitiveServices", "Microsoft.CognitiveServices"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountName"),
		resourceids.StaticSegment("staticRaiBlocklists", "raiBlocklists", "raiBlocklists"),
		resourceids.UserSpecifiedSegment("raiBlocklistName", "raiBlocklistName"),
		resourceids.StaticSegment("staticRaiBlocklistItems", "raiBlocklistItems", "raiBlocklistItems"),
		resourceids.UserSpecifiedSegment("raiBlocklistItemName", "raiBlocklistItemName"),
	}
}

// String returns a human-readable description of this Rai Blocklist Item ID
func (id RaiBlocklistItemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Rai Blocklist Name: %q", id.RaiBlocklistName),
		fmt.Sprintf("Rai Blocklist Item Name: %q", id.RaiBlocklistItemName),
	}
	return fmt.Sprintf("Rai Blocklist Item (%s)", strings.Join(components, "\n"))
}
