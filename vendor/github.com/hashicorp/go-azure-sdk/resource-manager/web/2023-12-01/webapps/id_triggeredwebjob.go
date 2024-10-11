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
	recaser.RegisterResourceId(&TriggeredWebJobId{})
}

var _ resourceids.ResourceId = &TriggeredWebJobId{}

// TriggeredWebJobId is a struct representing the Resource ID for a Triggered Web Job
type TriggeredWebJobId struct {
	SubscriptionId      string
	ResourceGroupName   string
	SiteName            string
	TriggeredWebJobName string
}

// NewTriggeredWebJobID returns a new TriggeredWebJobId struct
func NewTriggeredWebJobID(subscriptionId string, resourceGroupName string, siteName string, triggeredWebJobName string) TriggeredWebJobId {
	return TriggeredWebJobId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		SiteName:            siteName,
		TriggeredWebJobName: triggeredWebJobName,
	}
}

// ParseTriggeredWebJobID parses 'input' into a TriggeredWebJobId
func ParseTriggeredWebJobID(input string) (*TriggeredWebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TriggeredWebJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TriggeredWebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTriggeredWebJobIDInsensitively parses 'input' case-insensitively into a TriggeredWebJobId
// note: this method should only be used for API response data and not user input
func ParseTriggeredWebJobIDInsensitively(input string) (*TriggeredWebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TriggeredWebJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TriggeredWebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TriggeredWebJobId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TriggeredWebJobName, ok = input.Parsed["triggeredWebJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "triggeredWebJobName", input)
	}

	return nil
}

// ValidateTriggeredWebJobID checks that 'input' can be parsed as a Triggered Web Job ID
func ValidateTriggeredWebJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTriggeredWebJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Triggered Web Job ID
func (id TriggeredWebJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/triggeredWebJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.TriggeredWebJobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Triggered Web Job ID
func (id TriggeredWebJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticTriggeredWebJobs", "triggeredWebJobs", "triggeredWebJobs"),
		resourceids.UserSpecifiedSegment("triggeredWebJobName", "triggeredWebJobName"),
	}
}

// String returns a human-readable description of this Triggered Web Job ID
func (id TriggeredWebJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Triggered Web Job Name: %q", id.TriggeredWebJobName),
	}
	return fmt.Sprintf("Triggered Web Job (%s)", strings.Join(components, "\n"))
}
