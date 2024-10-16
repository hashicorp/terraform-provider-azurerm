package replicationrecoveryservicesproviders

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReplicationRecoveryServicesProviderId{})
}

var _ resourceids.ResourceId = &ReplicationRecoveryServicesProviderId{}

// ReplicationRecoveryServicesProviderId is a struct representing the Resource ID for a Replication Recovery Services Provider
type ReplicationRecoveryServicesProviderId struct {
	SubscriptionId                          string
	ResourceGroupName                       string
	VaultName                               string
	ReplicationFabricName                   string
	ReplicationRecoveryServicesProviderName string
}

// NewReplicationRecoveryServicesProviderID returns a new ReplicationRecoveryServicesProviderId struct
func NewReplicationRecoveryServicesProviderID(subscriptionId string, resourceGroupName string, vaultName string, replicationFabricName string, replicationRecoveryServicesProviderName string) ReplicationRecoveryServicesProviderId {
	return ReplicationRecoveryServicesProviderId{
		SubscriptionId:                          subscriptionId,
		ResourceGroupName:                       resourceGroupName,
		VaultName:                               vaultName,
		ReplicationFabricName:                   replicationFabricName,
		ReplicationRecoveryServicesProviderName: replicationRecoveryServicesProviderName,
	}
}

// ParseReplicationRecoveryServicesProviderID parses 'input' into a ReplicationRecoveryServicesProviderId
func ParseReplicationRecoveryServicesProviderID(input string) (*ReplicationRecoveryServicesProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationRecoveryServicesProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationRecoveryServicesProviderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReplicationRecoveryServicesProviderIDInsensitively parses 'input' case-insensitively into a ReplicationRecoveryServicesProviderId
// note: this method should only be used for API response data and not user input
func ParseReplicationRecoveryServicesProviderIDInsensitively(input string) (*ReplicationRecoveryServicesProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationRecoveryServicesProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationRecoveryServicesProviderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReplicationRecoveryServicesProviderId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VaultName, ok = input.Parsed["vaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vaultName", input)
	}

	if id.ReplicationFabricName, ok = input.Parsed["replicationFabricName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "replicationFabricName", input)
	}

	if id.ReplicationRecoveryServicesProviderName, ok = input.Parsed["replicationRecoveryServicesProviderName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "replicationRecoveryServicesProviderName", input)
	}

	return nil
}

// ValidateReplicationRecoveryServicesProviderID checks that 'input' can be parsed as a Replication Recovery Services Provider ID
func ValidateReplicationRecoveryServicesProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationRecoveryServicesProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Recovery Services Provider ID
func (id ReplicationRecoveryServicesProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationRecoveryServicesProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName, id.ReplicationRecoveryServicesProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Recovery Services Provider ID
func (id ReplicationRecoveryServicesProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultName"),
		resourceids.StaticSegment("staticReplicationFabrics", "replicationFabrics", "replicationFabrics"),
		resourceids.UserSpecifiedSegment("replicationFabricName", "replicationFabricName"),
		resourceids.StaticSegment("staticReplicationRecoveryServicesProviders", "replicationRecoveryServicesProviders", "replicationRecoveryServicesProviders"),
		resourceids.UserSpecifiedSegment("replicationRecoveryServicesProviderName", "replicationRecoveryServicesProviderName"),
	}
}

// String returns a human-readable description of this Replication Recovery Services Provider ID
func (id ReplicationRecoveryServicesProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Fabric Name: %q", id.ReplicationFabricName),
		fmt.Sprintf("Replication Recovery Services Provider Name: %q", id.ReplicationRecoveryServicesProviderName),
	}
	return fmt.Sprintf("Replication Recovery Services Provider (%s)", strings.Join(components, "\n"))
}
