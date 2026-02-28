package brokerlistener

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ListenerId{})
}

var _ resourceids.ResourceId = &ListenerId{}

// ListenerId is a struct representing the Resource ID for a Listener
type ListenerId struct {
	SubscriptionId    string
	ResourceGroupName string
	InstanceName      string
	BrokerName        string
	ListenerName      string
}

// NewListenerID returns a new ListenerId struct
func NewListenerID(subscriptionId string, resourceGroupName string, instanceName string, brokerName string, listenerName string) ListenerId {
	return ListenerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		InstanceName:      instanceName,
		BrokerName:        brokerName,
		ListenerName:      listenerName,
	}
}

// ParseListenerID parses 'input' into a ListenerId
func ParseListenerID(input string) (*ListenerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ListenerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ListenerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseListenerIDInsensitively parses 'input' case-insensitively into a ListenerId
// note: this method should only be used for API response data and not user input
func ParseListenerIDInsensitively(input string) (*ListenerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ListenerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ListenerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ListenerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ListenerName, ok = input.Parsed["listenerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "listenerName", input)
	}

	return nil
}

// ValidateListenerID checks that 'input' can be parsed as a Listener ID
func ValidateListenerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseListenerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Listener ID
func (id ListenerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTOperations/instances/%s/brokers/%s/listeners/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.InstanceName, id.BrokerName, id.ListenerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Listener ID
func (id ListenerId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticListeners", "listeners", "listeners"),
		resourceids.UserSpecifiedSegment("listenerName", "listenerName"),
	}
}

// String returns a human-readable description of this Listener ID
func (id ListenerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Instance Name: %q", id.InstanceName),
		fmt.Sprintf("Broker Name: %q", id.BrokerName),
		fmt.Sprintf("Listener Name: %q", id.ListenerName),
	}
	return fmt.Sprintf("Listener (%s)", strings.Join(components, "\n"))
}
