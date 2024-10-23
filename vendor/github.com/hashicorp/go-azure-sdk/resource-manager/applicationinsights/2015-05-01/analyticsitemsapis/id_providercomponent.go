package analyticsitemsapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProviderComponentId{})
}

var _ resourceids.ResourceId = &ProviderComponentId{}

// ProviderComponentId is a struct representing the Resource ID for a Provider Component
type ProviderComponentId struct {
	SubscriptionId    string
	ResourceGroupName string
	ComponentName     string
	ScopePath         string
}

// NewProviderComponentID returns a new ProviderComponentId struct
func NewProviderComponentID(subscriptionId string, resourceGroupName string, componentName string, scopePath string) ProviderComponentId {
	return ProviderComponentId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ComponentName:     componentName,
		ScopePath:         scopePath,
	}
}

// ParseProviderComponentID parses 'input' into a ProviderComponentId
func ParseProviderComponentID(input string) (*ProviderComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderComponentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderComponentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviderComponentIDInsensitively parses 'input' case-insensitively into a ProviderComponentId
// note: this method should only be used for API response data and not user input
func ParseProviderComponentIDInsensitively(input string) (*ProviderComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderComponentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderComponentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProviderComponentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ComponentName, ok = input.Parsed["componentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "componentName", input)
	}

	if id.ScopePath, ok = input.Parsed["scopePath"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scopePath", input)
	}

	return nil
}

// ValidateProviderComponentID checks that 'input' can be parsed as a Provider Component ID
func ValidateProviderComponentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderComponentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Component ID
func (id ProviderComponentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/components/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ComponentName, strings.TrimPrefix(id.ScopePath, "/"))
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Component ID
func (id ProviderComponentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticComponents", "components", "components"),
		resourceids.UserSpecifiedSegment("componentName", "componentName"),
		resourceids.ScopeSegment("scopePath", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
	}
}

// String returns a human-readable description of this Provider Component ID
func (id ProviderComponentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Component Name: %q", id.ComponentName),
		fmt.Sprintf("Scope Path: %q", id.ScopePath),
	}
	return fmt.Sprintf("Provider Component (%s)", strings.Join(components, "\n"))
}
