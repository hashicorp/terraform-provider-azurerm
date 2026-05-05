package cloudexadatainfrastructures

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CloudExadataInfrastructureId{})
}

var _ resourceids.ResourceId = &CloudExadataInfrastructureId{}

// CloudExadataInfrastructureId is a struct representing the Resource ID for a Cloud Exadata Infrastructure
type CloudExadataInfrastructureId struct {
	SubscriptionId                 string
	ResourceGroupName              string
	CloudExadataInfrastructureName string
}

// NewCloudExadataInfrastructureID returns a new CloudExadataInfrastructureId struct
func NewCloudExadataInfrastructureID(subscriptionId string, resourceGroupName string, cloudExadataInfrastructureName string) CloudExadataInfrastructureId {
	return CloudExadataInfrastructureId{
		SubscriptionId:                 subscriptionId,
		ResourceGroupName:              resourceGroupName,
		CloudExadataInfrastructureName: cloudExadataInfrastructureName,
	}
}

// ParseCloudExadataInfrastructureID parses 'input' into a CloudExadataInfrastructureId
func ParseCloudExadataInfrastructureID(input string) (*CloudExadataInfrastructureId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudExadataInfrastructureId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudExadataInfrastructureId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCloudExadataInfrastructureIDInsensitively parses 'input' case-insensitively into a CloudExadataInfrastructureId
// note: this method should only be used for API response data and not user input
func ParseCloudExadataInfrastructureIDInsensitively(input string) (*CloudExadataInfrastructureId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudExadataInfrastructureId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudExadataInfrastructureId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CloudExadataInfrastructureId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CloudExadataInfrastructureName, ok = input.Parsed["cloudExadataInfrastructureName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudExadataInfrastructureName", input)
	}

	return nil
}

// ValidateCloudExadataInfrastructureID checks that 'input' can be parsed as a Cloud Exadata Infrastructure ID
func ValidateCloudExadataInfrastructureID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCloudExadataInfrastructureID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cloud Exadata Infrastructure ID
func (id CloudExadataInfrastructureId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/cloudExadataInfrastructures/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudExadataInfrastructureName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud Exadata Infrastructure ID
func (id CloudExadataInfrastructureId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticCloudExadataInfrastructures", "cloudExadataInfrastructures", "cloudExadataInfrastructures"),
		resourceids.UserSpecifiedSegment("cloudExadataInfrastructureName", "cloudExadataInfrastructureName"),
	}
}

// String returns a human-readable description of this Cloud Exadata Infrastructure ID
func (id CloudExadataInfrastructureId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Exadata Infrastructure Name: %q", id.CloudExadataInfrastructureName),
	}
	return fmt.Sprintf("Cloud Exadata Infrastructure (%s)", strings.Join(components, "\n"))
}
