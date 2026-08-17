package backupvaultresources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OperationIdId{})
}

var _ resourceids.ResourceId = &OperationIdId{}

// OperationIdId is a struct representing the Resource ID for a Operation Id
type OperationIdId struct {
	SubscriptionId    string
	ResourceGroupName string
	BackupVaultName   string
	OperationId       string
}

// NewOperationIdID returns a new OperationIdId struct
func NewOperationIdID(subscriptionId string, resourceGroupName string, backupVaultName string, operationId string) OperationIdId {
	return OperationIdId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		BackupVaultName:   backupVaultName,
		OperationId:       operationId,
	}
}

// ParseOperationIdID parses 'input' into a OperationIdId
func ParseOperationIdID(input string) (*OperationIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationIdId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOperationIdIDInsensitively parses 'input' case-insensitively into a OperationIdId
// note: this method should only be used for API response data and not user input
func ParseOperationIdIDInsensitively(input string) (*OperationIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationIdId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OperationIdId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.BackupVaultName, ok = input.Parsed["backupVaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backupVaultName", input)
	}

	if id.OperationId, ok = input.Parsed["operationId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "operationId", input)
	}

	return nil
}

// ValidateOperationIdID checks that 'input' can be parsed as a Operation Id ID
func ValidateOperationIdID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationIdID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operation Id ID
func (id OperationIdId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/backupVaults/%s/backupJobs/operations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Operation Id ID
func (id OperationIdId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticBackupVaults", "backupVaults", "backupVaults"),
		resourceids.UserSpecifiedSegment("backupVaultName", "backupVaultName"),
		resourceids.StaticSegment("staticBackupJobs", "backupJobs", "backupJobs"),
		resourceids.StaticSegment("staticOperations", "operations", "operations"),
		resourceids.UserSpecifiedSegment("operationId", "operationId"),
	}
}

// String returns a human-readable description of this Operation Id ID
func (id OperationIdId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Backup Vault Name: %q", id.BackupVaultName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Operation Id (%s)", strings.Join(components, "\n"))
}
