package agentpools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AgentPoolId{})
}

var _ resourceids.ResourceId = &AgentPoolId{}

// AgentPoolId is a struct representing the Resource ID for a Agent Pool
type AgentPoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	AgentPoolName     string
}

// NewAgentPoolID returns a new AgentPoolId struct
func NewAgentPoolID(subscriptionId string, resourceGroupName string, registryName string, agentPoolName string) AgentPoolId {
	return AgentPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		AgentPoolName:     agentPoolName,
	}
}

// ParseAgentPoolID parses 'input' into a AgentPoolId
func ParseAgentPoolID(input string) (*AgentPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AgentPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AgentPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAgentPoolIDInsensitively parses 'input' case-insensitively into a AgentPoolId
// note: this method should only be used for API response data and not user input
func ParseAgentPoolIDInsensitively(input string) (*AgentPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AgentPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AgentPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AgentPoolId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AgentPoolName, ok = input.Parsed["agentPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "agentPoolName", input)
	}

	return nil
}

// ValidateAgentPoolID checks that 'input' can be parsed as a Agent Pool ID
func ValidateAgentPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAgentPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Agent Pool ID
func (id AgentPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/agentPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.AgentPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Agent Pool ID
func (id AgentPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticAgentPools", "agentPools", "agentPools"),
		resourceids.UserSpecifiedSegment("agentPoolName", "agentPoolName"),
	}
}

// String returns a human-readable description of this Agent Pool ID
func (id AgentPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Agent Pool Name: %q", id.AgentPoolName),
	}
	return fmt.Sprintf("Agent Pool (%s)", strings.Join(components, "\n"))
}
