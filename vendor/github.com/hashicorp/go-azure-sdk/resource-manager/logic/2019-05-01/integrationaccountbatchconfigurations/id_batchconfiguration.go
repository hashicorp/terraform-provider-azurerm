package integrationaccountbatchconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = BatchConfigurationId{}

// BatchConfigurationId is a struct representing the Resource ID for a Batch Configuration
type BatchConfigurationId struct {
	SubscriptionId         string
	ResourceGroupName      string
	IntegrationAccountName string
	BatchConfigurationName string
}

// NewBatchConfigurationID returns a new BatchConfigurationId struct
func NewBatchConfigurationID(subscriptionId string, resourceGroupName string, integrationAccountName string, batchConfigurationName string) BatchConfigurationId {
	return BatchConfigurationId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		IntegrationAccountName: integrationAccountName,
		BatchConfigurationName: batchConfigurationName,
	}
}

// ParseBatchConfigurationID parses 'input' into a BatchConfigurationId
func ParseBatchConfigurationID(input string) (*BatchConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(BatchConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BatchConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "integrationAccountName", *parsed)
	}

	if id.BatchConfigurationName, ok = parsed.Parsed["batchConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "batchConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseBatchConfigurationIDInsensitively parses 'input' case-insensitively into a BatchConfigurationId
// note: this method should only be used for API response data and not user input
func ParseBatchConfigurationIDInsensitively(input string) (*BatchConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(BatchConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BatchConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "integrationAccountName", *parsed)
	}

	if id.BatchConfigurationName, ok = parsed.Parsed["batchConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "batchConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateBatchConfigurationID checks that 'input' can be parsed as a Batch Configuration ID
func ValidateBatchConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBatchConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Batch Configuration ID
func (id BatchConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/batchConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IntegrationAccountName, id.BatchConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Batch Configuration ID
func (id BatchConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationAccounts", "integrationAccounts", "integrationAccounts"),
		resourceids.UserSpecifiedSegment("integrationAccountName", "integrationAccountValue"),
		resourceids.StaticSegment("staticBatchConfigurations", "batchConfigurations", "batchConfigurations"),
		resourceids.UserSpecifiedSegment("batchConfigurationName", "batchConfigurationValue"),
	}
}

// String returns a human-readable description of this Batch Configuration ID
func (id BatchConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Integration Account Name: %q", id.IntegrationAccountName),
		fmt.Sprintf("Batch Configuration Name: %q", id.BatchConfigurationName),
	}
	return fmt.Sprintf("Batch Configuration (%s)", strings.Join(components, "\n"))
}
