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
	recaser.RegisterResourceId(&RestoreOperationStatusId{})
}

var _ resourceids.ResourceId = &RestoreOperationStatusId{}

// RestoreOperationStatusId is a struct representing the Resource ID for a Restore Operation Status
type RestoreOperationStatusId struct {
	SubscriptionId      string
	ResourceGroupName   string
	CloudHsmClusterName string
	JobId               string
}

// NewRestoreOperationStatusID returns a new RestoreOperationStatusId struct
func NewRestoreOperationStatusID(subscriptionId string, resourceGroupName string, cloudHsmClusterName string, jobId string) RestoreOperationStatusId {
	return RestoreOperationStatusId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		CloudHsmClusterName: cloudHsmClusterName,
		JobId:               jobId,
	}
}

// ParseRestoreOperationStatusID parses 'input' into a RestoreOperationStatusId
func ParseRestoreOperationStatusID(input string) (*RestoreOperationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RestoreOperationStatusId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RestoreOperationStatusId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRestoreOperationStatusIDInsensitively parses 'input' case-insensitively into a RestoreOperationStatusId
// note: this method should only be used for API response data and not user input
func ParseRestoreOperationStatusIDInsensitively(input string) (*RestoreOperationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RestoreOperationStatusId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RestoreOperationStatusId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RestoreOperationStatusId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateRestoreOperationStatusID checks that 'input' can be parsed as a Restore Operation Status ID
func ValidateRestoreOperationStatusID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRestoreOperationStatusID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Restore Operation Status ID
func (id RestoreOperationStatusId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/%s/restoreOperationStatus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudHsmClusterName, id.JobId)
}

// Segments returns a slice of Resource ID Segments which comprise this Restore Operation Status ID
func (id RestoreOperationStatusId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHardwareSecurityModules", "Microsoft.HardwareSecurityModules", "Microsoft.HardwareSecurityModules"),
		resourceids.StaticSegment("staticCloudHsmClusters", "cloudHsmClusters", "cloudHsmClusters"),
		resourceids.UserSpecifiedSegment("cloudHsmClusterName", "cloudHsmClusterName"),
		resourceids.StaticSegment("staticRestoreOperationStatus", "restoreOperationStatus", "restoreOperationStatus"),
		resourceids.UserSpecifiedSegment("jobId", "jobId"),
	}
}

// String returns a human-readable description of this Restore Operation Status ID
func (id RestoreOperationStatusId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Hsm Cluster Name: %q", id.CloudHsmClusterName),
		fmt.Sprintf("Job: %q", id.JobId),
	}
	return fmt.Sprintf("Restore Operation Status (%s)", strings.Join(components, "\n"))
}
