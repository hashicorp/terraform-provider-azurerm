package hybridrunbookworker

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&HybridRunbookWorkerGroupId{})
}

var _ resourceids.ResourceId = &HybridRunbookWorkerGroupId{}

// HybridRunbookWorkerGroupId is a struct representing the Resource ID for a Hybrid Runbook Worker Group
type HybridRunbookWorkerGroupId struct {
	SubscriptionId               string
	ResourceGroupName            string
	AutomationAccountName        string
	HybridRunbookWorkerGroupName string
}

// NewHybridRunbookWorkerGroupID returns a new HybridRunbookWorkerGroupId struct
func NewHybridRunbookWorkerGroupID(subscriptionId string, resourceGroupName string, automationAccountName string, hybridRunbookWorkerGroupName string) HybridRunbookWorkerGroupId {
	return HybridRunbookWorkerGroupId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		AutomationAccountName:        automationAccountName,
		HybridRunbookWorkerGroupName: hybridRunbookWorkerGroupName,
	}
}

// ParseHybridRunbookWorkerGroupID parses 'input' into a HybridRunbookWorkerGroupId
func ParseHybridRunbookWorkerGroupID(input string) (*HybridRunbookWorkerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridRunbookWorkerGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridRunbookWorkerGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHybridRunbookWorkerGroupIDInsensitively parses 'input' case-insensitively into a HybridRunbookWorkerGroupId
// note: this method should only be used for API response data and not user input
func ParseHybridRunbookWorkerGroupIDInsensitively(input string) (*HybridRunbookWorkerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridRunbookWorkerGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridRunbookWorkerGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HybridRunbookWorkerGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AutomationAccountName, ok = input.Parsed["automationAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", input)
	}

	if id.HybridRunbookWorkerGroupName, ok = input.Parsed["hybridRunbookWorkerGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hybridRunbookWorkerGroupName", input)
	}

	return nil
}

// ValidateHybridRunbookWorkerGroupID checks that 'input' can be parsed as a Hybrid Runbook Worker Group ID
func ValidateHybridRunbookWorkerGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHybridRunbookWorkerGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hybrid Runbook Worker Group ID
func (id HybridRunbookWorkerGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/hybridRunbookWorkerGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.HybridRunbookWorkerGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hybrid Runbook Worker Group ID
func (id HybridRunbookWorkerGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticHybridRunbookWorkerGroups", "hybridRunbookWorkerGroups", "hybridRunbookWorkerGroups"),
		resourceids.UserSpecifiedSegment("hybridRunbookWorkerGroupName", "hybridRunbookWorkerGroupName"),
	}
}

// String returns a human-readable description of this Hybrid Runbook Worker Group ID
func (id HybridRunbookWorkerGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Hybrid Runbook Worker Group Name: %q", id.HybridRunbookWorkerGroupName),
	}
	return fmt.Sprintf("Hybrid Runbook Worker Group (%s)", strings.Join(components, "\n"))
}
