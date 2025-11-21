package brokerauthentication

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BrokerId{})
}

var _ resourceids.ResourceId = &BrokerId{}

// BrokerId is a struct representing the Resource ID for a Broker
type BrokerId struct {
	SubscriptionId    string
	ResourceGroupName string
	InstanceName      string
	BrokerName        string
}

// NewBrokerID returns a new BrokerId struct
func NewBrokerID(subscriptionId string, resourceGroupName string, instanceName string, brokerName string) BrokerId {
	return BrokerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		InstanceName:      instanceName,
		BrokerName:        brokerName,
	}
}

// ParseBrokerID parses 'input' into a BrokerId
func ParseBrokerID(input string) (*BrokerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BrokerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BrokerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBrokerIDInsensitively parses 'input' case-insensitively into a BrokerId
// note: this method should only be used for API response data and not user input
func ParseBrokerIDInsensitively(input string) (*BrokerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BrokerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BrokerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BrokerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.InstanceName, ok = input.Parsed["instanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instanceName", input)
	}

	if id.BrokerName, ok = input.Parsed["brokerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "brokerName", input)
	}

	return nil
}

// ValidateBrokerID checks that 'input' can be parsed as a Broker ID
func ValidateBrokerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBrokerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Broker ID
func (id BrokerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTOperations/instances/%s/brokers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.InstanceName, id.BrokerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Broker ID
func (id BrokerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftIoTOperations", "Microsoft.IoTOperations", "Microsoft.IoTOperations"),
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceName", "instanceName"),
		resourceids.StaticSegment("staticBrokers", "brokers", "brokers"),
		resourceids.UserSpecifiedSegment("brokerName", "brokerName"),
	}
}

// String returns a human-readable description of this Broker ID
func (id BrokerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Instance Name: %q", id.InstanceName),
		fmt.Sprintf("Broker Name: %q", id.BrokerName),
	}
	return fmt.Sprintf("Broker (%s)", strings.Join(components, "\n"))
}
