package quotabyperiodkeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PeriodId{})
}

var _ resourceids.ResourceId = &PeriodId{}

// PeriodId is a struct representing the Resource ID for a Period
type PeriodId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	QuotaCounterKey   string
	QuotaPeriodKey    string
}

// NewPeriodID returns a new PeriodId struct
func NewPeriodID(subscriptionId string, resourceGroupName string, serviceName string, quotaCounterKey string, quotaPeriodKey string) PeriodId {
	return PeriodId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		QuotaCounterKey:   quotaCounterKey,
		QuotaPeriodKey:    quotaPeriodKey,
	}
}

// ParsePeriodID parses 'input' into a PeriodId
func ParsePeriodID(input string) (*PeriodId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PeriodId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeriodId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePeriodIDInsensitively parses 'input' case-insensitively into a PeriodId
// note: this method should only be used for API response data and not user input
func ParsePeriodIDInsensitively(input string) (*PeriodId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PeriodId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeriodId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PeriodId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.QuotaCounterKey, ok = input.Parsed["quotaCounterKey"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "quotaCounterKey", input)
	}

	if id.QuotaPeriodKey, ok = input.Parsed["quotaPeriodKey"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "quotaPeriodKey", input)
	}

	return nil
}

// ValidatePeriodID checks that 'input' can be parsed as a Period ID
func ValidatePeriodID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePeriodID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Period ID
func (id PeriodId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/quotas/%s/periods/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.QuotaCounterKey, id.QuotaPeriodKey)
}

// Segments returns a slice of Resource ID Segments which comprise this Period ID
func (id PeriodId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticQuotas", "quotas", "quotas"),
		resourceids.UserSpecifiedSegment("quotaCounterKey", "quotaCounterKey"),
		resourceids.StaticSegment("staticPeriods", "periods", "periods"),
		resourceids.UserSpecifiedSegment("quotaPeriodKey", "quotaPeriodKey"),
	}
}

// String returns a human-readable description of this Period ID
func (id PeriodId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Quota Counter Key: %q", id.QuotaCounterKey),
		fmt.Sprintf("Quota Period Key: %q", id.QuotaPeriodKey),
	}
	return fmt.Sprintf("Period (%s)", strings.Join(components, "\n"))
}
