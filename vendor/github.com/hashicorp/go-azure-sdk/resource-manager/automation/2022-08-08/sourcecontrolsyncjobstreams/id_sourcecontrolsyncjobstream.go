package sourcecontrolsyncjobstreams

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SourceControlSyncJobStreamId{}

// SourceControlSyncJobStreamId is a struct representing the Resource ID for a Source Control Sync Job Stream
type SourceControlSyncJobStreamId struct {
	SubscriptionId         string
	ResourceGroupName      string
	AutomationAccountName  string
	SourceControlName      string
	SourceControlSyncJobId string
	StreamId               string
}

// NewSourceControlSyncJobStreamID returns a new SourceControlSyncJobStreamId struct
func NewSourceControlSyncJobStreamID(subscriptionId string, resourceGroupName string, automationAccountName string, sourceControlName string, sourceControlSyncJobId string, streamId string) SourceControlSyncJobStreamId {
	return SourceControlSyncJobStreamId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		AutomationAccountName:  automationAccountName,
		SourceControlName:      sourceControlName,
		SourceControlSyncJobId: sourceControlSyncJobId,
		StreamId:               streamId,
	}
}

// ParseSourceControlSyncJobStreamID parses 'input' into a SourceControlSyncJobStreamId
func ParseSourceControlSyncJobStreamID(input string) (*SourceControlSyncJobStreamId, error) {
	parser := resourceids.NewParserFromResourceIdType(SourceControlSyncJobStreamId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SourceControlSyncJobStreamId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.SourceControlName, ok = parsed.Parsed["sourceControlName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sourceControlName", *parsed)
	}

	if id.SourceControlSyncJobId, ok = parsed.Parsed["sourceControlSyncJobId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sourceControlSyncJobId", *parsed)
	}

	if id.StreamId, ok = parsed.Parsed["streamId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "streamId", *parsed)
	}

	return &id, nil
}

// ParseSourceControlSyncJobStreamIDInsensitively parses 'input' case-insensitively into a SourceControlSyncJobStreamId
// note: this method should only be used for API response data and not user input
func ParseSourceControlSyncJobStreamIDInsensitively(input string) (*SourceControlSyncJobStreamId, error) {
	parser := resourceids.NewParserFromResourceIdType(SourceControlSyncJobStreamId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SourceControlSyncJobStreamId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.SourceControlName, ok = parsed.Parsed["sourceControlName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sourceControlName", *parsed)
	}

	if id.SourceControlSyncJobId, ok = parsed.Parsed["sourceControlSyncJobId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sourceControlSyncJobId", *parsed)
	}

	if id.StreamId, ok = parsed.Parsed["streamId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "streamId", *parsed)
	}

	return &id, nil
}

// ValidateSourceControlSyncJobStreamID checks that 'input' can be parsed as a Source Control Sync Job Stream ID
func ValidateSourceControlSyncJobStreamID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSourceControlSyncJobStreamID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Source Control Sync Job Stream ID
func (id SourceControlSyncJobStreamId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/sourceControls/%s/sourceControlSyncJobs/%s/streams/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.SourceControlName, id.SourceControlSyncJobId, id.StreamId)
}

// Segments returns a slice of Resource ID Segments which comprise this Source Control Sync Job Stream ID
func (id SourceControlSyncJobStreamId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticSourceControls", "sourceControls", "sourceControls"),
		resourceids.UserSpecifiedSegment("sourceControlName", "sourceControlValue"),
		resourceids.StaticSegment("staticSourceControlSyncJobs", "sourceControlSyncJobs", "sourceControlSyncJobs"),
		resourceids.UserSpecifiedSegment("sourceControlSyncJobId", "sourceControlSyncJobIdValue"),
		resourceids.StaticSegment("staticStreams", "streams", "streams"),
		resourceids.UserSpecifiedSegment("streamId", "streamIdValue"),
	}
}

// String returns a human-readable description of this Source Control Sync Job Stream ID
func (id SourceControlSyncJobStreamId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Source Control Name: %q", id.SourceControlName),
		fmt.Sprintf("Source Control Sync Job: %q", id.SourceControlSyncJobId),
		fmt.Sprintf("Stream: %q", id.StreamId),
	}
	return fmt.Sprintf("Source Control Sync Job Stream (%s)", strings.Join(components, "\n"))
}
