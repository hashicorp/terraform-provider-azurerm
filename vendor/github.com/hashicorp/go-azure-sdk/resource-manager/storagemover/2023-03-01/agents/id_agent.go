package agents

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AgentId{})
}

var _ resourceids.ResourceId = &AgentId{}

// AgentId is a struct representing the Resource ID for a Agent
type AgentId struct {
	SubscriptionId    string
	ResourceGroupName string
	StorageMoverName  string
	AgentName         string
}

// NewAgentID returns a new AgentId struct
func NewAgentID(subscriptionId string, resourceGroupName string, storageMoverName string, agentName string) AgentId {
	return AgentId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StorageMoverName:  storageMoverName,
		AgentName:         agentName,
	}
}

// ParseAgentID parses 'input' into a AgentId
func ParseAgentID(input string) (*AgentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AgentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AgentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAgentIDInsensitively parses 'input' case-insensitively into a AgentId
// note: this method should only be used for API response data and not user input
func ParseAgentIDInsensitively(input string) (*AgentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AgentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AgentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AgentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageMoverName, ok = input.Parsed["storageMoverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageMoverName", input)
	}

	if id.AgentName, ok = input.Parsed["agentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "agentName", input)
	}

	return nil
}

// ValidateAgentID checks that 'input' can be parsed as a Agent ID
func ValidateAgentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAgentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Agent ID
func (id AgentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageMover/storageMovers/%s/agents/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName, id.AgentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Agent ID
func (id AgentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageMover", "Microsoft.StorageMover", "Microsoft.StorageMover"),
		resourceids.StaticSegment("staticStorageMovers", "storageMovers", "storageMovers"),
		resourceids.UserSpecifiedSegment("storageMoverName", "storageMoverName"),
		resourceids.StaticSegment("staticAgents", "agents", "agents"),
		resourceids.UserSpecifiedSegment("agentName", "agentName"),
	}
}

// String returns a human-readable description of this Agent ID
func (id AgentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Mover Name: %q", id.StorageMoverName),
		fmt.Sprintf("Agent Name: %q", id.AgentName),
	}
	return fmt.Sprintf("Agent (%s)", strings.Join(components, "\n"))
}
