// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &AutomationCompilationJobId{}

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
	parser := resourceids.NewParserFromResourceIdType(&AutomationCompilationJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutomationCompilationJobId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutomationCompilationJobIDInsensitively parses 'input' case-insensitively into a AutomationCompilationJobId
// note: this method should only be used for API response data and not user input
func ParseAutomationCompilationJobIDInsensitively(input string) (*AutomationCompilationJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutomationCompilationJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutomationCompilationJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutomationCompilationJobId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.CompilationJobId, ok = input.Parsed["compilationJobId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "compilationJobId", input)
	}

	return nil
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
