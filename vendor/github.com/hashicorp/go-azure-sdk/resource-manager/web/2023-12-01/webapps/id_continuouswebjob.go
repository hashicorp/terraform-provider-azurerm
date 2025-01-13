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
	recaser.RegisterResourceId(&ContinuousWebJobId{})
}

var _ resourceids.ResourceId = &ContinuousWebJobId{}

// ContinuousWebJobId is a struct representing the Resource ID for a Continuous Web Job
type ContinuousWebJobId struct {
	SubscriptionId       string
	ResourceGroupName    string
	SiteName             string
	ContinuousWebJobName string
}

// NewContinuousWebJobID returns a new ContinuousWebJobId struct
func NewContinuousWebJobID(subscriptionId string, resourceGroupName string, siteName string, continuousWebJobName string) ContinuousWebJobId {
	return ContinuousWebJobId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		SiteName:             siteName,
		ContinuousWebJobName: continuousWebJobName,
	}
}

// ParseContinuousWebJobID parses 'input' into a ContinuousWebJobId
func ParseContinuousWebJobID(input string) (*ContinuousWebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContinuousWebJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContinuousWebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseContinuousWebJobIDInsensitively parses 'input' case-insensitively into a ContinuousWebJobId
// note: this method should only be used for API response data and not user input
func ParseContinuousWebJobIDInsensitively(input string) (*ContinuousWebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContinuousWebJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContinuousWebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ContinuousWebJobId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ContinuousWebJobName, ok = input.Parsed["continuousWebJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "continuousWebJobName", input)
	}

	return nil
}

// ValidateContinuousWebJobID checks that 'input' can be parsed as a Continuous Web Job ID
func ValidateContinuousWebJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContinuousWebJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Continuous Web Job ID
func (id ContinuousWebJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/continuousWebJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.ContinuousWebJobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Continuous Web Job ID
func (id ContinuousWebJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticContinuousWebJobs", "continuousWebJobs", "continuousWebJobs"),
		resourceids.UserSpecifiedSegment("continuousWebJobName", "continuousWebJobName"),
	}
}

// String returns a human-readable description of this Continuous Web Job ID
func (id ContinuousWebJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Continuous Web Job Name: %q", id.ContinuousWebJobName),
	}
	return fmt.Sprintf("Continuous Web Job (%s)", strings.Join(components, "\n"))
}
