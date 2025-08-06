package datacollectionrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DataCollectionRuleId{})
}

var _ resourceids.ResourceId = &DataCollectionRuleId{}

// DataCollectionRuleId is a struct representing the Resource ID for a Data Collection Rule
type DataCollectionRuleId struct {
	SubscriptionId         string
	ResourceGroupName      string
	DataCollectionRuleName string
}

// NewDataCollectionRuleID returns a new DataCollectionRuleId struct
func NewDataCollectionRuleID(subscriptionId string, resourceGroupName string, dataCollectionRuleName string) DataCollectionRuleId {
	return DataCollectionRuleId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		DataCollectionRuleName: dataCollectionRuleName,
	}
}

// ParseDataCollectionRuleID parses 'input' into a DataCollectionRuleId
func ParseDataCollectionRuleID(input string) (*DataCollectionRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataCollectionRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataCollectionRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataCollectionRuleIDInsensitively parses 'input' case-insensitively into a DataCollectionRuleId
// note: this method should only be used for API response data and not user input
func ParseDataCollectionRuleIDInsensitively(input string) (*DataCollectionRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataCollectionRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataCollectionRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataCollectionRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DataCollectionRuleName, ok = input.Parsed["dataCollectionRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataCollectionRuleName", input)
	}

	return nil
}

// ValidateDataCollectionRuleID checks that 'input' can be parsed as a Data Collection Rule ID
func ValidateDataCollectionRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataCollectionRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Collection Rule ID
func (id DataCollectionRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/dataCollectionRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DataCollectionRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Collection Rule ID
func (id DataCollectionRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticDataCollectionRules", "dataCollectionRules", "dataCollectionRules"),
		resourceids.UserSpecifiedSegment("dataCollectionRuleName", "dataCollectionRuleName"),
	}
}

// String returns a human-readable description of this Data Collection Rule ID
func (id DataCollectionRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Data Collection Rule Name: %q", id.DataCollectionRuleName),
	}
	return fmt.Sprintf("Data Collection Rule (%s)", strings.Join(components, "\n"))
}
