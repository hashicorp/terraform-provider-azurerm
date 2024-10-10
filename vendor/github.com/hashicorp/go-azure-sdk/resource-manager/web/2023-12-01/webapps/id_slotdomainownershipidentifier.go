package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SlotDomainOwnershipIdentifierId{})
}

var _ resourceids.ResourceId = &SlotDomainOwnershipIdentifierId{}

// SlotDomainOwnershipIdentifierId is a struct representing the Resource ID for a Slot Domain Ownership Identifier
type SlotDomainOwnershipIdentifierId struct {
	SubscriptionId                string
	ResourceGroupName             string
	SiteName                      string
	SlotName                      string
	DomainOwnershipIdentifierName string
}

// NewSlotDomainOwnershipIdentifierID returns a new SlotDomainOwnershipIdentifierId struct
func NewSlotDomainOwnershipIdentifierID(subscriptionId string, resourceGroupName string, siteName string, slotName string, domainOwnershipIdentifierName string) SlotDomainOwnershipIdentifierId {
	return SlotDomainOwnershipIdentifierId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		SiteName:                      siteName,
		SlotName:                      slotName,
		DomainOwnershipIdentifierName: domainOwnershipIdentifierName,
	}
}

// ParseSlotDomainOwnershipIdentifierID parses 'input' into a SlotDomainOwnershipIdentifierId
func ParseSlotDomainOwnershipIdentifierID(input string) (*SlotDomainOwnershipIdentifierId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotDomainOwnershipIdentifierId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotDomainOwnershipIdentifierId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotDomainOwnershipIdentifierIDInsensitively parses 'input' case-insensitively into a SlotDomainOwnershipIdentifierId
// note: this method should only be used for API response data and not user input
func ParseSlotDomainOwnershipIdentifierIDInsensitively(input string) (*SlotDomainOwnershipIdentifierId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotDomainOwnershipIdentifierId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotDomainOwnershipIdentifierId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotDomainOwnershipIdentifierId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.SlotName, ok = input.Parsed["slotName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "slotName", input)
	}

	if id.DomainOwnershipIdentifierName, ok = input.Parsed["domainOwnershipIdentifierName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "domainOwnershipIdentifierName", input)
	}

	return nil
}

// ValidateSlotDomainOwnershipIdentifierID checks that 'input' can be parsed as a Slot Domain Ownership Identifier ID
func ValidateSlotDomainOwnershipIdentifierID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotDomainOwnershipIdentifierID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Domain Ownership Identifier ID
func (id SlotDomainOwnershipIdentifierId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/domainOwnershipIdentifiers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.DomainOwnershipIdentifierName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Domain Ownership Identifier ID
func (id SlotDomainOwnershipIdentifierId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticSlots", "slots", "slots"),
		resourceids.UserSpecifiedSegment("slotName", "slotName"),
		resourceids.StaticSegment("staticDomainOwnershipIdentifiers", "domainOwnershipIdentifiers", "domainOwnershipIdentifiers"),
		resourceids.UserSpecifiedSegment("domainOwnershipIdentifierName", "domainOwnershipIdentifierName"),
	}
}

// String returns a human-readable description of this Slot Domain Ownership Identifier ID
func (id SlotDomainOwnershipIdentifierId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Domain Ownership Identifier Name: %q", id.DomainOwnershipIdentifierName),
	}
	return fmt.Sprintf("Slot Domain Ownership Identifier (%s)", strings.Join(components, "\n"))
}
