package appattachpackage

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AppAttachPackageId{})
}

var _ resourceids.ResourceId = &AppAttachPackageId{}

// AppAttachPackageId is a struct representing the Resource ID for a App Attach Package
type AppAttachPackageId struct {
	SubscriptionId       string
	ResourceGroupName    string
	AppAttachPackageName string
}

// NewAppAttachPackageID returns a new AppAttachPackageId struct
func NewAppAttachPackageID(subscriptionId string, resourceGroupName string, appAttachPackageName string) AppAttachPackageId {
	return AppAttachPackageId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		AppAttachPackageName: appAttachPackageName,
	}
}

// ParseAppAttachPackageID parses 'input' into a AppAttachPackageId
func ParseAppAttachPackageID(input string) (*AppAttachPackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AppAttachPackageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AppAttachPackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAppAttachPackageIDInsensitively parses 'input' case-insensitively into a AppAttachPackageId
// note: this method should only be used for API response data and not user input
func ParseAppAttachPackageIDInsensitively(input string) (*AppAttachPackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AppAttachPackageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AppAttachPackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AppAttachPackageId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AppAttachPackageName, ok = input.Parsed["appAttachPackageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "appAttachPackageName", input)
	}

	return nil
}

// ValidateAppAttachPackageID checks that 'input' can be parsed as a App Attach Package ID
func ValidateAppAttachPackageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAppAttachPackageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted App Attach Package ID
func (id AppAttachPackageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/appAttachPackages/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AppAttachPackageName)
}

// Segments returns a slice of Resource ID Segments which comprise this App Attach Package ID
func (id AppAttachPackageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDesktopVirtualization", "Microsoft.DesktopVirtualization", "Microsoft.DesktopVirtualization"),
		resourceids.StaticSegment("staticAppAttachPackages", "appAttachPackages", "appAttachPackages"),
		resourceids.UserSpecifiedSegment("appAttachPackageName", "appAttachPackageName"),
	}
}

// String returns a human-readable description of this App Attach Package ID
func (id AppAttachPackageId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("App Attach Package Name: %q", id.AppAttachPackageName),
	}
	return fmt.Sprintf("App Attach Package (%s)", strings.Join(components, "\n"))
}
