package clouds

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CloudId{})
}

var _ resourceids.ResourceId = &CloudId{}

// CloudId is a struct representing the Resource ID for a Cloud
type CloudId struct {
	SubscriptionId    string
	ResourceGroupName string
	CloudName         string
}

// NewCloudID returns a new CloudId struct
func NewCloudID(subscriptionId string, resourceGroupName string, cloudName string) CloudId {
	return CloudId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		CloudName:         cloudName,
	}
}

// ParseCloudID parses 'input' into a CloudId
func ParseCloudID(input string) (*CloudId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCloudIDInsensitively parses 'input' case-insensitively into a CloudId
// note: this method should only be used for API response data and not user input
func ParseCloudIDInsensitively(input string) (*CloudId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CloudId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CloudName, ok = input.Parsed["cloudName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudName", input)
	}

	return nil
}

// ValidateCloudID checks that 'input' can be parsed as a Cloud ID
func ValidateCloudID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCloudID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cloud ID
func (id CloudId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ScVmm/clouds/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud ID
func (id CloudId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftScVmm", "Microsoft.ScVmm", "Microsoft.ScVmm"),
		resourceids.StaticSegment("staticClouds", "clouds", "clouds"),
		resourceids.UserSpecifiedSegment("cloudName", "cloudName"),
	}
}

// String returns a human-readable description of this Cloud ID
func (id CloudId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Name: %q", id.CloudName),
	}
	return fmt.Sprintf("Cloud (%s)", strings.Join(components, "\n"))
}
