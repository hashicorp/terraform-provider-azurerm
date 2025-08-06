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
	recaser.RegisterResourceId(&SiteDomainOwnershipIdentifierId{})
}

var _ resourceids.ResourceId = &SiteDomainOwnershipIdentifierId{}

// SiteDomainOwnershipIdentifierId is a struct representing the Resource ID for a Site Domain Ownership Identifier
type SiteDomainOwnershipIdentifierId struct {
	SubscriptionId                string
	ResourceGroupName             string
	SiteName                      string
	DomainOwnershipIdentifierName string
}

// NewSiteDomainOwnershipIdentifierID returns a new SiteDomainOwnershipIdentifierId struct
func NewSiteDomainOwnershipIdentifierID(subscriptionId string, resourceGroupName string, siteName string, domainOwnershipIdentifierName string) SiteDomainOwnershipIdentifierId {
	return SiteDomainOwnershipIdentifierId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		SiteName:                      siteName,
		DomainOwnershipIdentifierName: domainOwnershipIdentifierName,
	}
}

// ParseSiteDomainOwnershipIdentifierID parses 'input' into a SiteDomainOwnershipIdentifierId
func ParseSiteDomainOwnershipIdentifierID(input string) (*SiteDomainOwnershipIdentifierId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SiteDomainOwnershipIdentifierId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SiteDomainOwnershipIdentifierId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSiteDomainOwnershipIdentifierIDInsensitively parses 'input' case-insensitively into a SiteDomainOwnershipIdentifierId
// note: this method should only be used for API response data and not user input
func ParseSiteDomainOwnershipIdentifierIDInsensitively(input string) (*SiteDomainOwnershipIdentifierId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SiteDomainOwnershipIdentifierId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SiteDomainOwnershipIdentifierId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SiteDomainOwnershipIdentifierId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DomainOwnershipIdentifierName, ok = input.Parsed["domainOwnershipIdentifierName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "domainOwnershipIdentifierName", input)
	}

	return nil
}

// ValidateSiteDomainOwnershipIdentifierID checks that 'input' can be parsed as a Site Domain Ownership Identifier ID
func ValidateSiteDomainOwnershipIdentifierID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSiteDomainOwnershipIdentifierID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Site Domain Ownership Identifier ID
func (id SiteDomainOwnershipIdentifierId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/domainOwnershipIdentifiers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.DomainOwnershipIdentifierName)
}

// Segments returns a slice of Resource ID Segments which comprise this Site Domain Ownership Identifier ID
func (id SiteDomainOwnershipIdentifierId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticDomainOwnershipIdentifiers", "domainOwnershipIdentifiers", "domainOwnershipIdentifiers"),
		resourceids.UserSpecifiedSegment("domainOwnershipIdentifierName", "domainOwnershipIdentifierName"),
	}
}

// String returns a human-readable description of this Site Domain Ownership Identifier ID
func (id SiteDomainOwnershipIdentifierId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Domain Ownership Identifier Name: %q", id.DomainOwnershipIdentifierName),
	}
	return fmt.Sprintf("Site Domain Ownership Identifier (%s)", strings.Join(components, "\n"))
}
