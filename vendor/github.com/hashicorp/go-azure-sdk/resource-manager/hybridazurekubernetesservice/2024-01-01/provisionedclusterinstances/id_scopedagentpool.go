package provisionedclusterinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedAgentPoolId{})
}

var _ resourceids.ResourceId = &ScopedAgentPoolId{}

// ScopedAgentPoolId is a struct representing the Resource ID for a Scoped Agent Pool
type ScopedAgentPoolId struct {
	ConnectedClusterResourceUri string
	AgentPoolName               string
}

// NewScopedAgentPoolID returns a new ScopedAgentPoolId struct
func NewScopedAgentPoolID(connectedClusterResourceUri string, agentPoolName string) ScopedAgentPoolId {
	return ScopedAgentPoolId{
		ConnectedClusterResourceUri: connectedClusterResourceUri,
		AgentPoolName:               agentPoolName,
	}
}

// ParseScopedAgentPoolID parses 'input' into a ScopedAgentPoolId
func ParseScopedAgentPoolID(input string) (*ScopedAgentPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedAgentPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedAgentPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedAgentPoolIDInsensitively parses 'input' case-insensitively into a ScopedAgentPoolId
// note: this method should only be used for API response data and not user input
func ParseScopedAgentPoolIDInsensitively(input string) (*ScopedAgentPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedAgentPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedAgentPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedAgentPoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ConnectedClusterResourceUri, ok = input.Parsed["connectedClusterResourceUri"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectedClusterResourceUri", input)
	}

	if id.AgentPoolName, ok = input.Parsed["agentPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "agentPoolName", input)
	}

	return nil
}

// ValidateScopedAgentPoolID checks that 'input' can be parsed as a Scoped Agent Pool ID
func ValidateScopedAgentPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedAgentPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Agent Pool ID
func (id ScopedAgentPoolId) ID() string {
	fmtString := "/%s/providers/Microsoft.HybridContainerService/provisionedClusterInstances/default/agentPools/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ConnectedClusterResourceUri, "/"), id.AgentPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Agent Pool ID
func (id ScopedAgentPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("connectedClusterResourceUri", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridContainerService", "Microsoft.HybridContainerService", "Microsoft.HybridContainerService"),
		resourceids.StaticSegment("staticProvisionedClusterInstances", "provisionedClusterInstances", "provisionedClusterInstances"),
		resourceids.StaticSegment("staticDefault", "default", "default"),
		resourceids.StaticSegment("staticAgentPools", "agentPools", "agentPools"),
		resourceids.UserSpecifiedSegment("agentPoolName", "agentPoolName"),
	}
}

// String returns a human-readable description of this Scoped Agent Pool ID
func (id ScopedAgentPoolId) String() string {
	components := []string{
		fmt.Sprintf("Connected Cluster Resource Uri: %q", id.ConnectedClusterResourceUri),
		fmt.Sprintf("Agent Pool Name: %q", id.AgentPoolName),
	}
	return fmt.Sprintf("Scoped Agent Pool (%s)", strings.Join(components, "\n"))
}
