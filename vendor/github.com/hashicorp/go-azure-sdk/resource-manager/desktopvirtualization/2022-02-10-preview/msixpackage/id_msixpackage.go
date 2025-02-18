package msixpackage

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MsixPackageId{})
}

var _ resourceids.ResourceId = &MsixPackageId{}

// MsixPackageId is a struct representing the Resource ID for a Msix Package
type MsixPackageId struct {
	SubscriptionId    string
	ResourceGroupName string
	HostPoolName      string
	MsixPackageName   string
}

// NewMsixPackageID returns a new MsixPackageId struct
func NewMsixPackageID(subscriptionId string, resourceGroupName string, hostPoolName string, msixPackageName string) MsixPackageId {
	return MsixPackageId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HostPoolName:      hostPoolName,
		MsixPackageName:   msixPackageName,
	}
}

// ParseMsixPackageID parses 'input' into a MsixPackageId
func ParseMsixPackageID(input string) (*MsixPackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MsixPackageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MsixPackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMsixPackageIDInsensitively parses 'input' case-insensitively into a MsixPackageId
// note: this method should only be used for API response data and not user input
func ParseMsixPackageIDInsensitively(input string) (*MsixPackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MsixPackageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MsixPackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MsixPackageId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.HostPoolName, ok = input.Parsed["hostPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostPoolName", input)
	}

	if id.MsixPackageName, ok = input.Parsed["msixPackageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "msixPackageName", input)
	}

	return nil
}

// ValidateMsixPackageID checks that 'input' can be parsed as a Msix Package ID
func ValidateMsixPackageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMsixPackageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Msix Package ID
func (id MsixPackageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/hostPools/%s/msixPackages/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostPoolName, id.MsixPackageName)
}

// Segments returns a slice of Resource ID Segments which comprise this Msix Package ID
func (id MsixPackageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDesktopVirtualization", "Microsoft.DesktopVirtualization", "Microsoft.DesktopVirtualization"),
		resourceids.StaticSegment("staticHostPools", "hostPools", "hostPools"),
		resourceids.UserSpecifiedSegment("hostPoolName", "hostPoolName"),
		resourceids.StaticSegment("staticMsixPackages", "msixPackages", "msixPackages"),
		resourceids.UserSpecifiedSegment("msixPackageName", "msixPackageName"),
	}
}

// String returns a human-readable description of this Msix Package ID
func (id MsixPackageId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Host Pool Name: %q", id.HostPoolName),
		fmt.Sprintf("Msix Package Name: %q", id.MsixPackageName),
	}
	return fmt.Sprintf("Msix Package (%s)", strings.Join(components, "\n"))
}
