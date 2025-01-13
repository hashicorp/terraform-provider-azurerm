package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApiPortalDomainId{})
}

var _ resourceids.ResourceId = &ApiPortalDomainId{}

// ApiPortalDomainId is a struct representing the Resource ID for a Api Portal Domain
type ApiPortalDomainId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	ApiPortalName     string
	DomainName        string
}

// NewApiPortalDomainID returns a new ApiPortalDomainId struct
func NewApiPortalDomainID(subscriptionId string, resourceGroupName string, springName string, apiPortalName string, domainName string) ApiPortalDomainId {
	return ApiPortalDomainId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		ApiPortalName:     apiPortalName,
		DomainName:        domainName,
	}
}

// ParseApiPortalDomainID parses 'input' into a ApiPortalDomainId
func ParseApiPortalDomainID(input string) (*ApiPortalDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiPortalDomainId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiPortalDomainId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApiPortalDomainIDInsensitively parses 'input' case-insensitively into a ApiPortalDomainId
// note: this method should only be used for API response data and not user input
func ParseApiPortalDomainIDInsensitively(input string) (*ApiPortalDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiPortalDomainId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiPortalDomainId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApiPortalDomainId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.ApiPortalName, ok = input.Parsed["apiPortalName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiPortalName", input)
	}

	if id.DomainName, ok = input.Parsed["domainName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "domainName", input)
	}

	return nil
}

// ValidateApiPortalDomainID checks that 'input' can be parsed as a Api Portal Domain ID
func ValidateApiPortalDomainID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApiPortalDomainID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Api Portal Domain ID
func (id ApiPortalDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apiPortals/%s/domains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApiPortalName, id.DomainName)
}

// Segments returns a slice of Resource ID Segments which comprise this Api Portal Domain ID
func (id ApiPortalDomainId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticApiPortals", "apiPortals", "apiPortals"),
		resourceids.UserSpecifiedSegment("apiPortalName", "apiPortalName"),
		resourceids.StaticSegment("staticDomains", "domains", "domains"),
		resourceids.UserSpecifiedSegment("domainName", "domainName"),
	}
}

// String returns a human-readable description of this Api Portal Domain ID
func (id ApiPortalDomainId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Api Portal Name: %q", id.ApiPortalName),
		fmt.Sprintf("Domain Name: %q", id.DomainName),
	}
	return fmt.Sprintf("Api Portal Domain (%s)", strings.Join(components, "\n"))
}
