package disasterrecoveryconfigs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DisasterRecoveryConfigId{})
}

var _ resourceids.ResourceId = &DisasterRecoveryConfigId{}

// DisasterRecoveryConfigId is a struct representing the Resource ID for a Disaster Recovery Config
type DisasterRecoveryConfigId struct {
	SubscriptionId             string
	ResourceGroupName          string
	NamespaceName              string
	DisasterRecoveryConfigName string
}

// NewDisasterRecoveryConfigID returns a new DisasterRecoveryConfigId struct
func NewDisasterRecoveryConfigID(subscriptionId string, resourceGroupName string, namespaceName string, disasterRecoveryConfigName string) DisasterRecoveryConfigId {
	return DisasterRecoveryConfigId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		NamespaceName:              namespaceName,
		DisasterRecoveryConfigName: disasterRecoveryConfigName,
	}
}

// ParseDisasterRecoveryConfigID parses 'input' into a DisasterRecoveryConfigId
func ParseDisasterRecoveryConfigID(input string) (*DisasterRecoveryConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DisasterRecoveryConfigId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DisasterRecoveryConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDisasterRecoveryConfigIDInsensitively parses 'input' case-insensitively into a DisasterRecoveryConfigId
// note: this method should only be used for API response data and not user input
func ParseDisasterRecoveryConfigIDInsensitively(input string) (*DisasterRecoveryConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DisasterRecoveryConfigId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DisasterRecoveryConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DisasterRecoveryConfigId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NamespaceName, ok = input.Parsed["namespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", input)
	}

	if id.DisasterRecoveryConfigName, ok = input.Parsed["disasterRecoveryConfigName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "disasterRecoveryConfigName", input)
	}

	return nil
}

// ValidateDisasterRecoveryConfigID checks that 'input' can be parsed as a Disaster Recovery Config ID
func ValidateDisasterRecoveryConfigID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDisasterRecoveryConfigID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Disaster Recovery Config ID
func (id DisasterRecoveryConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/disasterRecoveryConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.DisasterRecoveryConfigName)
}

// Segments returns a slice of Resource ID Segments which comprise this Disaster Recovery Config ID
func (id DisasterRecoveryConfigId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticDisasterRecoveryConfigs", "disasterRecoveryConfigs", "disasterRecoveryConfigs"),
		resourceids.UserSpecifiedSegment("disasterRecoveryConfigName", "disasterRecoveryConfigName"),
	}
}

// String returns a human-readable description of this Disaster Recovery Config ID
func (id DisasterRecoveryConfigId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Disaster Recovery Config Name: %q", id.DisasterRecoveryConfigName),
	}
	return fmt.Sprintf("Disaster Recovery Config (%s)", strings.Join(components, "\n"))
}
