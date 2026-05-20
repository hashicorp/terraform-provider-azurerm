package permissionbindings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PermissionBindingId{})
}

var _ resourceids.ResourceId = &PermissionBindingId{}

// PermissionBindingId is a struct representing the Resource ID for a Permission Binding
type PermissionBindingId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	PermissionBindingName string
}

// NewPermissionBindingID returns a new PermissionBindingId struct
func NewPermissionBindingID(subscriptionId string, resourceGroupName string, namespaceName string, permissionBindingName string) PermissionBindingId {
	return PermissionBindingId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		PermissionBindingName: permissionBindingName,
	}
}

// ParsePermissionBindingID parses 'input' into a PermissionBindingId
func ParsePermissionBindingID(input string) (*PermissionBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PermissionBindingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PermissionBindingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePermissionBindingIDInsensitively parses 'input' case-insensitively into a PermissionBindingId
// note: this method should only be used for API response data and not user input
func ParsePermissionBindingIDInsensitively(input string) (*PermissionBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PermissionBindingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PermissionBindingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PermissionBindingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NamespaceName, ok = input.Parsed["namespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", input)
	}

	if id.PermissionBindingName, ok = input.Parsed["permissionBindingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "permissionBindingName", input)
	}

	return nil
}

// ValidatePermissionBindingID checks that 'input' can be parsed as a Permission Binding ID
func ValidatePermissionBindingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePermissionBindingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Permission Binding ID
func (id PermissionBindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/namespaces/%s/permissionBindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.PermissionBindingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Permission Binding ID
func (id PermissionBindingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticPermissionBindings", "permissionBindings", "permissionBindings"),
		resourceids.UserSpecifiedSegment("permissionBindingName", "permissionBindingName"),
	}
}

// String returns a human-readable description of this Permission Binding ID
func (id PermissionBindingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Permission Binding Name: %q", id.PermissionBindingName),
	}
	return fmt.Sprintf("Permission Binding (%s)", strings.Join(components, "\n"))
}
