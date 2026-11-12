package goldengateconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GoldenGateConnectionId{})
}

var _ resourceids.ResourceId = &GoldenGateConnectionId{}

// GoldenGateConnectionId is a struct representing the Resource ID for a Golden Gate Connection
type GoldenGateConnectionId struct {
	SubscriptionId           string
	ResourceGroupName        string
	GoldenGateConnectionName string
}

// NewGoldenGateConnectionID returns a new GoldenGateConnectionId struct
func NewGoldenGateConnectionID(subscriptionId string, resourceGroupName string, goldenGateConnectionName string) GoldenGateConnectionId {
	return GoldenGateConnectionId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		GoldenGateConnectionName: goldenGateConnectionName,
	}
}

// ParseGoldenGateConnectionID parses 'input' into a GoldenGateConnectionId
func ParseGoldenGateConnectionID(input string) (*GoldenGateConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GoldenGateConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GoldenGateConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGoldenGateConnectionIDInsensitively parses 'input' case-insensitively into a GoldenGateConnectionId
// note: this method should only be used for API response data and not user input
func ParseGoldenGateConnectionIDInsensitively(input string) (*GoldenGateConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GoldenGateConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GoldenGateConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GoldenGateConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.GoldenGateConnectionName, ok = input.Parsed["goldenGateConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "goldenGateConnectionName", input)
	}

	return nil
}

// ValidateGoldenGateConnectionID checks that 'input' can be parsed as a Golden Gate Connection ID
func ValidateGoldenGateConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGoldenGateConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Golden Gate Connection ID
func (id GoldenGateConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/goldenGateConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GoldenGateConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Golden Gate Connection ID
func (id GoldenGateConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticGoldenGateConnections", "goldenGateConnections", "goldenGateConnections"),
		resourceids.UserSpecifiedSegment("goldenGateConnectionName", "goldenGateConnectionName"),
	}
}

// String returns a human-readable description of this Golden Gate Connection ID
func (id GoldenGateConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Golden Gate Connection Name: %q", id.GoldenGateConnectionName),
	}
	return fmt.Sprintf("Golden Gate Connection (%s)", strings.Join(components, "\n"))
}
