package expressroutecrossconnectionarptable

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PeeringArpTableId{})
}

var _ resourceids.ResourceId = &PeeringArpTableId{}

// PeeringArpTableId is a struct representing the Resource ID for a Peering Arp Table
type PeeringArpTableId struct {
	SubscriptionId                  string
	ResourceGroupName               string
	ExpressRouteCrossConnectionName string
	PeeringName                     string
	ArpTableName                    string
}

// NewPeeringArpTableID returns a new PeeringArpTableId struct
func NewPeeringArpTableID(subscriptionId string, resourceGroupName string, expressRouteCrossConnectionName string, peeringName string, arpTableName string) PeeringArpTableId {
	return PeeringArpTableId{
		SubscriptionId:                  subscriptionId,
		ResourceGroupName:               resourceGroupName,
		ExpressRouteCrossConnectionName: expressRouteCrossConnectionName,
		PeeringName:                     peeringName,
		ArpTableName:                    arpTableName,
	}
}

// ParsePeeringArpTableID parses 'input' into a PeeringArpTableId
func ParsePeeringArpTableID(input string) (*PeeringArpTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PeeringArpTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeeringArpTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePeeringArpTableIDInsensitively parses 'input' case-insensitively into a PeeringArpTableId
// note: this method should only be used for API response data and not user input
func ParsePeeringArpTableIDInsensitively(input string) (*PeeringArpTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PeeringArpTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeeringArpTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PeeringArpTableId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExpressRouteCrossConnectionName, ok = input.Parsed["expressRouteCrossConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCrossConnectionName", input)
	}

	if id.PeeringName, ok = input.Parsed["peeringName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "peeringName", input)
	}

	if id.ArpTableName, ok = input.Parsed["arpTableName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "arpTableName", input)
	}

	return nil
}

// ValidatePeeringArpTableID checks that 'input' can be parsed as a Peering Arp Table ID
func ValidatePeeringArpTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePeeringArpTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Peering Arp Table ID
func (id PeeringArpTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCrossConnections/%s/peerings/%s/arpTables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCrossConnectionName, id.PeeringName, id.ArpTableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Peering Arp Table ID
func (id PeeringArpTableId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteCrossConnections", "expressRouteCrossConnections", "expressRouteCrossConnections"),
		resourceids.UserSpecifiedSegment("expressRouteCrossConnectionName", "expressRouteCrossConnectionName"),
		resourceids.StaticSegment("staticPeerings", "peerings", "peerings"),
		resourceids.UserSpecifiedSegment("peeringName", "peeringName"),
		resourceids.StaticSegment("staticArpTables", "arpTables", "arpTables"),
		resourceids.UserSpecifiedSegment("arpTableName", "arpTableName"),
	}
}

// String returns a human-readable description of this Peering Arp Table ID
func (id PeeringArpTableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Cross Connection Name: %q", id.ExpressRouteCrossConnectionName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
		fmt.Sprintf("Arp Table Name: %q", id.ArpTableName),
	}
	return fmt.Sprintf("Peering Arp Table (%s)", strings.Join(components, "\n"))
}
