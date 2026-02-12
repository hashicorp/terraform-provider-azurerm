package datasets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DatasetId{})
}

var _ resourceids.ResourceId = &DatasetId{}

// DatasetId is a struct representing the Resource ID for a Dataset
type DatasetId struct {
	SubscriptionId    string
	ResourceGroupName string
	FactoryName       string
	DatasetName       string
}

// NewDatasetID returns a new DatasetId struct
func NewDatasetID(subscriptionId string, resourceGroupName string, factoryName string, datasetName string) DatasetId {
	return DatasetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FactoryName:       factoryName,
		DatasetName:       datasetName,
	}
}

// ParseDatasetID parses 'input' into a DatasetId
func ParseDatasetID(input string) (*DatasetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatasetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatasetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatasetIDInsensitively parses 'input' case-insensitively into a DatasetId
// note: this method should only be used for API response data and not user input
func ParseDatasetIDInsensitively(input string) (*DatasetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatasetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatasetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatasetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FactoryName, ok = input.Parsed["factoryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "factoryName", input)
	}

	if id.DatasetName, ok = input.Parsed["datasetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "datasetName", input)
	}

	return nil
}

// ValidateDatasetID checks that 'input' can be parsed as a Dataset ID
func ValidateDatasetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatasetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dataset ID
func (id DatasetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/datasets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FactoryName, id.DatasetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dataset ID
func (id DatasetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataFactory", "Microsoft.DataFactory", "Microsoft.DataFactory"),
		resourceids.StaticSegment("staticFactories", "factories", "factories"),
		resourceids.UserSpecifiedSegment("factoryName", "factoryName"),
		resourceids.StaticSegment("staticDatasets", "datasets", "datasets"),
		resourceids.UserSpecifiedSegment("datasetName", "datasetName"),
	}
}

// String returns a human-readable description of this Dataset ID
func (id DatasetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Factory Name: %q", id.FactoryName),
		fmt.Sprintf("Dataset Name: %q", id.DatasetName),
	}
	return fmt.Sprintf("Dataset (%s)", strings.Join(components, "\n"))
}
