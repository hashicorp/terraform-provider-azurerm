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
	recaser.RegisterResourceId(&HybridRunbookWorkerId{})
}

var _ resourceids.ResourceId = &HybridRunbookWorkerId{}

// HybridRunbookWorkerId is a struct representing the Resource ID for a Hybrid Runbook Worker
type HybridRunbookWorkerId struct {
	SubscriptionId               string
	ResourceGroupName            string
	AutomationAccountName        string
	HybridRunbookWorkerGroupName string
	HybridRunbookWorkerId        string
}

// NewHybridRunbookWorkerID returns a new HybridRunbookWorkerId struct
func NewHybridRunbookWorkerID(subscriptionId string, resourceGroupName string, automationAccountName string, hybridRunbookWorkerGroupName string, hybridRunbookWorkerId string) HybridRunbookWorkerId {
	return HybridRunbookWorkerId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		AutomationAccountName:        automationAccountName,
		HybridRunbookWorkerGroupName: hybridRunbookWorkerGroupName,
		HybridRunbookWorkerId:        hybridRunbookWorkerId,
	}
}

// ParseHybridRunbookWorkerID parses 'input' into a HybridRunbookWorkerId
func ParseHybridRunbookWorkerID(input string) (*HybridRunbookWorkerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridRunbookWorkerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridRunbookWorkerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHybridRunbookWorkerIDInsensitively parses 'input' case-insensitively into a HybridRunbookWorkerId
// note: this method should only be used for API response data and not user input
func ParseHybridRunbookWorkerIDInsensitively(input string) (*HybridRunbookWorkerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridRunbookWorkerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridRunbookWorkerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HybridRunbookWorkerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.HybridRunbookWorkerId, ok = input.Parsed["hybridRunbookWorkerId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hybridRunbookWorkerId", input)
	}

	return nil
}

// ValidateHybridRunbookWorkerID checks that 'input' can be parsed as a Hybrid Runbook Worker ID
func ValidateHybridRunbookWorkerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHybridRunbookWorkerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hybrid Runbook Worker ID
func (id HybridRunbookWorkerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/hybridRunbookWorkerGroups/%s/hybridRunbookWorkers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.HybridRunbookWorkerGroupName, id.HybridRunbookWorkerId)
}

// Segments returns a slice of Resource ID Segments which comprise this Hybrid Runbook Worker ID
func (id HybridRunbookWorkerId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticHybridRunbookWorkers", "hybridRunbookWorkers", "hybridRunbookWorkers"),
		resourceids.UserSpecifiedSegment("hybridRunbookWorkerId", "hybridRunbookWorkerId"),
	}
}

// String returns a human-readable description of this Hybrid Runbook Worker ID
func (id HybridRunbookWorkerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Hybrid Runbook Worker Group Name: %q", id.HybridRunbookWorkerGroupName),
		fmt.Sprintf("Hybrid Runbook Worker: %q", id.HybridRunbookWorkerId),
	}
	return fmt.Sprintf("Hybrid Runbook Worker (%s)", strings.Join(components, "\n"))
}
