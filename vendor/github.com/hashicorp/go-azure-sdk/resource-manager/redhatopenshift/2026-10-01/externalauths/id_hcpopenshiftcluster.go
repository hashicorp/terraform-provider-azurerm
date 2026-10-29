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
	recaser.RegisterResourceId(&HcpOpenShiftClusterId{})
}

var _ resourceids.ResourceId = &HcpOpenShiftClusterId{}

// HcpOpenShiftClusterId is a struct representing the Resource ID for a Hcp Open Shift Cluster
type HcpOpenShiftClusterId struct {
	SubscriptionId          string
	ResourceGroupName       string
	HcpOpenShiftClusterName string
}

// NewHcpOpenShiftClusterID returns a new HcpOpenShiftClusterId struct
func NewHcpOpenShiftClusterID(subscriptionId string, resourceGroupName string, hcpOpenShiftClusterName string) HcpOpenShiftClusterId {
	return HcpOpenShiftClusterId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		HcpOpenShiftClusterName: hcpOpenShiftClusterName,
	}
}

// ParseHcpOpenShiftClusterID parses 'input' into a HcpOpenShiftClusterId
func ParseHcpOpenShiftClusterID(input string) (*HcpOpenShiftClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HcpOpenShiftClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HcpOpenShiftClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHcpOpenShiftClusterIDInsensitively parses 'input' case-insensitively into a HcpOpenShiftClusterId
// note: this method should only be used for API response data and not user input
func ParseHcpOpenShiftClusterIDInsensitively(input string) (*HcpOpenShiftClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HcpOpenShiftClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HcpOpenShiftClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HcpOpenShiftClusterId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateHcpOpenShiftClusterID checks that 'input' can be parsed as a Hcp Open Shift Cluster ID
func ValidateHcpOpenShiftClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHcpOpenShiftClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hcp Open Shift Cluster ID
func (id HcpOpenShiftClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HcpOpenShiftClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hcp Open Shift Cluster ID
func (id HcpOpenShiftClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRedHatOpenShift", "Microsoft.RedHatOpenShift", "Microsoft.RedHatOpenShift"),
		resourceids.StaticSegment("staticHcpOpenShiftClusters", "hcpOpenShiftClusters", "hcpOpenShiftClusters"),
		resourceids.UserSpecifiedSegment("hcpOpenShiftClusterName", "hcpOpenShiftClusterName"),
	}
}

// String returns a human-readable description of this Hcp Open Shift Cluster ID
func (id HcpOpenShiftClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Hcp Open Shift Cluster Name: %q", id.HcpOpenShiftClusterName),
	}
	return fmt.Sprintf("Hcp Open Shift Cluster (%s)", strings.Join(components, "\n"))
}
