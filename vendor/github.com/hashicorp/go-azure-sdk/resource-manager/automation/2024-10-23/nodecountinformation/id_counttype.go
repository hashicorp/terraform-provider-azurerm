package nodecountinformation

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CountTypeId{})
}

var _ resourceids.ResourceId = &CountTypeId{}

// CountTypeId is a struct representing the Resource ID for a Count Type
type CountTypeId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	CountType             CountType
}

// NewCountTypeID returns a new CountTypeId struct
func NewCountTypeID(subscriptionId string, resourceGroupName string, automationAccountName string, countType CountType) CountTypeId {
	return CountTypeId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		CountType:             countType,
	}
}

// ParseCountTypeID parses 'input' into a CountTypeId
func ParseCountTypeID(input string) (*CountTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CountTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CountTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCountTypeIDInsensitively parses 'input' case-insensitively into a CountTypeId
// note: this method should only be used for API response data and not user input
func ParseCountTypeIDInsensitively(input string) (*CountTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CountTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CountTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CountTypeId) FromParseResult(input resourceids.ParseResult) error {
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

	if v, ok := input.Parsed["countType"]; true {
		if !ok {
			return resourceids.NewSegmentNotSpecifiedError(id, "countType", input)
		}

		countType, err := parseCountType(v)
		if err != nil {
			return fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.CountType = *countType
	}

	return nil
}

// ValidateCountTypeID checks that 'input' can be parsed as a Count Type ID
func ValidateCountTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCountTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Count Type ID
func (id CountTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/nodeCounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, string(id.CountType))
}

// Segments returns a slice of Resource ID Segments which comprise this Count Type ID
func (id CountTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticNodeCounts", "nodeCounts", "nodeCounts"),
		resourceids.ConstantSegment("countType", PossibleValuesForCountType(), "nodeconfiguration"),
	}
}

// String returns a human-readable description of this Count Type ID
func (id CountTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Count Type: %q", string(id.CountType)),
	}
	return fmt.Sprintf("Count Type (%s)", strings.Join(components, "\n"))
}
