package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = TagRuleId{}

// TagRuleId is a struct representing the Resource ID for a Tag Rule
type TagRuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	MonitorName       string
	RuleSetName       string
}

// NewTagRuleID returns a new TagRuleId struct
func NewTagRuleID(subscriptionId string, resourceGroupName string, monitorName string, ruleSetName string) TagRuleId {
	return TagRuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MonitorName:       monitorName,
		RuleSetName:       ruleSetName,
	}
}

// ParseTagRuleID parses 'input' into a TagRuleId
func ParseTagRuleID(input string) (*TagRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(TagRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TagRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.MonitorName, ok = parsed.Parsed["monitorName"]; !ok {
		return nil, fmt.Errorf("the segment 'monitorName' was not found in the resource id %q", input)
	}

	if id.RuleSetName, ok = parsed.Parsed["ruleSetName"]; !ok {
		return nil, fmt.Errorf("the segment 'ruleSetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseTagRuleIDInsensitively parses 'input' case-insensitively into a TagRuleId
// note: this method should only be used for API response data and not user input
func ParseTagRuleIDInsensitively(input string) (*TagRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(TagRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TagRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.MonitorName, ok = parsed.Parsed["monitorName"]; !ok {
		return nil, fmt.Errorf("the segment 'monitorName' was not found in the resource id %q", input)
	}

	if id.RuleSetName, ok = parsed.Parsed["ruleSetName"]; !ok {
		return nil, fmt.Errorf("the segment 'ruleSetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateTagRuleID checks that 'input' can be parsed as a Tag Rule ID
func ValidateTagRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTagRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Tag Rule ID
func (id TagRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Elastic/monitors/%s/tagRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MonitorName, id.RuleSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Tag Rule ID
func (id TagRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftElastic", "Microsoft.Elastic", "Microsoft.Elastic"),
		resourceids.StaticSegment("staticMonitors", "monitors", "monitors"),
		resourceids.UserSpecifiedSegment("monitorName", "monitorValue"),
		resourceids.StaticSegment("staticTagRules", "tagRules", "tagRules"),
		resourceids.UserSpecifiedSegment("ruleSetName", "ruleSetValue"),
	}
}

// String returns a human-readable description of this Tag Rule ID
func (id TagRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Monitor Name: %q", id.MonitorName),
		fmt.Sprintf("Rule Set Name: %q", id.RuleSetName),
	}
	return fmt.Sprintf("Tag Rule (%s)", strings.Join(components, "\n"))
}
