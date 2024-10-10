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
	recaser.RegisterResourceId(&SlotDeploymentId{})
}

var _ resourceids.ResourceId = &SlotDeploymentId{}

// SlotDeploymentId is a struct representing the Resource ID for a Slot Deployment
type SlotDeploymentId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	DeploymentName    string
}

// NewSlotDeploymentID returns a new SlotDeploymentId struct
func NewSlotDeploymentID(subscriptionId string, resourceGroupName string, siteName string, slotName string, deploymentName string) SlotDeploymentId {
	return SlotDeploymentId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		DeploymentName:    deploymentName,
	}
}

// ParseSlotDeploymentID parses 'input' into a SlotDeploymentId
func ParseSlotDeploymentID(input string) (*SlotDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotDeploymentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotDeploymentIDInsensitively parses 'input' case-insensitively into a SlotDeploymentId
// note: this method should only be used for API response data and not user input
func ParseSlotDeploymentIDInsensitively(input string) (*SlotDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotDeploymentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotDeploymentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DeploymentName, ok = input.Parsed["deploymentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deploymentName", input)
	}

	return nil
}

// ValidateSlotDeploymentID checks that 'input' can be parsed as a Slot Deployment ID
func ValidateSlotDeploymentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotDeploymentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Deployment ID
func (id SlotDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.DeploymentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Deployment ID
func (id SlotDeploymentId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticDeployments", "deployments", "deployments"),
		resourceids.UserSpecifiedSegment("deploymentName", "deploymentName"),
	}
}

// String returns a human-readable description of this Slot Deployment ID
func (id SlotDeploymentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Deployment Name: %q", id.DeploymentName),
	}
	return fmt.Sprintf("Slot Deployment (%s)", strings.Join(components, "\n"))
}
