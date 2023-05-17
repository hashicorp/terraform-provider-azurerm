package replicationrecoveryplans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ReplicationRecoveryPlanId{}

// ReplicationRecoveryPlanId is a struct representing the Resource ID for a Replication Recovery Plan
type ReplicationRecoveryPlanId struct {
	SubscriptionId              string
	ResourceGroupName           string
	VaultName                   string
	ReplicationRecoveryPlanName string
}

// NewReplicationRecoveryPlanID returns a new ReplicationRecoveryPlanId struct
func NewReplicationRecoveryPlanID(subscriptionId string, resourceGroupName string, vaultName string, replicationRecoveryPlanName string) ReplicationRecoveryPlanId {
	return ReplicationRecoveryPlanId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		VaultName:                   vaultName,
		ReplicationRecoveryPlanName: replicationRecoveryPlanName,
	}
}

// ParseReplicationRecoveryPlanID parses 'input' into a ReplicationRecoveryPlanId
func ParseReplicationRecoveryPlanID(input string) (*ReplicationRecoveryPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationRecoveryPlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationRecoveryPlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.ReplicationRecoveryPlanName, ok = parsed.Parsed["replicationRecoveryPlanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationRecoveryPlanName", *parsed)
	}

	return &id, nil
}

// ParseReplicationRecoveryPlanIDInsensitively parses 'input' case-insensitively into a ReplicationRecoveryPlanId
// note: this method should only be used for API response data and not user input
func ParseReplicationRecoveryPlanIDInsensitively(input string) (*ReplicationRecoveryPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationRecoveryPlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationRecoveryPlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.ReplicationRecoveryPlanName, ok = parsed.Parsed["replicationRecoveryPlanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationRecoveryPlanName", *parsed)
	}

	return &id, nil
}

// ValidateReplicationRecoveryPlanID checks that 'input' can be parsed as a Replication Recovery Plan ID
func ValidateReplicationRecoveryPlanID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationRecoveryPlanID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Recovery Plan ID
func (id ReplicationRecoveryPlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationRecoveryPlans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationRecoveryPlanName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Recovery Plan ID
func (id ReplicationRecoveryPlanId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticReplicationRecoveryPlans", "replicationRecoveryPlans", "replicationRecoveryPlans"),
		resourceids.UserSpecifiedSegment("replicationRecoveryPlanName", "replicationRecoveryPlanValue"),
	}
}

// String returns a human-readable description of this Replication Recovery Plan ID
func (id ReplicationRecoveryPlanId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Recovery Plan Name: %q", id.ReplicationRecoveryPlanName),
	}
	return fmt.Sprintf("Replication Recovery Plan (%s)", strings.Join(components, "\n"))
}
