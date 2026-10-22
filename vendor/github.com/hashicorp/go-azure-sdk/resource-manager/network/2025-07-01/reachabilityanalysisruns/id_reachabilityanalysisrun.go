package reachabilityanalysisruns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReachabilityAnalysisRunId{})
}

var _ resourceids.ResourceId = &ReachabilityAnalysisRunId{}

// ReachabilityAnalysisRunId is a struct representing the Resource ID for a Reachability Analysis Run
type ReachabilityAnalysisRunId struct {
	SubscriptionId              string
	ResourceGroupName           string
	NetworkManagerName          string
	VerifierWorkspaceName       string
	ReachabilityAnalysisRunName string
}

// NewReachabilityAnalysisRunID returns a new ReachabilityAnalysisRunId struct
func NewReachabilityAnalysisRunID(subscriptionId string, resourceGroupName string, networkManagerName string, verifierWorkspaceName string, reachabilityAnalysisRunName string) ReachabilityAnalysisRunId {
	return ReachabilityAnalysisRunId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		NetworkManagerName:          networkManagerName,
		VerifierWorkspaceName:       verifierWorkspaceName,
		ReachabilityAnalysisRunName: reachabilityAnalysisRunName,
	}
}

// ParseReachabilityAnalysisRunID parses 'input' into a ReachabilityAnalysisRunId
func ParseReachabilityAnalysisRunID(input string) (*ReachabilityAnalysisRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReachabilityAnalysisRunId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReachabilityAnalysisRunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReachabilityAnalysisRunIDInsensitively parses 'input' case-insensitively into a ReachabilityAnalysisRunId
// note: this method should only be used for API response data and not user input
func ParseReachabilityAnalysisRunIDInsensitively(input string) (*ReachabilityAnalysisRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReachabilityAnalysisRunId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReachabilityAnalysisRunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReachabilityAnalysisRunId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkManagerName, ok = input.Parsed["networkManagerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", input)
	}

	if id.VerifierWorkspaceName, ok = input.Parsed["verifierWorkspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "verifierWorkspaceName", input)
	}

	if id.ReachabilityAnalysisRunName, ok = input.Parsed["reachabilityAnalysisRunName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "reachabilityAnalysisRunName", input)
	}

	return nil
}

// ValidateReachabilityAnalysisRunID checks that 'input' can be parsed as a Reachability Analysis Run ID
func ValidateReachabilityAnalysisRunID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReachabilityAnalysisRunID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Reachability Analysis Run ID
func (id ReachabilityAnalysisRunId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/verifierWorkspaces/%s/reachabilityAnalysisRuns/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.VerifierWorkspaceName, id.ReachabilityAnalysisRunName)
}

// Segments returns a slice of Resource ID Segments which comprise this Reachability Analysis Run ID
func (id ReachabilityAnalysisRunId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerName"),
		resourceids.StaticSegment("staticVerifierWorkspaces", "verifierWorkspaces", "verifierWorkspaces"),
		resourceids.UserSpecifiedSegment("verifierWorkspaceName", "verifierWorkspaceName"),
		resourceids.StaticSegment("staticReachabilityAnalysisRuns", "reachabilityAnalysisRuns", "reachabilityAnalysisRuns"),
		resourceids.UserSpecifiedSegment("reachabilityAnalysisRunName", "reachabilityAnalysisRunName"),
	}
}

// String returns a human-readable description of this Reachability Analysis Run ID
func (id ReachabilityAnalysisRunId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Verifier Workspace Name: %q", id.VerifierWorkspaceName),
		fmt.Sprintf("Reachability Analysis Run Name: %q", id.ReachabilityAnalysisRunName),
	}
	return fmt.Sprintf("Reachability Analysis Run (%s)", strings.Join(components, "\n"))
}
