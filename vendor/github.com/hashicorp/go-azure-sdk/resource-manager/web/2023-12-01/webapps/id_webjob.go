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
	recaser.RegisterResourceId(&WebJobId{})
}

var _ resourceids.ResourceId = &WebJobId{}

// WebJobId is a struct representing the Resource ID for a Web Job
type WebJobId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	WebJobName        string
}

// NewWebJobID returns a new WebJobId struct
func NewWebJobID(subscriptionId string, resourceGroupName string, siteName string, webJobName string) WebJobId {
	return WebJobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		WebJobName:        webJobName,
	}
}

// ParseWebJobID parses 'input' into a WebJobId
func ParseWebJobID(input string) (*WebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WebJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWebJobIDInsensitively parses 'input' case-insensitively into a WebJobId
// note: this method should only be used for API response data and not user input
func ParseWebJobIDInsensitively(input string) (*WebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WebJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WebJobId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.WebJobName, ok = input.Parsed["webJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "webJobName", input)
	}

	return nil
}

// ValidateWebJobID checks that 'input' can be parsed as a Web Job ID
func ValidateWebJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWebJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Web Job ID
func (id WebJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/webJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.WebJobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Web Job ID
func (id WebJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticWebJobs", "webJobs", "webJobs"),
		resourceids.UserSpecifiedSegment("webJobName", "webJobName"),
	}
}

// String returns a human-readable description of this Web Job ID
func (id WebJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Web Job Name: %q", id.WebJobName),
	}
	return fmt.Sprintf("Web Job (%s)", strings.Join(components, "\n"))
}
