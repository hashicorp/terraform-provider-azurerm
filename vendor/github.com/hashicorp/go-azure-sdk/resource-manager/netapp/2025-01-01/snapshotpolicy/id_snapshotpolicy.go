package snapshotpolicy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SnapshotPolicyId{})
}

var _ resourceids.ResourceId = &SnapshotPolicyId{}

// SnapshotPolicyId is a struct representing the Resource ID for a Snapshot Policy
type SnapshotPolicyId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetAppAccountName  string
	SnapshotPolicyName string
}

// NewSnapshotPolicyID returns a new SnapshotPolicyId struct
func NewSnapshotPolicyID(subscriptionId string, resourceGroupName string, netAppAccountName string, snapshotPolicyName string) SnapshotPolicyId {
	return SnapshotPolicyId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetAppAccountName:  netAppAccountName,
		SnapshotPolicyName: snapshotPolicyName,
	}
}

// ParseSnapshotPolicyID parses 'input' into a SnapshotPolicyId
func ParseSnapshotPolicyID(input string) (*SnapshotPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SnapshotPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SnapshotPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSnapshotPolicyIDInsensitively parses 'input' case-insensitively into a SnapshotPolicyId
// note: this method should only be used for API response data and not user input
func ParseSnapshotPolicyIDInsensitively(input string) (*SnapshotPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SnapshotPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SnapshotPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SnapshotPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetAppAccountName, ok = input.Parsed["netAppAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "netAppAccountName", input)
	}

	if id.SnapshotPolicyName, ok = input.Parsed["snapshotPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "snapshotPolicyName", input)
	}

	return nil
}

// ValidateSnapshotPolicyID checks that 'input' can be parsed as a Snapshot Policy ID
func ValidateSnapshotPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSnapshotPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Snapshot Policy ID
func (id SnapshotPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s/snapshotPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.SnapshotPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Snapshot Policy ID
func (id SnapshotPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetApp", "Microsoft.NetApp", "Microsoft.NetApp"),
		resourceids.StaticSegment("staticNetAppAccounts", "netAppAccounts", "netAppAccounts"),
		resourceids.UserSpecifiedSegment("netAppAccountName", "netAppAccountName"),
		resourceids.StaticSegment("staticSnapshotPolicies", "snapshotPolicies", "snapshotPolicies"),
		resourceids.UserSpecifiedSegment("snapshotPolicyName", "snapshotPolicyName"),
	}
}

// String returns a human-readable description of this Snapshot Policy ID
func (id SnapshotPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Net App Account Name: %q", id.NetAppAccountName),
		fmt.Sprintf("Snapshot Policy Name: %q", id.SnapshotPolicyName),
	}
	return fmt.Sprintf("Snapshot Policy (%s)", strings.Join(components, "\n"))
}
