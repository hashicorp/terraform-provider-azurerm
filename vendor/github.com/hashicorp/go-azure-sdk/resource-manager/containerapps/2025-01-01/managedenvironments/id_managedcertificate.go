package managedenvironments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ManagedCertificateId{})
}

var _ resourceids.ResourceId = &ManagedCertificateId{}

// ManagedCertificateId is a struct representing the Resource ID for a Managed Certificate
type ManagedCertificateId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ManagedEnvironmentName string
	ManagedCertificateName string
}

// NewManagedCertificateID returns a new ManagedCertificateId struct
func NewManagedCertificateID(subscriptionId string, resourceGroupName string, managedEnvironmentName string, managedCertificateName string) ManagedCertificateId {
	return ManagedCertificateId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ManagedEnvironmentName: managedEnvironmentName,
		ManagedCertificateName: managedCertificateName,
	}
}

// ParseManagedCertificateID parses 'input' into a ManagedCertificateId
func ParseManagedCertificateID(input string) (*ManagedCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedCertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedCertificateIDInsensitively parses 'input' case-insensitively into a ManagedCertificateId
// note: this method should only be used for API response data and not user input
func ParseManagedCertificateIDInsensitively(input string) (*ManagedCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedCertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedCertificateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedEnvironmentName, ok = input.Parsed["managedEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedEnvironmentName", input)
	}

	if id.ManagedCertificateName, ok = input.Parsed["managedCertificateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedCertificateName", input)
	}

	return nil
}

// ValidateManagedCertificateID checks that 'input' can be parsed as a Managed Certificate ID
func ValidateManagedCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Certificate ID
func (id ManagedCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/managedCertificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName, id.ManagedCertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Certificate ID
func (id ManagedCertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticManagedEnvironments", "managedEnvironments", "managedEnvironments"),
		resourceids.UserSpecifiedSegment("managedEnvironmentName", "managedEnvironmentName"),
		resourceids.StaticSegment("staticManagedCertificates", "managedCertificates", "managedCertificates"),
		resourceids.UserSpecifiedSegment("managedCertificateName", "managedCertificateName"),
	}
}

// String returns a human-readable description of this Managed Certificate ID
func (id ManagedCertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Environment Name: %q", id.ManagedEnvironmentName),
		fmt.Sprintf("Managed Certificate Name: %q", id.ManagedCertificateName),
	}
	return fmt.Sprintf("Managed Certificate (%s)", strings.Join(components, "\n"))
}
