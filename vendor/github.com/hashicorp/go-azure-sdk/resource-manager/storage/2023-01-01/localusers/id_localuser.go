package localusers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LocalUserId{})
}

var _ resourceids.ResourceId = &LocalUserId{}

// LocalUserId is a struct representing the Resource ID for a Local User
type LocalUserId struct {
	SubscriptionId     string
	ResourceGroupName  string
	StorageAccountName string
	LocalUserName      string
}

// NewLocalUserID returns a new LocalUserId struct
func NewLocalUserID(subscriptionId string, resourceGroupName string, storageAccountName string, localUserName string) LocalUserId {
	return LocalUserId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		StorageAccountName: storageAccountName,
		LocalUserName:      localUserName,
	}
}

// ParseLocalUserID parses 'input' into a LocalUserId
func ParseLocalUserID(input string) (*LocalUserId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocalUserId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalUserId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLocalUserIDInsensitively parses 'input' case-insensitively into a LocalUserId
// note: this method should only be used for API response data and not user input
func ParseLocalUserIDInsensitively(input string) (*LocalUserId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocalUserId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalUserId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LocalUserId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.LocalUserName, ok = input.Parsed["localUserName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "localUserName", input)
	}

	return nil
}

// ValidateLocalUserID checks that 'input' can be parsed as a Local User ID
func ValidateLocalUserID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalUserID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local User ID
func (id LocalUserId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/localUsers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, id.LocalUserName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local User ID
func (id LocalUserId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("storageAccountName", "storageAccountName"),
		resourceids.StaticSegment("staticLocalUsers", "localUsers", "localUsers"),
		resourceids.UserSpecifiedSegment("localUserName", "localUserName"),
	}
}

// String returns a human-readable description of this Local User ID
func (id LocalUserId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Account Name: %q", id.StorageAccountName),
		fmt.Sprintf("Local User Name: %q", id.LocalUserName),
	}
	return fmt.Sprintf("Local User (%s)", strings.Join(components, "\n"))
}
