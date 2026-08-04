package cloudhsmclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BackupOperationStatusId{})
}

var _ resourceids.ResourceId = &BackupOperationStatusId{}

// BackupOperationStatusId is a struct representing the Resource ID for a Backup Operation Status
type BackupOperationStatusId struct {
	SubscriptionId      string
	ResourceGroupName   string
	CloudHsmClusterName string
	JobId               string
}

// NewBackupOperationStatusID returns a new BackupOperationStatusId struct
func NewBackupOperationStatusID(subscriptionId string, resourceGroupName string, cloudHsmClusterName string, jobId string) BackupOperationStatusId {
	return BackupOperationStatusId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		CloudHsmClusterName: cloudHsmClusterName,
		JobId:               jobId,
	}
}

// ParseBackupOperationStatusID parses 'input' into a BackupOperationStatusId
func ParseBackupOperationStatusID(input string) (*BackupOperationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupOperationStatusId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupOperationStatusId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBackupOperationStatusIDInsensitively parses 'input' case-insensitively into a BackupOperationStatusId
// note: this method should only be used for API response data and not user input
func ParseBackupOperationStatusIDInsensitively(input string) (*BackupOperationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupOperationStatusId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupOperationStatusId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BackupOperationStatusId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CloudHsmClusterName, ok = input.Parsed["cloudHsmClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudHsmClusterName", input)
	}

	if id.JobId, ok = input.Parsed["jobId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobId", input)
	}

	return nil
}

// ValidateBackupOperationStatusID checks that 'input' can be parsed as a Backup Operation Status ID
func ValidateBackupOperationStatusID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupOperationStatusID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Operation Status ID
func (id BackupOperationStatusId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/%s/backupOperationStatus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudHsmClusterName, id.JobId)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Operation Status ID
func (id BackupOperationStatusId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHardwareSecurityModules", "Microsoft.HardwareSecurityModules", "Microsoft.HardwareSecurityModules"),
		resourceids.StaticSegment("staticCloudHsmClusters", "cloudHsmClusters", "cloudHsmClusters"),
		resourceids.UserSpecifiedSegment("cloudHsmClusterName", "cloudHsmClusterName"),
		resourceids.StaticSegment("staticBackupOperationStatus", "backupOperationStatus", "backupOperationStatus"),
		resourceids.UserSpecifiedSegment("jobId", "jobId"),
	}
}

// String returns a human-readable description of this Backup Operation Status ID
func (id BackupOperationStatusId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Hsm Cluster Name: %q", id.CloudHsmClusterName),
		fmt.Sprintf("Job: %q", id.JobId),
	}
	return fmt.Sprintf("Backup Operation Status (%s)", strings.Join(components, "\n"))
}
