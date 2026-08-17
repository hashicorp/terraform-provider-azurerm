package sourcecontrolsyncjob

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SourceControlSyncJobId{})
}

var _ resourceids.ResourceId = &SourceControlSyncJobId{}

// SourceControlSyncJobId is a struct representing the Resource ID for a Source Control Sync Job
type SourceControlSyncJobId struct {
	SubscriptionId         string
	ResourceGroupName      string
	AutomationAccountName  string
	SourceControlName      string
	SourceControlSyncJobId string
}

// NewSourceControlSyncJobID returns a new SourceControlSyncJobId struct
func NewSourceControlSyncJobID(subscriptionId string, resourceGroupName string, automationAccountName string, sourceControlName string, sourceControlSyncJobId string) SourceControlSyncJobId {
	return SourceControlSyncJobId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		AutomationAccountName:  automationAccountName,
		SourceControlName:      sourceControlName,
		SourceControlSyncJobId: sourceControlSyncJobId,
	}
}

// ParseSourceControlSyncJobID parses 'input' into a SourceControlSyncJobId
func ParseSourceControlSyncJobID(input string) (*SourceControlSyncJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SourceControlSyncJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SourceControlSyncJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSourceControlSyncJobIDInsensitively parses 'input' case-insensitively into a SourceControlSyncJobId
// note: this method should only be used for API response data and not user input
func ParseSourceControlSyncJobIDInsensitively(input string) (*SourceControlSyncJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SourceControlSyncJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SourceControlSyncJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SourceControlSyncJobId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SourceControlName, ok = input.Parsed["sourceControlName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sourceControlName", input)
	}

	if id.SourceControlSyncJobId, ok = input.Parsed["sourceControlSyncJobId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sourceControlSyncJobId", input)
	}

	return nil
}

// ValidateSourceControlSyncJobID checks that 'input' can be parsed as a Source Control Sync Job ID
func ValidateSourceControlSyncJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSourceControlSyncJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Source Control Sync Job ID
func (id SourceControlSyncJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/sourceControls/%s/sourceControlSyncJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.SourceControlName, id.SourceControlSyncJobId)
}

// Segments returns a slice of Resource ID Segments which comprise this Source Control Sync Job ID
func (id SourceControlSyncJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticSourceControls", "sourceControls", "sourceControls"),
		resourceids.UserSpecifiedSegment("sourceControlName", "sourceControlName"),
		resourceids.StaticSegment("staticSourceControlSyncJobs", "sourceControlSyncJobs", "sourceControlSyncJobs"),
		resourceids.UserSpecifiedSegment("sourceControlSyncJobId", "sourceControlSyncJobId"),
	}
}

// String returns a human-readable description of this Source Control Sync Job ID
func (id SourceControlSyncJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Source Control Name: %q", id.SourceControlName),
		fmt.Sprintf("Source Control Sync Job: %q", id.SourceControlSyncJobId),
	}
	return fmt.Sprintf("Source Control Sync Job (%s)", strings.Join(components, "\n"))
}
