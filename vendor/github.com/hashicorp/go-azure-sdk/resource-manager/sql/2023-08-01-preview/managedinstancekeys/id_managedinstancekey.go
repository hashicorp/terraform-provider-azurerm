package managedinstancekeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ManagedInstanceKeyId{})
}

var _ resourceids.ResourceId = &ManagedInstanceKeyId{}

// ManagedInstanceKeyId is a struct representing the Resource ID for a Managed Instance Key
type ManagedInstanceKeyId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ManagedInstanceName string
	KeyName             string
}

// NewManagedInstanceKeyID returns a new ManagedInstanceKeyId struct
func NewManagedInstanceKeyID(subscriptionId string, resourceGroupName string, managedInstanceName string, keyName string) ManagedInstanceKeyId {
	return ManagedInstanceKeyId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ManagedInstanceName: managedInstanceName,
		KeyName:             keyName,
	}
}

// ParseManagedInstanceKeyID parses 'input' into a ManagedInstanceKeyId
func ParseManagedInstanceKeyID(input string) (*ManagedInstanceKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedInstanceKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedInstanceKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedInstanceKeyIDInsensitively parses 'input' case-insensitively into a ManagedInstanceKeyId
// note: this method should only be used for API response data and not user input
func ParseManagedInstanceKeyIDInsensitively(input string) (*ManagedInstanceKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedInstanceKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedInstanceKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedInstanceKeyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedInstanceName, ok = input.Parsed["managedInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedInstanceName", input)
	}

	if id.KeyName, ok = input.Parsed["keyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "keyName", input)
	}

	return nil
}

// ValidateManagedInstanceKeyID checks that 'input' can be parsed as a Managed Instance Key ID
func ValidateManagedInstanceKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedInstanceKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Instance Key ID
func (id ManagedInstanceKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/managedInstances/%s/keys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedInstanceName, id.KeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Instance Key ID
func (id ManagedInstanceKeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticManagedInstances", "managedInstances", "managedInstances"),
		resourceids.UserSpecifiedSegment("managedInstanceName", "managedInstanceName"),
		resourceids.StaticSegment("staticKeys", "keys", "keys"),
		resourceids.UserSpecifiedSegment("keyName", "keyName"),
	}
}

// String returns a human-readable description of this Managed Instance Key ID
func (id ManagedInstanceKeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Instance Name: %q", id.ManagedInstanceName),
		fmt.Sprintf("Key Name: %q", id.KeyName),
	}
	return fmt.Sprintf("Managed Instance Key (%s)", strings.Join(components, "\n"))
}
