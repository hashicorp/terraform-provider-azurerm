package rulesets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = RuleSetId{}

// RuleSetId is a struct representing the Resource ID for a Rule Set
type RuleSetId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	RuleSetName       string
}

// NewRuleSetID returns a new RuleSetId struct
func NewRuleSetID(subscriptionId string, resourceGroupName string, profileName string, ruleSetName string) RuleSetId {
	return RuleSetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		RuleSetName:       ruleSetName,
	}
}

// ParseRuleSetID parses 'input' into a RuleSetId
func ParseRuleSetID(input string) (*RuleSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(RuleSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RuleSetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.RuleSetName, ok = parsed.Parsed["ruleSetName"]; !ok {
		return nil, fmt.Errorf("the segment 'ruleSetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseRuleSetIDInsensitively parses 'input' case-insensitively into a RuleSetId
// note: this method should only be used for API response data and not user input
func ParseRuleSetIDInsensitively(input string) (*RuleSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(RuleSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RuleSetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.RuleSetName, ok = parsed.Parsed["ruleSetName"]; !ok {
		return nil, fmt.Errorf("the segment 'ruleSetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateRuleSetID checks that 'input' can be parsed as a Rule Set ID
func ValidateRuleSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRuleSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rule Set ID
func (id RuleSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/ruleSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rule Set ID
func (id RuleSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCdn", "Microsoft.Cdn", "Microsoft.Cdn"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileValue"),
		resourceids.StaticSegment("staticRuleSets", "ruleSets", "ruleSets"),
		resourceids.UserSpecifiedSegment("ruleSetName", "ruleSetValue"),
	}
}

// String returns a human-readable description of this Rule Set ID
func (id RuleSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Rule Set Name: %q", id.RuleSetName),
	}
	return fmt.Sprintf("Rule Set (%s)", strings.Join(components, "\n"))
}
