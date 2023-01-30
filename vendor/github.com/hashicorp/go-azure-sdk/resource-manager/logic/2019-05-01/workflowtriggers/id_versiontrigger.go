package workflowtriggers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VersionTriggerId{}

// VersionTriggerId is a struct representing the Resource ID for a Version Trigger
type VersionTriggerId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkflowName      string
	VersionId         string
	TriggerName       string
}

// NewVersionTriggerID returns a new VersionTriggerId struct
func NewVersionTriggerID(subscriptionId string, resourceGroupName string, workflowName string, versionId string, triggerName string) VersionTriggerId {
	return VersionTriggerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkflowName:      workflowName,
		VersionId:         versionId,
		TriggerName:       triggerName,
	}
}

// ParseVersionTriggerID parses 'input' into a VersionTriggerId
func ParseVersionTriggerID(input string) (*VersionTriggerId, error) {
	parser := resourceids.NewParserFromResourceIdType(VersionTriggerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VersionTriggerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkflowName, ok = parsed.Parsed["workflowName"]; !ok {
		return nil, fmt.Errorf("the segment 'workflowName' was not found in the resource id %q", input)
	}

	if id.VersionId, ok = parsed.Parsed["versionId"]; !ok {
		return nil, fmt.Errorf("the segment 'versionId' was not found in the resource id %q", input)
	}

	if id.TriggerName, ok = parsed.Parsed["triggerName"]; !ok {
		return nil, fmt.Errorf("the segment 'triggerName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseVersionTriggerIDInsensitively parses 'input' case-insensitively into a VersionTriggerId
// note: this method should only be used for API response data and not user input
func ParseVersionTriggerIDInsensitively(input string) (*VersionTriggerId, error) {
	parser := resourceids.NewParserFromResourceIdType(VersionTriggerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VersionTriggerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkflowName, ok = parsed.Parsed["workflowName"]; !ok {
		return nil, fmt.Errorf("the segment 'workflowName' was not found in the resource id %q", input)
	}

	if id.VersionId, ok = parsed.Parsed["versionId"]; !ok {
		return nil, fmt.Errorf("the segment 'versionId' was not found in the resource id %q", input)
	}

	if id.TriggerName, ok = parsed.Parsed["triggerName"]; !ok {
		return nil, fmt.Errorf("the segment 'triggerName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateVersionTriggerID checks that 'input' can be parsed as a Version Trigger ID
func ValidateVersionTriggerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVersionTriggerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Version Trigger ID
func (id VersionTriggerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s/versions/%s/triggers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkflowName, id.VersionId, id.TriggerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Version Trigger ID
func (id VersionTriggerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticWorkflows", "workflows", "workflows"),
		resourceids.UserSpecifiedSegment("workflowName", "workflowValue"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionId", "versionIdValue"),
		resourceids.StaticSegment("staticTriggers", "triggers", "triggers"),
		resourceids.UserSpecifiedSegment("triggerName", "triggerValue"),
	}
}

// String returns a human-readable description of this Version Trigger ID
func (id VersionTriggerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workflow Name: %q", id.WorkflowName),
		fmt.Sprintf("Version: %q", id.VersionId),
		fmt.Sprintf("Trigger Name: %q", id.TriggerName),
	}
	return fmt.Sprintf("Version Trigger (%s)", strings.Join(components, "\n"))
}
