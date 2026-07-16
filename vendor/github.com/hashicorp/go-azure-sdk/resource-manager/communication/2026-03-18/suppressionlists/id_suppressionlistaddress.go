package suppressionlists

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SuppressionListAddressId{})
}

var _ resourceids.ResourceId = &SuppressionListAddressId{}

// SuppressionListAddressId is a struct representing the Resource ID for a Suppression List Address
type SuppressionListAddressId struct {
	SubscriptionId      string
	ResourceGroupName   string
	EmailServiceName    string
	DomainName          string
	SuppressionListName string
	AddressId           string
}

// NewSuppressionListAddressID returns a new SuppressionListAddressId struct
func NewSuppressionListAddressID(subscriptionId string, resourceGroupName string, emailServiceName string, domainName string, suppressionListName string, addressId string) SuppressionListAddressId {
	return SuppressionListAddressId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		EmailServiceName:    emailServiceName,
		DomainName:          domainName,
		SuppressionListName: suppressionListName,
		AddressId:           addressId,
	}
}

// ParseSuppressionListAddressID parses 'input' into a SuppressionListAddressId
func ParseSuppressionListAddressID(input string) (*SuppressionListAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SuppressionListAddressId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SuppressionListAddressId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSuppressionListAddressIDInsensitively parses 'input' case-insensitively into a SuppressionListAddressId
// note: this method should only be used for API response data and not user input
func ParseSuppressionListAddressIDInsensitively(input string) (*SuppressionListAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SuppressionListAddressId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SuppressionListAddressId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SuppressionListAddressId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.EmailServiceName, ok = input.Parsed["emailServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "emailServiceName", input)
	}

	if id.DomainName, ok = input.Parsed["domainName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "domainName", input)
	}

	if id.SuppressionListName, ok = input.Parsed["suppressionListName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "suppressionListName", input)
	}

	if id.AddressId, ok = input.Parsed["addressId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "addressId", input)
	}

	return nil
}

// ValidateSuppressionListAddressID checks that 'input' can be parsed as a Suppression List Address ID
func ValidateSuppressionListAddressID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSuppressionListAddressID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Suppression List Address ID
func (id SuppressionListAddressId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Communication/emailServices/%s/domains/%s/suppressionLists/%s/suppressionListAddresses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.EmailServiceName, id.DomainName, id.SuppressionListName, id.AddressId)
}

// Segments returns a slice of Resource ID Segments which comprise this Suppression List Address ID
func (id SuppressionListAddressId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCommunication", "Microsoft.Communication", "Microsoft.Communication"),
		resourceids.StaticSegment("staticEmailServices", "emailServices", "emailServices"),
		resourceids.UserSpecifiedSegment("emailServiceName", "emailServiceName"),
		resourceids.StaticSegment("staticDomains", "domains", "domains"),
		resourceids.UserSpecifiedSegment("domainName", "domainName"),
		resourceids.StaticSegment("staticSuppressionLists", "suppressionLists", "suppressionLists"),
		resourceids.UserSpecifiedSegment("suppressionListName", "suppressionListName"),
		resourceids.StaticSegment("staticSuppressionListAddresses", "suppressionListAddresses", "suppressionListAddresses"),
		resourceids.UserSpecifiedSegment("addressId", "addressId"),
	}
}

// String returns a human-readable description of this Suppression List Address ID
func (id SuppressionListAddressId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Email Service Name: %q", id.EmailServiceName),
		fmt.Sprintf("Domain Name: %q", id.DomainName),
		fmt.Sprintf("Suppression List Name: %q", id.SuppressionListName),
		fmt.Sprintf("Address: %q", id.AddressId),
	}
	return fmt.Sprintf("Suppression List Address (%s)", strings.Join(components, "\n"))
}
