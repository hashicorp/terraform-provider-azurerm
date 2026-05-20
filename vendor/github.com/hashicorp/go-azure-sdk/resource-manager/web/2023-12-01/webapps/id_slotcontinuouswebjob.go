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
	recaser.RegisterResourceId(&SlotContinuousWebJobId{})
}

var _ resourceids.ResourceId = &SlotContinuousWebJobId{}

// SlotContinuousWebJobId is a struct representing the Resource ID for a Slot Continuous Web Job
type SlotContinuousWebJobId struct {
	SubscriptionId       string
	ResourceGroupName    string
	SiteName             string
	SlotName             string
	ContinuousWebJobName string
}

// NewSlotContinuousWebJobID returns a new SlotContinuousWebJobId struct
func NewSlotContinuousWebJobID(subscriptionId string, resourceGroupName string, siteName string, slotName string, continuousWebJobName string) SlotContinuousWebJobId {
	return SlotContinuousWebJobId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		SiteName:             siteName,
		SlotName:             slotName,
		ContinuousWebJobName: continuousWebJobName,
	}
}

// ParseSlotContinuousWebJobID parses 'input' into a SlotContinuousWebJobId
func ParseSlotContinuousWebJobID(input string) (*SlotContinuousWebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotContinuousWebJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotContinuousWebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotContinuousWebJobIDInsensitively parses 'input' case-insensitively into a SlotContinuousWebJobId
// note: this method should only be used for API response data and not user input
func ParseSlotContinuousWebJobIDInsensitively(input string) (*SlotContinuousWebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotContinuousWebJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotContinuousWebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotContinuousWebJobId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ContinuousWebJobName, ok = input.Parsed["continuousWebJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "continuousWebJobName", input)
	}

	return nil
}

// ValidateSlotContinuousWebJobID checks that 'input' can be parsed as a Slot Continuous Web Job ID
func ValidateSlotContinuousWebJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotContinuousWebJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Continuous Web Job ID
func (id SlotContinuousWebJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/continuousWebJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.ContinuousWebJobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Continuous Web Job ID
func (id SlotContinuousWebJobId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticContinuousWebJobs", "continuousWebJobs", "continuousWebJobs"),
		resourceids.UserSpecifiedSegment("continuousWebJobName", "continuousWebJobName"),
	}
}

// String returns a human-readable description of this Slot Continuous Web Job ID
func (id SlotContinuousWebJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Continuous Web Job Name: %q", id.ContinuousWebJobName),
	}
	return fmt.Sprintf("Slot Continuous Web Job (%s)", strings.Join(components, "\n"))
}
