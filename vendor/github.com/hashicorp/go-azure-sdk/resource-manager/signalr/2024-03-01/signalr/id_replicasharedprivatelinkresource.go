package signalr

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReplicaSharedPrivateLinkResourceId{})
}

var _ resourceids.ResourceId = &ReplicaSharedPrivateLinkResourceId{}

// ReplicaSharedPrivateLinkResourceId is a struct representing the Resource ID for a Replica Shared Private Link Resource
type ReplicaSharedPrivateLinkResourceId struct {
	SubscriptionId                string
	ResourceGroupName             string
	SignalRName                   string
	ReplicaName                   string
	SharedPrivateLinkResourceName string
}

// NewReplicaSharedPrivateLinkResourceID returns a new ReplicaSharedPrivateLinkResourceId struct
func NewReplicaSharedPrivateLinkResourceID(subscriptionId string, resourceGroupName string, signalRName string, replicaName string, sharedPrivateLinkResourceName string) ReplicaSharedPrivateLinkResourceId {
	return ReplicaSharedPrivateLinkResourceId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		SignalRName:                   signalRName,
		ReplicaName:                   replicaName,
		SharedPrivateLinkResourceName: sharedPrivateLinkResourceName,
	}
}

// ParseReplicaSharedPrivateLinkResourceID parses 'input' into a ReplicaSharedPrivateLinkResourceId
func ParseReplicaSharedPrivateLinkResourceID(input string) (*ReplicaSharedPrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicaSharedPrivateLinkResourceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicaSharedPrivateLinkResourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReplicaSharedPrivateLinkResourceIDInsensitively parses 'input' case-insensitively into a ReplicaSharedPrivateLinkResourceId
// note: this method should only be used for API response data and not user input
func ParseReplicaSharedPrivateLinkResourceIDInsensitively(input string) (*ReplicaSharedPrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicaSharedPrivateLinkResourceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicaSharedPrivateLinkResourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReplicaSharedPrivateLinkResourceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SignalRName, ok = input.Parsed["signalRName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "signalRName", input)
	}

	if id.ReplicaName, ok = input.Parsed["replicaName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "replicaName", input)
	}

	if id.SharedPrivateLinkResourceName, ok = input.Parsed["sharedPrivateLinkResourceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sharedPrivateLinkResourceName", input)
	}

	return nil
}

// ValidateReplicaSharedPrivateLinkResourceID checks that 'input' can be parsed as a Replica Shared Private Link Resource ID
func ValidateReplicaSharedPrivateLinkResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicaSharedPrivateLinkResourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replica Shared Private Link Resource ID
func (id ReplicaSharedPrivateLinkResourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/signalR/%s/replicas/%s/sharedPrivateLinkResources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SignalRName, id.ReplicaName, id.SharedPrivateLinkResourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replica Shared Private Link Resource ID
func (id ReplicaSharedPrivateLinkResourceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSignalRService", "Microsoft.SignalRService", "Microsoft.SignalRService"),
		resourceids.StaticSegment("staticSignalR", "signalR", "signalR"),
		resourceids.UserSpecifiedSegment("signalRName", "signalRName"),
		resourceids.StaticSegment("staticReplicas", "replicas", "replicas"),
		resourceids.UserSpecifiedSegment("replicaName", "replicaName"),
		resourceids.StaticSegment("staticSharedPrivateLinkResources", "sharedPrivateLinkResources", "sharedPrivateLinkResources"),
		resourceids.UserSpecifiedSegment("sharedPrivateLinkResourceName", "sharedPrivateLinkResourceName"),
	}
}

// String returns a human-readable description of this Replica Shared Private Link Resource ID
func (id ReplicaSharedPrivateLinkResourceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Signal R Name: %q", id.SignalRName),
		fmt.Sprintf("Replica Name: %q", id.ReplicaName),
		fmt.Sprintf("Shared Private Link Resource Name: %q", id.SharedPrivateLinkResourceName),
	}
	return fmt.Sprintf("Replica Shared Private Link Resource (%s)", strings.Join(components, "\n"))
}
