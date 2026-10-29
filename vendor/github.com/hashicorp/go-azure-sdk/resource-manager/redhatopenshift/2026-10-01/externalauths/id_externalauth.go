package externalauths

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExternalAuthId{})
}

var _ resourceids.ResourceId = &ExternalAuthId{}

// ExternalAuthId is a struct representing the Resource ID for a External Auth
type ExternalAuthId struct {
	SubscriptionId          string
	ResourceGroupName       string
	HcpOpenShiftClusterName string
	ExternalAuthName        string
}

// NewExternalAuthID returns a new ExternalAuthId struct
func NewExternalAuthID(subscriptionId string, resourceGroupName string, hcpOpenShiftClusterName string, externalAuthName string) ExternalAuthId {
	return ExternalAuthId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		HcpOpenShiftClusterName: hcpOpenShiftClusterName,
		ExternalAuthName:        externalAuthName,
	}
}

// ParseExternalAuthID parses 'input' into a ExternalAuthId
func ParseExternalAuthID(input string) (*ExternalAuthId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExternalAuthId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExternalAuthId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExternalAuthIDInsensitively parses 'input' case-insensitively into a ExternalAuthId
// note: this method should only be used for API response data and not user input
func ParseExternalAuthIDInsensitively(input string) (*ExternalAuthId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExternalAuthId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExternalAuthId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExternalAuthId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.HcpOpenShiftClusterName, ok = input.Parsed["hcpOpenShiftClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hcpOpenShiftClusterName", input)
	}

	if id.ExternalAuthName, ok = input.Parsed["externalAuthName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "externalAuthName", input)
	}

	return nil
}

// ValidateExternalAuthID checks that 'input' can be parsed as a External Auth ID
func ValidateExternalAuthID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExternalAuthID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted External Auth ID
func (id ExternalAuthId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/%s/externalAuths/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HcpOpenShiftClusterName, id.ExternalAuthName)
}

// Segments returns a slice of Resource ID Segments which comprise this External Auth ID
func (id ExternalAuthId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRedHatOpenShift", "Microsoft.RedHatOpenShift", "Microsoft.RedHatOpenShift"),
		resourceids.StaticSegment("staticHcpOpenShiftClusters", "hcpOpenShiftClusters", "hcpOpenShiftClusters"),
		resourceids.UserSpecifiedSegment("hcpOpenShiftClusterName", "hcpOpenShiftClusterName"),
		resourceids.StaticSegment("staticExternalAuths", "externalAuths", "externalAuths"),
		resourceids.UserSpecifiedSegment("externalAuthName", "externalAuthName"),
	}
}

// String returns a human-readable description of this External Auth ID
func (id ExternalAuthId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Hcp Open Shift Cluster Name: %q", id.HcpOpenShiftClusterName),
		fmt.Sprintf("External Auth Name: %q", id.ExternalAuthName),
	}
	return fmt.Sprintf("External Auth (%s)", strings.Join(components, "\n"))
}
