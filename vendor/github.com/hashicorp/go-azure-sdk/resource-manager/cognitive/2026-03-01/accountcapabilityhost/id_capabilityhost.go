package accountcapabilityhost

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CapabilityHostId{})
}

var _ resourceids.ResourceId = &CapabilityHostId{}

// CapabilityHostId is a struct representing the Resource ID for a Capability Host
type CapabilityHostId struct {
	SubscriptionId     string
	ResourceGroupName  string
	AccountName        string
	CapabilityHostName string
}

// NewCapabilityHostID returns a new CapabilityHostId struct
func NewCapabilityHostID(subscriptionId string, resourceGroupName string, accountName string, capabilityHostName string) CapabilityHostId {
	return CapabilityHostId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		AccountName:        accountName,
		CapabilityHostName: capabilityHostName,
	}
}

// ParseCapabilityHostID parses 'input' into a CapabilityHostId
func ParseCapabilityHostID(input string) (*CapabilityHostId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CapabilityHostId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapabilityHostId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCapabilityHostIDInsensitively parses 'input' case-insensitively into a CapabilityHostId
// note: this method should only be used for API response data and not user input
func ParseCapabilityHostIDInsensitively(input string) (*CapabilityHostId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CapabilityHostId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapabilityHostId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CapabilityHostId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.CapabilityHostName, ok = input.Parsed["capabilityHostName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "capabilityHostName", input)
	}

	return nil
}

// ValidateCapabilityHostID checks that 'input' can be parsed as a Capability Host ID
func ValidateCapabilityHostID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCapabilityHostID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Capability Host ID
func (id CapabilityHostId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s/capabilityHosts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.CapabilityHostName)
}

// Segments returns a slice of Resource ID Segments which comprise this Capability Host ID
func (id CapabilityHostId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCognitiveServices", "Microsoft.CognitiveServices", "Microsoft.CognitiveServices"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountName"),
		resourceids.StaticSegment("staticCapabilityHosts", "capabilityHosts", "capabilityHosts"),
		resourceids.UserSpecifiedSegment("capabilityHostName", "capabilityHostName"),
	}
}

// String returns a human-readable description of this Capability Host ID
func (id CapabilityHostId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Capability Host Name: %q", id.CapabilityHostName),
	}
	return fmt.Sprintf("Capability Host (%s)", strings.Join(components, "\n"))
}
