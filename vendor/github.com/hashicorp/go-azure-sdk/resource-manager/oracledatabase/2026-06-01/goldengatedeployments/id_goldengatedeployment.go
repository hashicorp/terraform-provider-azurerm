package goldengatedeployments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GoldenGateDeploymentId{})
}

var _ resourceids.ResourceId = &GoldenGateDeploymentId{}

// GoldenGateDeploymentId is a struct representing the Resource ID for a Golden Gate Deployment
type GoldenGateDeploymentId struct {
	SubscriptionId           string
	ResourceGroupName        string
	GoldenGateDeploymentName string
}

// NewGoldenGateDeploymentID returns a new GoldenGateDeploymentId struct
func NewGoldenGateDeploymentID(subscriptionId string, resourceGroupName string, goldenGateDeploymentName string) GoldenGateDeploymentId {
	return GoldenGateDeploymentId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		GoldenGateDeploymentName: goldenGateDeploymentName,
	}
}

// ParseGoldenGateDeploymentID parses 'input' into a GoldenGateDeploymentId
func ParseGoldenGateDeploymentID(input string) (*GoldenGateDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GoldenGateDeploymentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GoldenGateDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGoldenGateDeploymentIDInsensitively parses 'input' case-insensitively into a GoldenGateDeploymentId
// note: this method should only be used for API response data and not user input
func ParseGoldenGateDeploymentIDInsensitively(input string) (*GoldenGateDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GoldenGateDeploymentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GoldenGateDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GoldenGateDeploymentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.GoldenGateDeploymentName, ok = input.Parsed["goldenGateDeploymentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "goldenGateDeploymentName", input)
	}

	return nil
}

// ValidateGoldenGateDeploymentID checks that 'input' can be parsed as a Golden Gate Deployment ID
func ValidateGoldenGateDeploymentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGoldenGateDeploymentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Golden Gate Deployment ID
func (id GoldenGateDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/goldenGateDeployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GoldenGateDeploymentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Golden Gate Deployment ID
func (id GoldenGateDeploymentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticGoldenGateDeployments", "goldenGateDeployments", "goldenGateDeployments"),
		resourceids.UserSpecifiedSegment("goldenGateDeploymentName", "goldenGateDeploymentName"),
	}
}

// String returns a human-readable description of this Golden Gate Deployment ID
func (id GoldenGateDeploymentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Golden Gate Deployment Name: %q", id.GoldenGateDeploymentName),
	}
	return fmt.Sprintf("Golden Gate Deployment (%s)", strings.Join(components, "\n"))
}
