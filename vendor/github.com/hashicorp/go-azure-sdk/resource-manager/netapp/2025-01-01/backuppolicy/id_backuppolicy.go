package backuppolicy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BackupPolicyId{})
}

var _ resourceids.ResourceId = &BackupPolicyId{}

// BackupPolicyId is a struct representing the Resource ID for a Backup Policy
type BackupPolicyId struct {
	SubscriptionId    string
	ResourceGroupName string
	NetAppAccountName string
	BackupPolicyName  string
}

// NewBackupPolicyID returns a new BackupPolicyId struct
func NewBackupPolicyID(subscriptionId string, resourceGroupName string, netAppAccountName string, backupPolicyName string) BackupPolicyId {
	return BackupPolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NetAppAccountName: netAppAccountName,
		BackupPolicyName:  backupPolicyName,
	}
}

// ParseBackupPolicyID parses 'input' into a BackupPolicyId
func ParseBackupPolicyID(input string) (*BackupPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBackupPolicyIDInsensitively parses 'input' case-insensitively into a BackupPolicyId
// note: this method should only be used for API response data and not user input
func ParseBackupPolicyIDInsensitively(input string) (*BackupPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BackupPolicyId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.BackupPolicyName, ok = input.Parsed["backupPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backupPolicyName", input)
	}

	return nil
}

// ValidateBackupPolicyID checks that 'input' can be parsed as a Backup Policy ID
func ValidateBackupPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Policy ID
func (id BackupPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s/backupPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.BackupPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Policy ID
func (id BackupPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetApp", "Microsoft.NetApp", "Microsoft.NetApp"),
		resourceids.StaticSegment("staticNetAppAccounts", "netAppAccounts", "netAppAccounts"),
		resourceids.UserSpecifiedSegment("netAppAccountName", "netAppAccountName"),
		resourceids.StaticSegment("staticBackupPolicies", "backupPolicies", "backupPolicies"),
		resourceids.UserSpecifiedSegment("backupPolicyName", "backupPolicyName"),
	}
}

// String returns a human-readable description of this Backup Policy ID
func (id BackupPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Net App Account Name: %q", id.NetAppAccountName),
		fmt.Sprintf("Backup Policy Name: %q", id.BackupPolicyName),
	}
	return fmt.Sprintf("Backup Policy (%s)", strings.Join(components, "\n"))
}
