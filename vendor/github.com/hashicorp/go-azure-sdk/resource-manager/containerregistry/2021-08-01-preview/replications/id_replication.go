package replications

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ReplicationId{}

// ReplicationId is a struct representing the Resource ID for a Replication
type ReplicationId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	ReplicationName   string
}

// NewReplicationID returns a new ReplicationId struct
func NewReplicationID(subscriptionId string, resourceGroupName string, registryName string, replicationName string) ReplicationId {
	return ReplicationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		ReplicationName:   replicationName,
	}
}

// ParseReplicationID parses 'input' into a ReplicationId
func ParseReplicationID(input string) (*ReplicationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReplicationIDInsensitively parses 'input' case-insensitively into a ReplicationId
// note: this method should only be used for API response data and not user input
func ParseReplicationIDInsensitively(input string) (*ReplicationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReplicationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RegistryName, ok = input.Parsed["registryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "registryName", input)
	}

	if id.ReplicationName, ok = input.Parsed["replicationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "replicationName", input)
	}

	return nil
}

// ValidateReplicationID checks that 'input' can be parsed as a Replication ID
func ValidateReplicationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication ID
func (id ReplicationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/replications/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.ReplicationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication ID
func (id ReplicationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticReplications", "replications", "replications"),
		resourceids.UserSpecifiedSegment("replicationName", "replicationValue"),
	}
}

// String returns a human-readable description of this Replication ID
func (id ReplicationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Replication Name: %q", id.ReplicationName),
	}
	return fmt.Sprintf("Replication (%s)", strings.Join(components, "\n"))
}
