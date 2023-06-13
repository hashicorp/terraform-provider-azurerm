// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AutomationCompilationJobId{}

// AutomationCompilationJobId is a struct representing the Resource ID for a Compilation Job
type AutomationCompilationJobId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	CompilationJobId      string
}

// NewAutomationCompilationJobID returns a new AutomationCompilationJobId struct
func NewAutomationCompilationJobID(subscriptionId string, resourceGroupName string, automationAccountName string, compilationJobId string) AutomationCompilationJobId {
	return AutomationCompilationJobId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		CompilationJobId:      compilationJobId,
	}
}

// ParseAutomationCompilationJobID parses 'input' into a AutomationCompilationJobId
func ParseAutomationCompilationJobID(input string) (*AutomationCompilationJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(AutomationCompilationJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AutomationCompilationJobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.CompilationJobId, ok = parsed.Parsed["compilationJobId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "compilationJobId", *parsed)
	}

	return &id, nil
}

// ParseAutomationCompilationJobIDInsensitively parses 'input' case-insensitively into a AutomationCompilationJobId
// note: this method should only be used for API response data and not user input
func ParseAutomationCompilationJobIDInsensitively(input string) (*AutomationCompilationJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(AutomationCompilationJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AutomationCompilationJobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.CompilationJobId, ok = parsed.Parsed["compilationJobId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "compilationJobId", *parsed)
	}

	return &id, nil
}

// ValidateAutomationCompilationJobID checks that 'input' can be parsed as a Compilation Job ID
func ValidateAutomationCompilationJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutomationCompilationJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Compilation Job ID
func (id AutomationCompilationJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/compilationJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.CompilationJobId)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud Services I P Configuration ID
func (id AutomationCompilationJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticCompilationJobs", "compilationJobs", "compilationJobs"),
		resourceids.UserSpecifiedSegment("compilationJobId", "compilationJobIdValue"),
	}
}

// String returns a human-readable description of this Compilation Job ID
func (id AutomationCompilationJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Compilation Job: %q", id.CompilationJobId),
	}
	return fmt.Sprintf("Compilation Job (%s)", strings.Join(components, "\n"))
}
