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
	recaser.RegisterResourceId(&SitecontainerId{})
}

var _ resourceids.ResourceId = &SitecontainerId{}

// SitecontainerId is a struct representing the Resource ID for a Sitecontainer
type SitecontainerId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SitecontainerName string
}

// NewSitecontainerID returns a new SitecontainerId struct
func NewSitecontainerID(subscriptionId string, resourceGroupName string, siteName string, sitecontainerName string) SitecontainerId {
	return SitecontainerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SitecontainerName: sitecontainerName,
	}
}

// ParseSitecontainerID parses 'input' into a SitecontainerId
func ParseSitecontainerID(input string) (*SitecontainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SitecontainerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SitecontainerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSitecontainerIDInsensitively parses 'input' case-insensitively into a SitecontainerId
// note: this method should only be used for API response data and not user input
func ParseSitecontainerIDInsensitively(input string) (*SitecontainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SitecontainerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SitecontainerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SitecontainerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SitecontainerName, ok = input.Parsed["sitecontainerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sitecontainerName", input)
	}

	return nil
}

// ValidateSitecontainerID checks that 'input' can be parsed as a Sitecontainer ID
func ValidateSitecontainerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSitecontainerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sitecontainer ID
func (id SitecontainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/sitecontainers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SitecontainerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sitecontainer ID
func (id SitecontainerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticSitecontainers", "sitecontainers", "sitecontainers"),
		resourceids.UserSpecifiedSegment("sitecontainerName", "sitecontainerName"),
	}
}

// String returns a human-readable description of this Sitecontainer ID
func (id SitecontainerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Sitecontainer Name: %q", id.SitecontainerName),
	}
	return fmt.Sprintf("Sitecontainer (%s)", strings.Join(components, "\n"))
}
