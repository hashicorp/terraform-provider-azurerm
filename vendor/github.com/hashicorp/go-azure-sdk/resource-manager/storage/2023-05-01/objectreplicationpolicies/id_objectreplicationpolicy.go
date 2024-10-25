package objectreplicationpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ObjectReplicationPolicyId{})
}

var _ resourceids.ResourceId = &ObjectReplicationPolicyId{}

// ObjectReplicationPolicyId is a struct representing the Resource ID for a Object Replication Policy
type ObjectReplicationPolicyId struct {
	SubscriptionId            string
	ResourceGroupName         string
	StorageAccountName        string
	ObjectReplicationPolicyId string
}

// NewObjectReplicationPolicyID returns a new ObjectReplicationPolicyId struct
func NewObjectReplicationPolicyID(subscriptionId string, resourceGroupName string, storageAccountName string, objectReplicationPolicyId string) ObjectReplicationPolicyId {
	return ObjectReplicationPolicyId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		StorageAccountName:        storageAccountName,
		ObjectReplicationPolicyId: objectReplicationPolicyId,
	}
}

// ParseObjectReplicationPolicyID parses 'input' into a ObjectReplicationPolicyId
func ParseObjectReplicationPolicyID(input string) (*ObjectReplicationPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ObjectReplicationPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ObjectReplicationPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseObjectReplicationPolicyIDInsensitively parses 'input' case-insensitively into a ObjectReplicationPolicyId
// note: this method should only be used for API response data and not user input
func ParseObjectReplicationPolicyIDInsensitively(input string) (*ObjectReplicationPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ObjectReplicationPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ObjectReplicationPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ObjectReplicationPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageAccountName, ok = input.Parsed["storageAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageAccountName", input)
	}

	if id.ObjectReplicationPolicyId, ok = input.Parsed["objectReplicationPolicyId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "objectReplicationPolicyId", input)
	}

	return nil
}

// ValidateObjectReplicationPolicyID checks that 'input' can be parsed as a Object Replication Policy ID
func ValidateObjectReplicationPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseObjectReplicationPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Object Replication Policy ID
func (id ObjectReplicationPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/objectReplicationPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, id.ObjectReplicationPolicyId)
}

// Segments returns a slice of Resource ID Segments which comprise this Object Replication Policy ID
func (id ObjectReplicationPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("storageAccountName", "storageAccountName"),
		resourceids.StaticSegment("staticObjectReplicationPolicies", "objectReplicationPolicies", "objectReplicationPolicies"),
		resourceids.UserSpecifiedSegment("objectReplicationPolicyId", "objectReplicationPolicyId"),
	}
}

// String returns a human-readable description of this Object Replication Policy ID
func (id ObjectReplicationPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Account Name: %q", id.StorageAccountName),
		fmt.Sprintf("Object Replication Policy: %q", id.ObjectReplicationPolicyId),
	}
	return fmt.Sprintf("Object Replication Policy (%s)", strings.Join(components, "\n"))
}
