package archives

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PackageId{})
}

var _ resourceids.ResourceId = &PackageId{}

// PackageId is a struct representing the Resource ID for a Package
type PackageId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	PackageName       string
}

// NewPackageID returns a new PackageId struct
func NewPackageID(subscriptionId string, resourceGroupName string, registryName string, packageName string) PackageId {
	return PackageId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		PackageName:       packageName,
	}
}

// ParsePackageID parses 'input' into a PackageId
func ParsePackageID(input string) (*PackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PackageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePackageIDInsensitively parses 'input' case-insensitively into a PackageId
// note: this method should only be used for API response data and not user input
func ParsePackageIDInsensitively(input string) (*PackageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PackageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PackageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PackageId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RegistryName, ok = input.Parsed["registryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "registryName", input)
	}

	if id.PackageName, ok = input.Parsed["packageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "packageName", input)
	}

	return nil
}

// ValidatePackageID checks that 'input' can be parsed as a Package ID
func ValidatePackageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePackageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Package ID
func (id PackageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/packages/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.PackageName)
}

// Segments returns a slice of Resource ID Segments which comprise this Package ID
func (id PackageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticPackages", "packages", "packages"),
		resourceids.UserSpecifiedSegment("packageName", "packageName"),
	}
}

// String returns a human-readable description of this Package ID
func (id PackageId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Package Name: %q", id.PackageName),
	}
	return fmt.Sprintf("Package (%s)", strings.Join(components, "\n"))
}
