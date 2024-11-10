package licenses

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LicenseId{})
}

var _ resourceids.ResourceId = &LicenseId{}

// LicenseId is a struct representing the Resource ID for a License
type LicenseId struct {
	SubscriptionId    string
	ResourceGroupName string
	LicenseName       string
}

// NewLicenseID returns a new LicenseId struct
func NewLicenseID(subscriptionId string, resourceGroupName string, licenseName string) LicenseId {
	return LicenseId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LicenseName:       licenseName,
	}
}

// ParseLicenseID parses 'input' into a LicenseId
func ParseLicenseID(input string) (*LicenseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LicenseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LicenseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLicenseIDInsensitively parses 'input' case-insensitively into a LicenseId
// note: this method should only be used for API response data and not user input
func ParseLicenseIDInsensitively(input string) (*LicenseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LicenseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LicenseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LicenseId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LicenseName, ok = input.Parsed["licenseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "licenseName", input)
	}

	return nil
}

// ValidateLicenseID checks that 'input' can be parsed as a License ID
func ValidateLicenseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLicenseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted License ID
func (id LicenseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/licenses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LicenseName)
}

// Segments returns a slice of Resource ID Segments which comprise this License ID
func (id LicenseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticLicenses", "licenses", "licenses"),
		resourceids.UserSpecifiedSegment("licenseName", "licenseName"),
	}
}

// String returns a human-readable description of this License ID
func (id LicenseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("License Name: %q", id.LicenseName),
	}
	return fmt.Sprintf("License (%s)", strings.Join(components, "\n"))
}
